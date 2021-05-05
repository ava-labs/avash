package node

// Metadata struct for storing metadata, available to commands
type Metadata struct {
	Serverhost     string `json:"public-ip"`
	Stakingport    string `json:"staking-port"`
	HTTPport       string `json:"http-port"`
	HTTPTLS        bool   `json:"http-tls-enabled"`
	Dbdir          string `json:"db-dir"`
	Datadir        string `json:"data-dir"`
	Logsdir        string `json:"log-dir"`
	Loglevel       string `json:"log-level"`
	StakingEnabled bool   `json:"staking-enabled"`
	StakerCertPath string `json:"staking-tls-cert-file"`
	StakerKeyPath  string `json:"staking-tls-key-file"`
}
