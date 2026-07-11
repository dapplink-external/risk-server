package flags

import (
	"github.com/urfave/cli/v2"
)

const envVarPrefix = "RISK"

func prefixEnvVars(name string) []string {
	return []string{envVarPrefix + "_" + name}
}

var (
	RpcHostFlag = &cli.StringFlag{
		Name:     "rpc-host",
		Usage:    "The host of the rpc",
		EnvVars:  prefixEnvVars("RPC_HOST"),
		Required: true,
	}
	RpcPortFlag = &cli.IntFlag{
		Name:     "rpc-port",
		Usage:    "The port of the rpc",
		EnvVars:  prefixEnvVars("RPC_PORT"),
		Value:    8987,
		Required: true,
	}

	MetricsHostFlag = &cli.StringFlag{
		Name:     "metrics-host",
		Usage:    "The host of the metrics",
		EnvVars:  prefixEnvVars("METRICS_HOST"),
		Required: true,
	}
	MetricsPortFlag = &cli.IntFlag{
		Name:     "metrics-port",
		Usage:    "The port of the metrics",
		EnvVars:  prefixEnvVars("METRICS_PORT"),
		Value:    7214,
		Required: true,
	}

	ApiGateWayRpcFlag = &cli.StringFlag{
		Name:     "api-gateway-rpc",
		Usage:    "The gateway rpc of the api",
		EnvVars:  prefixEnvVars("API_GATEWAY_RPC"),
		Required: true,
	}

	LevelDbPathFlag = &cli.StringFlag{
		Name:     "leveldb-path",
		Usage:    "The path of leveldb ",
		EnvVars:  prefixEnvVars("LEVELDB_PATH"),
		Required: true,
	}
	RiskKeyFlag = &cli.StringFlag{
		Name:     "risk-key",
		Usage:    "The risk key of the risk controller",
		EnvVars:  prefixEnvVars("RISK_KEY"),
		Required: true,
	}

	AccessTokenFlag = &cli.StringFlag{
		Name:     "access-token",
		Usage:    "The access token risk controller",
		EnvVars:  prefixEnvVars("ACCESS_TOKEN"),
		Required: true,
	}
)

var requireFlags = []cli.Flag{
	RpcHostFlag,
	RpcPortFlag,
	MetricsPortFlag,
	MetricsHostFlag,

	ApiGateWayRpcFlag,
	LevelDbPathFlag,
	RiskKeyFlag,
	AccessTokenFlag,
}

var optionalFlags = []cli.Flag{}

var Flags []cli.Flag

func init() {
	Flags = append(requireFlags, optionalFlags...)
}
