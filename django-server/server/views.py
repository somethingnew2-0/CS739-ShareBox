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
        'recovery' : isRecovery
    }

@csrf_exempt
@require_POST
@json_view
def initClient(request, clientId):
    client = consulRead('Client', clientId)
    data = json.loads(request.body)
    client['ip'] = data['IP']
    totalSpace = int(data['space'])
    client['userQuota'] = totalSpace / REPLICATION_FACTOR
    client['systemQuota'] = totalSpace - client['userQuota']
    client['userSpace'] = client['userQuota']
    client['systemSpace'] = client['systemQuota']
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
    if int(data['size']) > int(client['userSpace']):
        raise PermissionDenied()
    else:
        client['userSpace'] = int(client['userSpace']) - int(data['size'])

    newFile = File()
    newFile.name = str(data['name'])
    newFile.size = int(data['size'])
    newFile.status = 'added'
    newFile.blocks = []

    for blockInfo in data['blocks']:
        addBlock(blockInfo, newFile)

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

# Helpers
def addBlock(blockInfo, newFile, clientId):
    block = Block()
    block.offset = blockInfo["offset"]
    block.fileId = newFile.id
    block.shardCount = 0
    block.onlineShards = 0
    block.shards = []
    shardIndex = 0
    for shardInfo in blockInfo['shards']:
        clients = getShardClients(len(blockInfo['shards']))
        addShard(shardInfo, block, shardIndex, clients)
        block.shardCount = block.shardCount + 1
        block.onlineShards = block.onlineShards + 1 ## Assume the shards will be written correctly, TODO : Check for correct write completion
        shardIndex = shardIndex + 1
    consulWrite('Block', block)
    newFile.blocks.append(block.id)

def addShard(shardInfo, block, shardIndex, clients):
    shard = Shard()
    shard.offset = shardInfo["offset"]
    shard.size = shardInfo["size"]
    shard.clientId = clients[shardIndex]
    shard.blockId = block.id
    shard.status = 'online'
    consulWrite('Shard', shard)
    block.shards.append(shard.id)
    shardClient = consulRead('Client', clients[shardIndex])
    shardClient['shards'].append(shard.id)
    consulWrite('Client', shardClient)

def getShardClients(shardCount):
    #Round robin assignment of available online clients
    onlineClients = getOnlineClients()
    numClients = len(onlineClients)
    shardClients = []
    for i in range(0, shardCount):
        shardClients.append(onlineClients[i % numClients])
    return shardClients

def getOnlineClients():
    s = getConsulateSession()
    clients = s.kv.find('Client').values()
    onlineClients = []
    for client in clients:
        if client["systemStatus"] == "online":
            onlineClients.append(client.id)

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