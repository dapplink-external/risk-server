package services

import (
	"context"
	"fmt"
	"net"
	"sync/atomic"

	"google.golang.org/grpc"
	"google.golang.org/grpc/reflection"

	"github.com/ethereum/go-ethereum/log"
	"github.com/the-web3/mock-risk-server/client/walletapiclient"
	"github.com/the-web3/mock-risk-server/protobuf/riskcontroller"
)

const MaxRecvMessageSize = 1024 * 1024 * 300

type RiskServerConfig struct {
	GrpcHostname string
	GrpcPort     int
}

type RiskServerWireServices struct {
	*RiskServerConfig
	rpcApiClient *walletapiclient.WalletApiGateWayServiceClient
	stopped      atomic.Bool
}

func NewRiskServerWireServices(config *RiskServerConfig, rpcApiClient *walletapiclient.WalletApiGateWayServiceClient) (*RiskServerWireServices, error) {
	return &RiskServerWireServices{
		RiskServerConfig: config,
		rpcApiClient:     rpcApiClient,
	}, nil
}

func (rss *RiskServerWireServices) Start(ctx context.Context) error {
	go func(rss *RiskServerWireServices) {
		addr := fmt.Sprintf("%s:%d", rss.GrpcHostname, rss.GrpcPort)
		log.Info("start rpc server", "addr", addr)
		listener, err := net.Listen("tcp", addr)
		if err != nil {
			log.Error("Could not start tcp listener. ")
		}
		gs := grpc.NewServer(
			grpc.MaxRecvMsgSize(MaxRecvMessageSize),
			grpc.ChainUnaryInterceptor(
				nil,
			),
		)

		reflection.Register(gs)

		riskcontroller.RegisterRiskControllerServicesServer(gs, rss)

		log.Info("grpc info", "port", rss.GrpcPort, "address", listener.Addr())

		if err := gs.Serve(listener); err != nil {
			log.Error("Could not GRPC server")
		}
	}(rss)
	return nil
}

func (rss *RiskServerWireServices) Stop(ctx context.Context) error {
	rss.stopped.Store(true)
	return nil
}

func (rss *RiskServerWireServices) Stopped() bool {
	return rss.stopped.Load()
}
