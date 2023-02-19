package grpcapi

import (
	"fmt"
	"log"
	"net"

	"github.com/3rubasa/shagent/pkg/businesslogic"
	grpc "google.golang.org/grpc"
)

type API struct {
	mc   businesslogic.MainController
	impl *impl
	port int
	s    *grpc.Server
}

func New(mc businesslogic.MainController, port int) *API {
	return &API{
		mc:   mc,
		impl: NewImpl(mc),
		port: port,
	}
}

func (a *API) Start() error {
	a.s = grpc.NewServer()
	RegisterStateProviderServer(a.s, a.impl)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", a.port))
	if err != nil {
		log.Println("ERROR: Failed to start GRPC API: ", err)
		return err
	}

	go func() {
		log.Println("NOTICE: Starting GRPC API server to listen at ", lis.Addr())

		err := a.s.Serve(lis)
		if err != nil {
			log.Println("ERROR: GPC API: Failed to serve: ", err)
			return
		}

		log.Println("NOTICE: GRPC API server has been shut down")
	}()

	return nil
}

func (a *API) Stop() error {
	a.s.GracefulStop()
	return nil
}
