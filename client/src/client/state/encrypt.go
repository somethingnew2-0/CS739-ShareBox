package state

type Encrypt struct {
	Plaintext []byte
}

func (e Encrypt) Run(sm *StateMachine) {
	encode := &Encode{Ciphertext: make([]byte, len(e.Plaintext))}
	sm.Cipher.Encrypt(encode.Ciphertext, e.Plaintext)
	sm.Add(encode)
}
