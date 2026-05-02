package config

import (
	"github.com/urfave/cli/v2"

	"github.com/the-web3/mock-risk-server/flags"
)

type Config struct {
	Migrations    string
	RpcServer     ServerConfig
	MetricsServer ServerConfig
	ApiGateWayRpc string
	DbConf        DBConfig
}

type DBConfig struct {
	Host     string
	Port     int
	Name     string
	User     string
	Password string
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
		Migrations: ctx.String(flags.MigrationsFlag.Name),
		RpcServer: ServerConfig{
			Host: ctx.String(flags.RpcHostFlag.Name),
			Port: ctx.Int(flags.RpcPortFlag.Name),
		},
		MetricsServer: ServerConfig{
			Host: ctx.String(flags.MetricsHostFlag.Name),
			Port: ctx.Int(flags.MetricsPortFlag.Name),
		},
	}
}
