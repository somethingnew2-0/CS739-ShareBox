namespace cpp replica
namespace d replica
namespace java replica
namespace php replica
namespace perl replica

struct Replica {
  1: binary shard,
  2: string shardHash,
  3: i32 shardOffset,
  4: string shardId,
  5: string blockId,
  6: string fileId,
  7: string clientId,
}

service Replicator {
   void ping(),
   void add(1:Replica r) 
   void modify(1:Replica r) 
   void remove(1:string shardId) 
   Replica download(1:string shardId) 
}
