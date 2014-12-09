namespace cpp replica
namespace d replica
namespace java replica
namespace php replica
namespace perl replica

struct Replica {
  1: binary shard,
  2: string shardHash,
  3: i32 shardOffset,
  4: string blockId,
  5: string fileId,
  6: string clientId,
}

exception InvalidOperation {
  1: string why
}

service Replicator {
   void ping(),
   void add(1:Replica r) throws (1:InvalidOperation iv),
   void modify(1:Replica r) throws (1:InvalidOperation iv),
   void remove(1:Replica r) throws (1:InvalidOperation iv),
}
