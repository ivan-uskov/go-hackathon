package transport

import (
	log "github.com/sirupsen/logrus"
	"google.golang.org/grpc"
)

func DialGRPC(address string) *grpc.ClientConn {
	conn, err := grpc.Dial(address, grpc.WithInsecure())
	if err != nil {
		log.Fatal(err)
	}

	return conn
}
