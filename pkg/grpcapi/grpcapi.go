package grpcapi

import (
	"fmt"
	"log"
	"net"

	"github.com/3rubasa/shagent/pkg/businesslogic"
	grpc "google.golang.org/grpc"
)

type API struct {
	bl   *businesslogic.BusinessLogic
	impl *impl
	port int
	s    *grpc.Server
}

func New(bl *businesslogic.BusinessLogic, port int) *API {
	return &API{
		bl:   bl,
		impl: NewImpl(bl),
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
