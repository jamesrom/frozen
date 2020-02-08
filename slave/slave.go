package main

import (
	"context"
	"flag"
	"log"
	"net"

	"google.golang.org/grpc"

	"github.com/arr-ai/frozen/slave/proto/slave"
)

type slaveServer struct {
}

func (s *slaveServer) Union(_ context.Context, req *slave.UnionRequest) (*slave.Tree, error) {
	panic("unfinished")
	// a :=
	// result = result.Merge(i.Value().(frozen.Map), func(_, a, b interface{}) interface{} {
	// 		return a
	// 	})
	// }
	// return nil, fmt.Errorf("unfinished")
}

func main() {
	listen := flag.String("listen", "", "[host]:port to listen on")
	flag.Parse()

	skt, err := net.Listen("tcp", *listen)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	grpcServer := grpc.NewServer()
	slave.RegisterSlaveServer(grpcServer, &slaveServer{})
	// determine whether to use TLS
	panic(grpcServer.Serve(skt))
}
