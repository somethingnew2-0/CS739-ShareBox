package settings

const (
	ServerAddress = "http://localhost:8000"
	ConfigFile    = "config.json"

	ClientTLS  = false
	ClientPort = 12345

	MinimumWorkers = 8
	MaxStates      = 256

	M           = 12
	K           = 8
	ShardLength = 8192
	BlockSize   = K * ShardLength
)
