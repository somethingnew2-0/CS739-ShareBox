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

type Replicator interface {
	Ping() (err error)
	// Parameters:
	//  - R
	Add(r *Replica) (iv *InvalidOperation, err error)
	// Parameters:
	//  - R
	Modify(r *Replica) (iv *InvalidOperation, err error)
	// Parameters:
	//  - R
	Remove(r *Replica) (iv *InvalidOperation, err error)
}

type ReplicatorClient struct {
	Transport       thrift.TTransport
	ProtocolFactory thrift.TProtocolFactory
	InputProtocol   thrift.TProtocol
	OutputProtocol  thrift.TProtocol
	SeqId           int32
}

func NewReplicatorClientFactory(t thrift.TTransport, f thrift.TProtocolFactory) *ReplicatorClient {
	return &ReplicatorClient{Transport: t,
		ProtocolFactory: f,
		InputProtocol:   f.GetProtocol(t),
		OutputProtocol:  f.GetProtocol(t),
		SeqId:           0,
	}
}

func NewReplicatorClientProtocol(t thrift.TTransport, iprot thrift.TProtocol, oprot thrift.TProtocol) *ReplicatorClient {
	return &ReplicatorClient{Transport: t,
		ProtocolFactory: nil,
		InputProtocol:   iprot,
		OutputProtocol:  oprot,
		SeqId:           0,
	}
}

func (p *ReplicatorClient) Ping() (err error) {
	if err = p.sendPing(); err != nil {
		return
	}
	return p.recvPing()
}

func (p *ReplicatorClient) sendPing() (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("ping", thrift.CALL, p.SeqId)
	args0 := NewPingArgs()
	err = args0.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ReplicatorClient) recvPing() (err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error2 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error3 error
		error3, err = error2.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error3
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result1 := NewPingResult()
	err = result1.Read(iprot)
	iprot.ReadMessageEnd()
	return
}

// Parameters:
//  - R
func (p *ReplicatorClient) Add(r *Replica) (iv *InvalidOperation, err error) {
	if err = p.sendAdd(r); err != nil {
		return
	}
	return p.recvAdd()
}

func (p *ReplicatorClient) sendAdd(r *Replica) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("add", thrift.CALL, p.SeqId)
	args4 := NewAddArgs()
	args4.R = r
	err = args4.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ReplicatorClient) recvAdd() (iv *InvalidOperation, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error6 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error7 error
		error7, err = error6.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error7
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result5 := NewAddResult()
	err = result5.Read(iprot)
	iprot.ReadMessageEnd()
	if result5.Iv != nil {
		iv = result5.Iv
	}
	return
}

// Parameters:
//  - R
func (p *ReplicatorClient) Modify(r *Replica) (iv *InvalidOperation, err error) {
	if err = p.sendModify(r); err != nil {
		return
	}
	return p.recvModify()
}

func (p *ReplicatorClient) sendModify(r *Replica) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("modify", thrift.CALL, p.SeqId)
	args8 := NewModifyArgs()
	args8.R = r
	err = args8.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ReplicatorClient) recvModify() (iv *InvalidOperation, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error10 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error11 error
		error11, err = error10.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error11
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result9 := NewModifyResult()
	err = result9.Read(iprot)
	iprot.ReadMessageEnd()
	if result9.Iv != nil {
		iv = result9.Iv
	}
	return
}

// Parameters:
//  - R
func (p *ReplicatorClient) Remove(r *Replica) (iv *InvalidOperation, err error) {
	if err = p.sendRemove(r); err != nil {
		return
	}
	return p.recvRemove()
}

func (p *ReplicatorClient) sendRemove(r *Replica) (err error) {
	oprot := p.OutputProtocol
	if oprot == nil {
		oprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.OutputProtocol = oprot
	}
	p.SeqId++
	oprot.WriteMessageBegin("remove", thrift.CALL, p.SeqId)
	args12 := NewRemoveArgs()
	args12.R = r
	err = args12.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return
}

func (p *ReplicatorClient) recvRemove() (iv *InvalidOperation, err error) {
	iprot := p.InputProtocol
	if iprot == nil {
		iprot = p.ProtocolFactory.GetProtocol(p.Transport)
		p.InputProtocol = iprot
	}
	_, mTypeId, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return
	}
	if mTypeId == thrift.EXCEPTION {
		error14 := thrift.NewTApplicationException(thrift.UNKNOWN_APPLICATION_EXCEPTION, "Unknown Exception")
		var error15 error
		error15, err = error14.Read(iprot)
		if err != nil {
			return
		}
		if err = iprot.ReadMessageEnd(); err != nil {
			return
		}
		err = error15
		return
	}
	if p.SeqId != seqId {
		err = thrift.NewTApplicationException(thrift.BAD_SEQUENCE_ID, "ping failed: out of sequence response")
		return
	}
	result13 := NewRemoveResult()
	err = result13.Read(iprot)
	iprot.ReadMessageEnd()
	if result13.Iv != nil {
		iv = result13.Iv
	}
	return
}

type ReplicatorProcessor struct {
	processorMap map[string]thrift.TProcessorFunction
	handler      Replicator
}

func (p *ReplicatorProcessor) AddToProcessorMap(key string, processor thrift.TProcessorFunction) {
	p.processorMap[key] = processor
}

func (p *ReplicatorProcessor) GetProcessorFunction(key string) (processor thrift.TProcessorFunction, ok bool) {
	processor, ok = p.processorMap[key]
	return processor, ok
}

func (p *ReplicatorProcessor) ProcessorMap() map[string]thrift.TProcessorFunction {
	return p.processorMap
}

func NewReplicatorProcessor(handler Replicator) *ReplicatorProcessor {

	self16 := &ReplicatorProcessor{handler: handler, processorMap: make(map[string]thrift.TProcessorFunction)}
	self16.processorMap["ping"] = &replicatorProcessorPing{handler: handler}
	self16.processorMap["add"] = &replicatorProcessorAdd{handler: handler}
	self16.processorMap["modify"] = &replicatorProcessorModify{handler: handler}
	self16.processorMap["remove"] = &replicatorProcessorRemove{handler: handler}
	return self16
}

func (p *ReplicatorProcessor) Process(iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	name, _, seqId, err := iprot.ReadMessageBegin()
	if err != nil {
		return false, err
	}
	if processor, ok := p.GetProcessorFunction(name); ok {
		return processor.Process(seqId, iprot, oprot)
	}
	iprot.Skip(thrift.STRUCT)
	iprot.ReadMessageEnd()
	x17 := thrift.NewTApplicationException(thrift.UNKNOWN_METHOD, "Unknown function "+name)
	oprot.WriteMessageBegin(name, thrift.EXCEPTION, seqId)
	x17.Write(oprot)
	oprot.WriteMessageEnd()
	oprot.Flush()
	return false, x17

}

type replicatorProcessorPing struct {
	handler Replicator
}

func (p *replicatorProcessorPing) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewPingArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("ping", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewPingResult()
	if err = p.handler.Ping(); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing ping: "+err.Error())
		oprot.WriteMessageBegin("ping", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("ping", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type replicatorProcessorAdd struct {
	handler Replicator
}

func (p *replicatorProcessorAdd) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewAddArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("add", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewAddResult()
	if result.Iv, err = p.handler.Add(args.R); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing add: "+err.Error())
		oprot.WriteMessageBegin("add", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("add", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type replicatorProcessorModify struct {
	handler Replicator
}

func (p *replicatorProcessorModify) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewModifyArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("modify", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewModifyResult()
	if result.Iv, err = p.handler.Modify(args.R); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing modify: "+err.Error())
		oprot.WriteMessageBegin("modify", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("modify", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

type replicatorProcessorRemove struct {
	handler Replicator
}

func (p *replicatorProcessorRemove) Process(seqId int32, iprot, oprot thrift.TProtocol) (success bool, err thrift.TException) {
	args := NewRemoveArgs()
	if err = args.Read(iprot); err != nil {
		iprot.ReadMessageEnd()
		x := thrift.NewTApplicationException(thrift.PROTOCOL_ERROR, err.Error())
		oprot.WriteMessageBegin("remove", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	iprot.ReadMessageEnd()
	result := NewRemoveResult()
	if result.Iv, err = p.handler.Remove(args.R); err != nil {
		x := thrift.NewTApplicationException(thrift.INTERNAL_ERROR, "Internal error processing remove: "+err.Error())
		oprot.WriteMessageBegin("remove", thrift.EXCEPTION, seqId)
		x.Write(oprot)
		oprot.WriteMessageEnd()
		oprot.Flush()
		return
	}
	if err2 := oprot.WriteMessageBegin("remove", thrift.REPLY, seqId); err2 != nil {
		err = err2
	}
	if err2 := result.Write(oprot); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.WriteMessageEnd(); err == nil && err2 != nil {
		err = err2
	}
	if err2 := oprot.Flush(); err == nil && err2 != nil {
		err = err2
	}
	if err != nil {
		return
	}
	return true, err
}

// HELPER FUNCTIONS AND STRUCTURES

type PingArgs struct {
}

func NewPingArgs() *PingArgs {
	return &PingArgs{}
}

func (p *PingArgs) Read(iprot thrift.TProtocol) error {
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
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *PingArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ping_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *PingArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PingArgs(%+v)", *p)
}

type PingResult struct {
}

func NewPingResult() *PingResult {
	return &PingResult{}
}

func (p *PingResult) Read(iprot thrift.TProtocol) error {
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
		if err := iprot.ReadFieldEnd(); err != nil {
			return err
		}
	}
	if err := iprot.ReadStructEnd(); err != nil {
		return fmt.Errorf("%T read struct end error: %s", p, err)
	}
	return nil
}

func (p *PingResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("ping_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *PingResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("PingResult(%+v)", *p)
}

type AddArgs struct {
	R *Replica `thrift:"r,1"`
}

func NewAddArgs() *AddArgs {
	return &AddArgs{}
}

func (p *AddArgs) Read(iprot thrift.TProtocol) error {
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

func (p *AddArgs) readField1(iprot thrift.TProtocol) error {
	p.R = NewReplica()
	if err := p.R.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.R)
	}
	return nil
}

func (p *AddArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("add_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
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

func (p *AddArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if p.R != nil {
		if err := oprot.WriteFieldBegin("r", thrift.STRUCT, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:r: %s", p, err)
		}
		if err := p.R.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.R)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:r: %s", p, err)
		}
	}
	return err
}

func (p *AddArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AddArgs(%+v)", *p)
}

type AddResult struct {
	Iv *InvalidOperation `thrift:"iv,1"`
}

func NewAddResult() *AddResult {
	return &AddResult{}
}

func (p *AddResult) Read(iprot thrift.TProtocol) error {
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

func (p *AddResult) readField1(iprot thrift.TProtocol) error {
	p.Iv = NewInvalidOperation()
	if err := p.Iv.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Iv)
	}
	return nil
}

func (p *AddResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("add_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	case p.Iv != nil:
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *AddResult) writeField1(oprot thrift.TProtocol) (err error) {
	if p.Iv != nil {
		if err := oprot.WriteFieldBegin("iv", thrift.STRUCT, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:iv: %s", p, err)
		}
		if err := p.Iv.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.Iv)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:iv: %s", p, err)
		}
	}
	return err
}

func (p *AddResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("AddResult(%+v)", *p)
}

type ModifyArgs struct {
	R *Replica `thrift:"r,1"`
}

func NewModifyArgs() *ModifyArgs {
	return &ModifyArgs{}
}

func (p *ModifyArgs) Read(iprot thrift.TProtocol) error {
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

func (p *ModifyArgs) readField1(iprot thrift.TProtocol) error {
	p.R = NewReplica()
	if err := p.R.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.R)
	}
	return nil
}

func (p *ModifyArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("modify_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
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

func (p *ModifyArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if p.R != nil {
		if err := oprot.WriteFieldBegin("r", thrift.STRUCT, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:r: %s", p, err)
		}
		if err := p.R.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.R)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:r: %s", p, err)
		}
	}
	return err
}

func (p *ModifyArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ModifyArgs(%+v)", *p)
}

type ModifyResult struct {
	Iv *InvalidOperation `thrift:"iv,1"`
}

func NewModifyResult() *ModifyResult {
	return &ModifyResult{}
}

func (p *ModifyResult) Read(iprot thrift.TProtocol) error {
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

func (p *ModifyResult) readField1(iprot thrift.TProtocol) error {
	p.Iv = NewInvalidOperation()
	if err := p.Iv.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Iv)
	}
	return nil
}

func (p *ModifyResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("modify_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	case p.Iv != nil:
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *ModifyResult) writeField1(oprot thrift.TProtocol) (err error) {
	if p.Iv != nil {
		if err := oprot.WriteFieldBegin("iv", thrift.STRUCT, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:iv: %s", p, err)
		}
		if err := p.Iv.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.Iv)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:iv: %s", p, err)
		}
	}
	return err
}

func (p *ModifyResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("ModifyResult(%+v)", *p)
}

type RemoveArgs struct {
	R *Replica `thrift:"r,1"`
}

func NewRemoveArgs() *RemoveArgs {
	return &RemoveArgs{}
}

func (p *RemoveArgs) Read(iprot thrift.TProtocol) error {
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

func (p *RemoveArgs) readField1(iprot thrift.TProtocol) error {
	p.R = NewReplica()
	if err := p.R.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.R)
	}
	return nil
}

func (p *RemoveArgs) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("remove_args"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	if err := p.writeField1(oprot); err != nil {
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

func (p *RemoveArgs) writeField1(oprot thrift.TProtocol) (err error) {
	if p.R != nil {
		if err := oprot.WriteFieldBegin("r", thrift.STRUCT, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:r: %s", p, err)
		}
		if err := p.R.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.R)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:r: %s", p, err)
		}
	}
	return err
}

func (p *RemoveArgs) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RemoveArgs(%+v)", *p)
}

type RemoveResult struct {
	Iv *InvalidOperation `thrift:"iv,1"`
}

func NewRemoveResult() *RemoveResult {
	return &RemoveResult{}
}

func (p *RemoveResult) Read(iprot thrift.TProtocol) error {
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

func (p *RemoveResult) readField1(iprot thrift.TProtocol) error {
	p.Iv = NewInvalidOperation()
	if err := p.Iv.Read(iprot); err != nil {
		return fmt.Errorf("%T error reading struct: %s", p.Iv)
	}
	return nil
}

func (p *RemoveResult) Write(oprot thrift.TProtocol) error {
	if err := oprot.WriteStructBegin("remove_result"); err != nil {
		return fmt.Errorf("%T write struct begin error: %s", p, err)
	}
	switch {
	case p.Iv != nil:
		if err := p.writeField1(oprot); err != nil {
			return err
		}
	}
	if err := oprot.WriteFieldStop(); err != nil {
		return fmt.Errorf("%T write field stop error: %s", err)
	}
	if err := oprot.WriteStructEnd(); err != nil {
		return fmt.Errorf("%T write struct stop error: %s", err)
	}
	return nil
}

func (p *RemoveResult) writeField1(oprot thrift.TProtocol) (err error) {
	if p.Iv != nil {
		if err := oprot.WriteFieldBegin("iv", thrift.STRUCT, 1); err != nil {
			return fmt.Errorf("%T write field begin error 1:iv: %s", p, err)
		}
		if err := p.Iv.Write(oprot); err != nil {
			return fmt.Errorf("%T error writing struct: %s", p.Iv)
		}
		if err := oprot.WriteFieldEnd(); err != nil {
			return fmt.Errorf("%T write field end error 1:iv: %s", p, err)
		}
	}
	return err
}

func (p *RemoveResult) String() string {
	if p == nil {
		return "<nil>"
	}
	return fmt.Sprintf("RemoveResult(%+v)", *p)
}
