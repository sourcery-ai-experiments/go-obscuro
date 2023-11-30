package config

// Config contains the configuration required by the WalletExtension.
type Config struct {
	TenGatewayHost          string
	TenGatewayPortHTTP      int
	TenGatewayPortWS        int
	NodeRPCHTTPAddress      string
	NodeRPCWebsocketAddress string
	LogPath                 string
	DBPathOverride          string // Overrides the database file location. Used in tests.
	VerboseFlag             bool
	DBType                  string
	DBConnectionURL         string
	TenChainID              int
}
