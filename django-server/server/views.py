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
import uuid

# Create your views here.
@csrf_exempt
@require_POST
@json_view
def createUser(request):
    data = json.loads(request.body)
    username = data['username']
    passwordHash = data['passwordHash']
    user = consulFindByValue('User', 'username', username)

    if user is not None:
        authenticateUser(user, passwordHash)
        client = consulRead('Client', user['clientId'])
    else:
        user = User().__dict__
        client = Client().__dict__
        client['initStatus'] = 'new'
        client['userId'] = user['id']
        consulWrite('Client', client)
        user['clientId'] = client['id']
        user['username'] = username
        user['passwordHash'] = passwordHash
    
    user['authToken'] = uuid.uuid4().hex
    consulWrite('User', user)
    return {
        'user' : {
            'files' : user['files'],
            'id' : user['id'],
            'clientId' : user['clientId'],
            'auth' : user['authToken']
        }
    }

def authenticateUser(user, passwordHash):
    if user['passwordHash'] != passwordHash:
        raise PermissionDenied()
    else:
        return True


def authenticateRequest(request):
    userId = request.META.get('HTTP_USERID', None)
    authToken = request.META.get('HTTP_AUTH', None)
    if userId is None or authToken is None:
        raise PermissionDenied()

    user = consulRead('User', userId)
    if user['authToken'] != authToken:
        raise PermissionDenied()

@require_safe
@json_view
def getClientInitStatus(request, clientId):
    authenticateRequest(request)
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

@require_safe
@json_view
def recoverClient(request, clientId):
    authenticateRequest(request)
    client = consulRead('Client', clientId)
    if client['initStatus'] != 'recover': #TODO: Add security checks
        return {
            'allowed' : False
        }
    user = consulRead('User', client['userId'])
    return {
        'allowed' : True,
        'fileList' : user['files']
    }

def get_client_ip(request):
    x_forwarded_for = request.META.get('HTTP_X_FORWARDED_FOR')
    if x_forwarded_for:
        ip = x_forwarded_for.split(',')[0]
    else:
        ip = request.META.get('REMOTE_ADDR')
    return ip

@csrf_exempt
@require_POST
@json_view
def initClient(request, clientId):
    authenticateRequest(request)
    client = consulRead('Client', clientId)
    if client['initStatus'] != 'new':
        return {
            'message' : 'Client already initialized',
            'error' : 500
        }
    data = json.loads(request.body)
    client['ip'] = get_client_ip(request)
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

def get_client_ip(request):
    x_forwarded_for = request.META.get('HTTP_X_FORWARDED_FOR')
    if x_forwarded_for:
        ip = x_forwarded_for.split(',')[0]
    else:
        ip = request.META.get('REMOTE_ADDR')
    return ip

@csrf_exempt
@require_POST
@json_view
def addFile(request, clientId):
    authenticateRequest(request)
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
    newFile.originalSize = int(data['originalSize'])
    newFile.status = 'added'
    newFile.clientId = clientId
    newFile.blocks = []
    newFile.hash = data['hash']

    if data['blocks'] is not None:
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
def updateFile(request, clientId):
    authenticateRequest(request)
    client = consulRead('Client', clientId)
    data = json.loads(request.body)
    updateFile = consulRead('File', data['id'])
    if updateFile['clientId'] != client['id'] or updateFile['userId'] != client['userId'] :
        return {
            'allowed' : False
        }
    
    availableSpace = int(client['userSpace']) - int(client['userReservedSpace'])
    deltasize = int(updateFile['size']) - int(data['size'])
    if deltasize > availableSpace:
        return {
            'allowed' : False,
            'message' : 'Out of space',
            'error' : 403
        }
    else:
        #This works even if the file has shrunk in size (deltasize will be negative in that case)
        client['userReservedSpace'] = int(client['userReservedSpace']) + deltasize

    consulWrite('Client', client)
    updateFile['name'] = str(data['name'])
    updateFile['deltasize'] = int(deltasize)
    updateFile['size'] = int(data['size'])
    updatedBlocks = []
    for blockInfo in data['blocks']:
        if blockInfo.get('blockId', None) is None:
            block = Block().__dict__
            block['fileId'] = updateFile['id']
            block['offset'] = blockInfo['blockOffset']
            block['shardCount'] = 0
            block['onlineShards'] = 0
            block['shards'] = []
            consulWrite('Block', block)
        else:
            block = consulRead('Block', blockInfo['blockId'])

        updatedBlocks.append(updateBlock(block, blockInfo, updateFile))

    returnInfo = {
        'allowed' : True,
        'id' : updateFile['id'],
        'blocks' : [],
        'shards' : 0,
        'clients' : []
    }

    for blk in updatedBlocks:
        returnInfo['blocks'] += blk['id']
        returnInfo['shards'] += blk['shardCount']
        returnInfo['clients'] += blk['clients']

    return returnInfo

def updateBlock(block, blockInfo, updateFile):
    newShardsInfo = []
    oldShardsInfo = []
    for shardInfo in blockInfo['shards']:
        if shardInfo.get('shardId', None) is None:
            newShardsInfo.append(shardInfo)
        else:
            oldShard = consulRead('Shard', shardInfo['shardId'])
            oldShardsInfo.append(oldShard)

    clients = getShardClients(newShardsInfo)
    newShards = addShardsToBlock(block, newShardsInfo, clients)

    offlineShards, onlineShards = shardsByStatus(oldShardsInfo)
    clients = getShardClients(offlineShards)
    for i in range(0,len(offlineShards)):
        offlineShards[i]['clientId'] = clients[i]
        offlineShards[i]['status'] = 'online'
        consulWrite('Shard', offlineShards[i])
        updateShardClient(offlineShards[i]['clientId'],offlineShards[i]['id'])
    

    updatedBlockInfo = {
        'id' : [block['id']],
        'shardCount' : len(offlineShards + onlineShards + newShards),
        'clients' : []
    }

    for shard in newShards+offlineShards+onlineShards :
        client = consulRead('Client', shard['clientId'])
        updatedBlockInfo['clients'].append({ 
                'id' : shard['id'],
                'blockId': shard['blockId'],
                'offset' : shard['offset'],
                'IP' : client['ip']
            })
    
    return updatedBlockInfo



def shardsByStatus(shards):
    offlines = []
    onlines = []
    for shard in shards:
        assert (shard['status'] == 'offline' or shard['status'] == 'online')
        if shard['status'] == 'offline':
            offlines.append(shard)
        else :
            onlines.append(shard)

    return (offlines, onlines)

@csrf_exempt
@require_POST
@json_view
def commitFile(request, fileId):
    authenticateRequest(request)
    newFile = consulRead('File', fileId)
    data = json.loads(request.body)
    if newFile['status'] != 'added' or newFile['status'] != 'updated' or newFile['clientId'] != data['clientId']:
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

@csrf_exempt
@require_POST
@json_view
def validateShard(request, shardId):
    authenticateRequest(request)
    shard = consulRead('Shard',shardId)
    data = json.loads(request.body)
    receiverId = data['receiverId']
    if shard['clientId'] == receiverId :
        return {
            'accept' : True
        }
    else:
        return {
            'accept' : False
        }

@csrf_exempt
@require_POST
@json_view
def removeFile(request, clientId):
    authenticateRequest(request)
    client = consulRead('Client', clientId)
    data = json.loads(request.body)
    newFile = consulRead('File', data['id'])

    if newFile['clientId'] != clientId or newFile['status'] != 'committed':
        return {
            'allowed' : False
        }

    shardCount = 0
    clients = []
    for blockId in newFile['blocks']:
        block = consulRead('Block', blockId)
        shardCount += int(block['shardCount'])
        for shardId in block['shards']:
            shard = consulRead('Shard', shardId)
            shardClient = consulRead('Client', shard['clientId'])
            clients.append({ 
                'id' : shard['id'],
                'blockId': shard['blockId'],
                'offset' : shard['offset'],
                'IP' : client['ip']
            })

    newFile['status'] = 'removed'
    consulWrite('File', newFile)

    return {
        'allowed' : True,
        'shards' : shardCount,
        'clients' : clients
    }

@csrf_exempt
@require_POST
@json_view
def downloadFile(request, fileId):
    authenticateRequest(request)
    dlFile = consulRead('File', fileId)
    data = json.loads(request.body)
    client = consulRead('Client', data['clientId'])

    if dlFile['clientId'] != client['id'] or dlFile['status'] != 'committed':
        return {
            'allowed' : False
        }

    shardCount = 0
    clients = []
    for blockId in dlFile['blocks']:
        block = consulRead('Block', blockId)
        shardCount += int(block['shardCount'])
        for shardId in block['shards']:
            shard = consulRead('Shard', shardId)
            shardClient = consulRead('Client', shard['clientId'])
            clients.append({ 
                'id' : shard['id'],
                'blockId': shard['blockId'],
                'offset' : shard['offset'],
                'IP' : client['ip']
            })

    return {
        'allowed' : True,
        'shards' : shardCount,
        'clients' : clients
    }

@csrf_exempt
@require_POST
@json_view
def deleteFile(request, fileId):
    authenticateRequest(request)
    delFile = consulRead('File', fileId)
    data = json.loads(request.body)
    if delFile['status'] != 'removed' or delFile['clientId'] != data['clientId']:
        return {
            'error' : 403,
            'message' : "File not available for delete",
            'success' : False
        }

    shardClients = {}
    for blockId in delFile['blocks']:
        block = consulRead('Block', blockId)
        for shardId in block['shards']:
            shard = consulRead('Shard', shardId)
            if shardClients.get(shard['clientId'], None) is None:
                shardClients[shard['clientId']] = int(shard['size'])
            else:
                shardClients[shard['clientId']] = int(shardClients[shard['clientId']]) + int(shard['size']) 
            consulDelete('Shard', shardId)
        consulDelete('Block', blockId)

    fileClient = consulRead('Client', delFile['clientId'])
    fileClient['userSpace'] = int(fileClient['userSpace']) + int(delFile['size'])
    consulWrite('Client', fileClient)
    user = consulRead('File', delFile['userId'])
    user['files'].remove(delFile['id'])
    consulWrite('User', user)
    batchFreeSystemSpace(shardClients)
    consulDelete('File', delFile['id'])

    return {
        'success' : True
    }

def batchFreeSystemSpace(shardClients):
    for clientId in shardClients.keys():
        shardClient = consulRead('Client', clientId)
        shardClient['systemSpace'] = int(shardClient['systemSpace']) + shardClients[clientId]
        consulWrite('Client', shardClient)

@csrf_exempt
@require_POST
@json_view
def invalidateShard(request, shardId):
    authenticateRequest(request)
    shard = consulRead('Shard',shardId)
    data = json.loads(request.body)
    receiverId = data['receiverId']
    ownerId = data['ownerId']
    userFile = consulRead('File', shard['fileId'])
    actualOwnerId = userFile['clientId']
    if shard['clientId'] == receiverId and actualOwnerId == ownerId and userFile['status'] == 'removed':
        return {
            'delete' : True
        }
    else:
        return {
            'delete' : False
        }

# Helpers
def addBlock(blockInfo, newFile):
    block = Block()
    block.offset = blockInfo["blockOffset"]
    block.hash = blockInfo['hash']
    if isinstance(newFile, dict):
        block.fileId = newFile['id']
    else:
        block.fileId = newFile.id
    block.shardCount = 0
    block.onlineShards = 0
    block.shards = []
    clients = getShardClients(shards)
    addShardsToBlock(block, blockInfo['shards'], clients)
    consulWrite('Block', block)
    if isinstance(newFile, dict):
        newFile['blocks'].append(block.id)
    else:
        newFile.blocks.append(block.id)

def addShardsToBlock(block, shards, clients):
    shardIndex = 0
    newShards = []
    for shardInfo in shards:
        newShards.append(addShard(shardInfo, block, shardIndex, clients))
        if isinstance(block, dict):
            block['shardCount'] = int(block['shardCount']) + 1
            block['onlineShards'] = int(block['onlineShards']) + 1
        else:
            block.shardCount = block.shardCount + 1
            block.onlineShards = block.onlineShards + 1 ## Assume the shards will be written correctly, TODO : Check for correct write completion
        shardIndex = shardIndex + 1

    #Return newly created shards, useful in updateFile scenario in updateBlock routine
    return newShards

def addShard(shardInfo, block, shardIndex, clients):
    shard = Shard()
    shard.offset = shardInfo["offset"]
    shard.size = int(shardInfo["size"])
    shard.clientId = clients[shardIndex]['id']
    shard.hash = shardInfo['hash']
    if isinstance(block, dict):
        shard.blockId = block['id']
        shard.fileId = block['fileId']
    else:
        shard.blockId = block.id
        shard.fileId = block.fileId
    
    shard.status = 'online'
    consulWrite('Shard', shard)
    if isinstance(block, dict):
        block['shards'].append(shard.id)
    else:
        block.shards.append(shard.id)    
    updateShardClient(clients[shardIndex]['id'], shard.id)
    return shard.__dict__

def updateShardClient(clientId, shardId):
    shardClient = consulRead('Client', )
    if shardClient['shards'] is not None:
        shardClient['shards'].append(shardId)
    else:
        shardClient['shards'] = [shardId]
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

def consulDelete(root, id):
    s = getConsulateSession()
    try:
        del s.kv[root + '/' + id]
    except AttributeError:
        raise Http404

def consulFindByValue(root, field, value):
    s = getConsulateSession()
    objects = s.kv.find(root).values()
    returnObj = None
    for obj in objects:
        objVal = obj.get(field, None)
        if  objVal == value:
            returnObj = obj
            break

    return returnObj
