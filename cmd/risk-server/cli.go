package main

import (
	"context"

	"github.com/urfave/cli/v2"

	"github.com/ethereum/go-ethereum/log"
	"github.com/the-web3/mock-risk-server/common/cliapp"
	"github.com/the-web3/mock-risk-server/config"
	flags2 "github.com/the-web3/mock-risk-server/flags"
	"github.com/the-web3/mock-risk-server/services"
)

func runRpc(ctx *cli.Context, shutdown context.CancelCauseFunc) (cliapp.Lifecycle, error) {
	cfg, err := config.LoadConfig(ctx)
	if err != nil {
		log.Error("Failed to load config", "err", err)
		return nil, err
	}
	ristServerConfig := &services.RiskServerConfig{
		GrpcHostname: cfg.RpcServer.Host,
		GrpcPort:     cfg.RpcServer.Port,
	}
	return services.NewRiskServerWireServices(ristServerConfig)
}

func NewCli() *cli.App {
	flags := flags2.Flags
	return &cli.App{
		Version:              "v0.0.1-beta",
		Description:          "wallet mock risk server",
		EnableBashCompletion: true,
		Commands: []*cli.Command{
			{
				Name:        "rpc",
				Flags:       flags,
				Description: "Run rpc services",
				Action:      cliapp.LifecycleCmd(runRpc),
			},
			{
				Name:        "version",
				Description: "Show project version",
				Action: func(ctx *cli.Context) error {
					cli.ShowVersion(ctx)
					return nil
				},
			},
		},
	}
}
