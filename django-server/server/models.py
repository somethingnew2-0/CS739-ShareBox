# Create your models here.
import uuid

class User(object):
  def __init__(self):
    self.id = uuid.uuid4().hex
    self.clientId = None
    self.files = None

class File(object):
  def __init__(self):
    self.id = uuid.uuid4().hex
    self.name = None
    self.size = None
    self.status = None
    self.blocks = None
    self.clientId = None
    self.userId = None

class Block(object):
  def __init__(self):
    self.id = uuid.uuid4().hex
    self.fileId = None
    self.offset = None
    self.shardCount = None
    self.onlineShards = None
    self.shards = None

class Shard(object):
  def __init__(self):
    self.id = uuid.uuid4().hex
    self.size = None
    self.offset = None
    self.clientId = None
    self.blockId = None
    self.status = None

class Client(object):
  def __init__(self):
    self.id = uuid.uuid4().hex
    self.userId = None
    self.ip = None
    self.systemStatus = None
    self.initStatus = None
    self.userQuota = None
    self.systemQuota = None
    self.userSpace = None
    self.systemSpace = None
    self.userReservedSpace = None
    self.systemReservedSpace = None
    self.shards = None