package settings

const (
	ConsulAddress = "docker:8500"
	ServerAddress = "http://localhost:8000"

	ConfigFile = "config.json"

	ClientTLS  = false
	ClientPort = 12345

	// Seconds to wait for create file to occur before updating it
	UpdateTimeout = 10

	ReplicasPath = "replicas"

	MinimumWorkers = 8
	MaxStates      = 256

	M           = 12
	K           = 8
	ShardLength = 8192
	BlockSize   = K * ShardLength
)
