package main

import (
	"context"

	"github.com/urfave/cli/v2"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"

	"github.com/ethereum/go-ethereum/log"

	"github.com/dapplink-external/risk-server/client/walletapiclient"
	"github.com/dapplink-external/risk-server/common/cliapp"
	"github.com/dapplink-external/risk-server/config"
	flags2 "github.com/dapplink-external/risk-server/flags"
	"github.com/dapplink-external/risk-server/protobuf/walletapi"
	"github.com/dapplink-external/risk-server/services"
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
		LevelDBPath:  cfg.LevelDbPath,
		AccessToken:  cfg.AccessToken,
	}

	connApi, err := grpc.NewClient(cfg.ApiGateWayRpc, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		log.Error("Connect to da retriever fail", "err", err)
		return nil, err
	}
	gateWayClient := walletapi.NewWalletApiGateWayServiceClient(connApi)

	apiGateWayClient, err := walletapiclient.NewWalletApiGateWayServiceClient(context.Background(), gateWayClient, "DappLinkEthereum")
	if err != nil {
		log.Error("Connect to da retriever fail", "err", err)
		return nil, err
	}

	return services.NewRiskServerWireServices(ristServerConfig, apiGateWayClient)
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
