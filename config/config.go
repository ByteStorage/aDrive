package config

type NameNodeConfig struct {
	// HeartbeatTimeout is the timeout for leader heartbeat
	HeartbeatTimeout int64 `json:"heartbeat_timeout"`
	// Peers is the peers of the raft group
	Peers []string `json:"peers"`
}

type DataNodeConfig struct {
	// HeartbeatTimeout is the timeout for leader heartbeat
	HeartbeatTimeout int64 `json:"heartbeat_timeout"`
	// DataPath is the path to store the data
	DataPath string `json:"data_path"`
}
