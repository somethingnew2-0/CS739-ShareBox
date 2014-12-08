from django.shortcuts import render
from django.http import HttpResponseForbidden, HttpResponse
from django.core.exceptions import PermissionDenied
from django.views.decorators.csrf import csrf_exempt
from django.views.decorators.http import require_POST, require_safe
from django.http import Http404
from jsonview.decorators import json_view
import consulate
import json
from models import *
from random import shuffle

# Create your views here.
@csrf_exempt
@require_POST
@json_view
def createUser(request):
    user = User()
    client = Client()
    client.initStatus = 'new'
    client.userId = user.id
    consulWrite('Client', client)
    user.clientId = client.id
    consulWrite('User', user)
    return { 'user' : user.__dict__ }

@require_safe
@json_view
def getClientInitStatus(request, clientId):
    client = consulRead('Client', clientId)
    if client['initStatus'] == 'new':
        isNew = True
        isRecovery = False
    elif client['initStatus'] == 'recovery':
        isNew = False
        isRecovery = True
    else:
        raise Http404

    return {
        'new' : isNew,
        'recovery' : isRecovery,
        'usuable' : client['userSpace'],
        'userReserved' : client['userReservedSpace'],
        'system' : client['systemSpace'],
        'systemReserved': client['systemReservedSpace']
    }

@csrf_exempt
@require_POST
@json_view
def initClient(request, clientId):
    client = consulRead('Client', clientId)
    if client['initStatus'] != 'new':
        return {
            'message' : 'Client already initialized',
            'error' : 500
        }
    data = json.loads(request.body)
    client['ip'] = data['IP']
    totalSpace = int(data['space'])
    client['userQuota'] = totalSpace / REPLICATION_FACTOR
    client['systemQuota'] = totalSpace - client['userQuota']
    client['userSpace'] = client['userQuota']
    client['systemSpace'] = client['systemQuota']
    client['userReservedSpace'] = 0
    client['systemReservedSpace'] = 0
    client['copies'] = 0
    client['systemStatus'] = 'online'
    client['initStatus'] = 'recovery'
    consulWrite('Client', client)
    return {
        'usable' : client['userSpace'],
        'system' : client['systemSpace']
    }

@csrf_exempt
@require_POST
@json_view
def addFile(request, clientId):
    client = consulRead('Client', clientId)
    data = json.loads(request.body)
    availableSpace = int(client['userSpace']) - int(client['userReservedSpace'])
    if int(data['size']) > availableSpace:
        return {
            'allowed' : False,
            'message' : 'Out of space',
            'error' : 403
        }
    else:
        client['userReservedSpace'] = int(client['userReservedSpace']) + int(data['size'])

    consulWrite('Client', client)

    newFile = File()
    newFile.name = str(data['name'])
    newFile.size = int(data['size'])
    newFile.status = 'added'
    newFile.clientId = clientId
    newFile.blocks = []

    for blockInfo in data['blocks']:
        addBlock(blockInfo, newFile)

    user = consulRead('User', client['userId'])
    if user['files'] is None:
        user['files'] = []
    user['files'].append(newFile.id)
    newFile.userId = user['id']
    consulWrite('User', user)
    consulWrite('File', newFile)
    shardCount = 0
    clients = []
    for blockId in newFile.blocks:
        block = consulRead('Block', blockId)
        shardCount += block['shardCount']
        for shardId in block['shards']:
            shard = consulRead('Shard', shardId)
            client = consulRead('Client', shard['clientId'])
            clients.append({ 
                'id' : shard['id'],
                'blockId': shard['blockId'],
                'offset' : shard['offset'],
                'IP' : client['ip']
            })
    return {
        'allowed' : True,
        'id' : newFile.id,
        'blocks' : newFile.blocks,
        'shards' : shardCount,
        'clients' : clients
    }

@csrf_exempt
@require_POST
@json_view
def commitFile(request, fileId):
    newFile = consulRead('File', fileId)
    if newFile['status'] != 'added':
        return {
            'error' : 403,
            'message' : "File not available for commit",
            'success' : False
        }

    client = consulRead('Client', newFile['clientId'])
    client['userSpace'] = int(client['userSpace']) - int(newFile['size'])
    client['userReservedSpace'] = int(client['userReservedSpace']) - int(newFile['size'])
    consulWrite('Client', client)
    reservations = {}
    for blockId in newFile['blocks']:
        block = consulRead('Block', blockId)
        for shardId in block['shards']:
            shard = consulRead('Shard', shardId)
            if reservations.get(shard['clientId'], None) is None:
                reservations[shard['clientId']] = int(shard['size'])
            else:
                reservations[shard['clientId']] = int(reservations[shard['clientId']]) + int(shard['size']) 

    batchCommitReservations(reservations)
    newFile['status'] = 'committed'
    consulWrite('File', newFile)
    return {
        'success': True
    }

def batchCommitReservations(reservations):
    for clientId in reservations.keys():
        shardClient = consulRead('Client', clientId)
        shardClient['systemReservedSpace'] = int(shardClient['systemReservedSpace']) - reservations[clientId]
        shardClient['systemSpace'] = int(shardClient['systemSpace']) - reservations[clientId]
        consulWrite('Client', shardClient)

# Helpers
def addBlock(blockInfo, newFile):
    block = Block()
    block.offset = blockInfo["blockOffset"]
    block.fileId = newFile.id
    block.shardCount = 0
    block.onlineShards = 0
    block.shards = []
    shardIndex = 0
    clients = getShardClients(blockInfo['shards'])
    for shardInfo in blockInfo['shards']:
        addShard(shardInfo, block, shardIndex, clients)
        block.shardCount = block.shardCount + 1
        block.onlineShards = block.onlineShards + 1 ## Assume the shards will be written correctly, TODO : Check for correct write completion
        shardIndex = shardIndex + 1
    consulWrite('Block', block)
    newFile.blocks.append(block.id)

def addShard(shardInfo, block, shardIndex, clients):
    shard = Shard()
    shard.offset = shardInfo["offset"]
    shard.size = int(shardInfo["size"])
    shard.clientId = clients[shardIndex]['id']
    shard.blockId = block.id
    shard.status = 'online'
    consulWrite('Shard', shard)
    block.shards.append(shard.id)
    shardClient = consulRead('Client', clients[shardIndex]['id'])
    if shardClient['shards'] is not None:
        shardClient['shards'].append(shard.id)
    else:
        shardClient['shards'] = [shard.id]
    consulWrite('Client', shardClient)

def getShardClients(shards):
    #Greedy assignment of available online clients
    #Always fills one client before going to the next one
    onlineClients = getOnlineClients()
    numClients = len(onlineClients)
    clientReservations = []
    for shardInfo in shards:
        clientFound = False
        for clientId in onlineClients:
            client = consulRead('Client', clientId)
            availableSystemSpace = int(client['systemSpace']) - int(client['systemReservedSpace'])
            if availableSystemSpace > int(shardInfo['size']):
                client['systemReservedSpace'] = int(client['systemReservedSpace']) + int(shardInfo['size'])
                clientReservations.append({'id' : clientId, 'space': int(shardInfo['size'])})
                #print str(client['id']) + " - " + str(client['ip']) + " - " + str(client['systemReservedSpace']) + " \n"
                consulWrite('Client', client)
                clientFound = True
                break
        if not clientFound:
            releaseReservations(clientReservations)
            return []

    return clientReservations

def releaseReservations(clientReservations):
    for reservation in clientReservations:
        client = consulRead('Client', reservation['id'])
        client['systemReservedSpace'] = int(client['systemReservedSpace']) - int(reservation['space'])
        consulWrite('Client', client)

def getOnlineClients():
    s = getConsulateSession()
    clients = s.kv.find('Client').values()
    onlineClients = []
    for client in clients:
        if client["systemStatus"] == "online":
            onlineClients.append(client["id"])

    shuffle(onlineClients)
    return onlineClients


REPLICATION_FACTOR = 3
def getConsulateSession():
    return consulate.Consulate()

def consulWrite(root, obj):
    s = getConsulateSession()
    if isinstance(obj, dict):
        s.kv[root + '/' + obj['id']] = obj
    else:
        s.kv[root + '/' + obj.id] = obj.__dict__

def consulRead(root, id):
    s = getConsulateSession()
    try:
        return s.kv[root + '/' + id]
    except AttributeError:
        raise Http404