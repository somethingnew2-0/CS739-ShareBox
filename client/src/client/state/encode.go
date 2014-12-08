package state

type Encode struct {
	Ciphertext []byte
}

func (e Encode) Run(sm *StateMachine) {
	upload := &Upload{}
	upload.EncodedCiphertext = append(e.Ciphertext, sm.ErasureCode.Encode(e.Ciphertext)...)
	sm.Add(upload)
}
