package config

import (
	"github.com/urfave/cli/v2"

	"github.com/dapplink-external/risk-server/flags"
)

type Config struct {
	Migrations    string
	RpcServer     ServerConfig
	MetricsServer ServerConfig
	ApiGateWayRpc string
	LevelDbPath   string
	RiskKey       string
	AccessToken   string
}

type ServerConfig struct {
	Host string
	Port int
}

func LoadConfig(cliCtx *cli.Context) (Config, error) {
	var cfg Config
	cfg = NewConfig(cliCtx)
	return cfg, nil
}

func NewConfig(ctx *cli.Context) Config {
	return Config{
		RpcServer: ServerConfig{
			Host: ctx.String(flags.RpcHostFlag.Name),
			Port: ctx.Int(flags.RpcPortFlag.Name),
		},
		MetricsServer: ServerConfig{
			Host: ctx.String(flags.MetricsHostFlag.Name),
			Port: ctx.Int(flags.MetricsPortFlag.Name),
		},
		ApiGateWayRpc: ctx.String(flags.ApiGateWayRpcFlag.Name),
		LevelDbPath:   ctx.String(flags.LevelDbPathFlag.Name),
		RiskKey:       ctx.String(flags.RiskKeyFlag.Name),
		AccessToken:   ctx.String(flags.AccessTokenFlag.Name),
	}
}
