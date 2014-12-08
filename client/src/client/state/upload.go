package state

type Upload struct {
	EncodedCiphertext []byte
}

func (u Upload) Run(sm *StateMachine) {
}
