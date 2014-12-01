package settings

const ServerAddress = "http://localhost:8000"
const ConfigFile = "config.json"

const MinimumWorkers = 16
const MaxStates = 256

const M = 12
const K = 8
const ShardLength = 8192
const BlockSize = K * ShardLength
