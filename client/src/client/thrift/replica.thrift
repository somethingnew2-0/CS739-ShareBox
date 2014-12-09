namespace cpp replica
namespace d replica
namespace java replica
namespace php replica
namespace perl replica

enum Operation {
  ADD = 1,
  MODIFY = 2,
  REMOVE = 3
}

struct Replica {
  1: binary shard,
  2: string shardHash,
  3: i32 shardOffset,
  4: string blockId,
  5: string fileId,
  6: string clientId,
  7: Operation op,
}

exception InvalidOperation {
  1: string why
}

service Replicator {
   void ping(),
   void add(1:Replica r) throws (1:InvalidOperation err),
   void modify(1:Replica r) throws (1:InvalidOperation err),
   void remove(1:Replica r) throws (1:InvalidOperation err),
}
