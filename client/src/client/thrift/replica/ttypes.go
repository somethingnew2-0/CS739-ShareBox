// Autogenerated by Thrift Compiler (0.9.1)
// DO NOT EDIT UNLESS YOU ARE SURE THAT YOU KNOW WHAT YOU ARE DOING

package replica

import (
	"fmt"
	"git.apache.org/thrift.git/lib/go/thrift"
	"math"
)

// (needed to ensure safety because of naive import list construction.)
var _ = math.MinInt32
var _ = thrift.ZERO
var _ = fmt.Printf

var GoUnusedProtection__ int

type Replica struct {
	Shard       []byte `thrift:"shard,1"`
	ShardHash   []byte `thrift:"shardHash,2"`
	ShardOffset int32  `thrift:"shardOffset,3"`
	ShardId     string `thrift:"shardId,4"`
	BlockId     string `thrift:"blockId,5"`
	FileId      string `thrift:"fileId,6"`
	ClientId    string `thrift:"clientId,7"`
}

func NewReplica() *Replica {
	return &Replica{}
}

func (p *Replica) Read(iprot thrift.TProtocol) error {
	if _, err := iprot.ReadStructBegin(); err != nil {
		return fmt.Errorf("%T read error", p)
	}
	for {
		_, fieldTypeId, fieldId, err := iprot.ReadFieldBegin()
		if err != nil {
			return fmt.Errorf("%T field %d read error: %s", p, fieldId, err)
		}
		if fieldTypeId == thrift.STOP {
			break
		}
		switch fieldId {
		case 1:
			if err := p.readField1(iprot); err != nil {
				return err
			}
		case 2:
			if err := p.readField2(iprot); err != nil {
				return err
			}
		case 3:
			if err := p.readField3(iprot); err != nil {
				return err
			}
		case 4:
			if err := p.readField4(iprot); err != nil {
				return err
			}
		case 5:
			if err := p.readField5(iprot); err != nil {
				return err
			}
		case 6:
			if err := p.readField6(iprot); err != nil {
				return err
			}
		case 7:
			if err := p.readField7(iprot); err != nil {
				return err
			}
		default:
			if err := iprot.Skip(fieldTypeId); err != nil {
				return err
			}
		}
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *Replica) readField1(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return fmt.Errorf("error reading field 1: %s")
	} else {
		p.Shard = v
	}
	return nil
}

func (p *Replica) readField2(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadBinary(); err != nil {
		return fmt.Errorf("error reading field 2: %s")
	} else {
		p.ShardHash = v
	}
	return nil
}

func (p *Replica) readField3(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadI32(); err != nil {
		return fmt.Errorf("error reading field 3: %s")
	} else {
		p.ShardOffset = v
	}
	return nil
}

func (p *Replica) readField4(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 4: %s")
	} else {
		p.ShardId = v
	}
	return nil
}

func (p *Replica) readField5(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 5: %s")
	} else {
		p.BlockId = v
	}
	return nil
}

func (p *Replica) readField6(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 6: %s")
	} else {
		p.FileId = v
	}
	return nil
}

func (p *Replica) readField7(iprot thrift.TProtocol) error {
	if v, err := iprot.ReadString(); err != nil {
		return fmt.Errorf("error reading field 7: %s")
	} else {
		p.ClientId = v
	}
	return nil
}

func (p *Replica) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("Replica"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
		return err
	}
	if err := p.writeField2(oprot); err != nil {
		return err
	}
	if err := p.writeField3(oprot); err != nil {
		return err
	}
	if err := p.writeField4(oprot); err != nil {
		return err
	}
	if err := p.writeField5(oprot); err != nil {
		return err
	}
	if err := p.writeField6(oprot); err != nil {
		return err
	}
	if err := p.writeField7(oprot); err != nil {
		return err
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *Replica) writeField1(oprot thrift.TProtocol) (err error) {
	if p.Shard != nil {
		if err := oprot.WriteFieldBegin("shard", thrift.BINARY, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:shard: %s", p, err)
		}
		if err := oprot.WriteBinary(p.Shard); err != nil {
			return fmt.Errorf("%T.shard (1) field write error: %s", p)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:shard: %s", p, err)
		}
	}
	return err
}

func (p *Replica) writeField2(oprot thrift.TProtocol) (err error) {
	if p.ShardHash != nil {
		if err := oprot.WriteFieldBegin("shardHash", thrift.BINARY, 2); err != nil {
			return fmt.Errorf("%T write field begin error 2:shardHash: %s", p, err)
		}
		if err := oprot.WriteBinary(p.ShardHash); err != nil {
			return fmt.Errorf("%T.shardHash (2) field write error: %s", p)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 2:shardHash: %s", p, err)
		}
	}
	return err
}

func (p *Replica) writeField3(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("shardOffset", thrift.I32, 3); err != nil {
		return fmt.Errorf("%T write field begin error 3:shardOffset: %s", p, err)
	}
	if err := oprot.WriteI32(int32(p.ShardOffset)); err != nil {
		return fmt.Errorf("%T.shardOffset (3) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 3:shardOffset: %s", p, err)
	}
	return err
}

func (p *Replica) writeField4(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("shardId", thrift.STRING, 4); err != nil {
		return fmt.Errorf("%T write field begin error 4:shardId: %s", p, err)
	}
	if err := oprot.WriteString(string(p.ShardId)); err != nil {
		return fmt.Errorf("%T.shardId (4) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 4:shardId: %s", p, err)
	}
	return err
}

func (p *Replica) writeField5(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("blockId", thrift.STRING, 5); err != nil {
		return fmt.Errorf("%T write field begin error 5:blockId: %s", p, err)
	}
	if err := oprot.WriteString(string(p.BlockId)); err != nil {
		return fmt.Errorf("%T.blockId (5) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 5:blockId: %s", p, err)
	}
	return err
}

func (p *Replica) writeField6(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("fileId", thrift.STRING, 6); err != nil {
		return fmt.Errorf("%T write field begin error 6:fileId: %s", p, err)
	}
	if err := oprot.WriteString(string(p.FileId)); err != nil {
		return fmt.Errorf("%T.fileId (6) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 6:fileId: %s", p, err)
	}
	return err
}

func (p *Replica) writeField7(oprot thrift.TProtocol) (err error) {
	if err := oprot.WriteFieldBegin("clientId", thrift.STRING, 7); err != nil {
		return fmt.Errorf("%T write field begin error 7:clientId: %s", p, err)
	}
	if err := oprot.WriteString(string(p.ClientId)); err != nil {
		return fmt.Errorf("%T.clientId (7) field write error: %s", p)
	}
	if err := oprot.WriteFieldEnd(); err != nil {
		return fmt.Errorf("%T write field end error 7:clientId: %s", p, err)
	}
	return err
}

func (p *Replica) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("Replica(%+v)", *p)
}
