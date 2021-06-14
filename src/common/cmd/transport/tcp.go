package transport

import (
	"fmt"
	log "github.com/sirupsen/logrus"
	"net"
)

func ListenTCP(port int) net.Listener {
	l, err := net.Listen("tcp", fmt.Sprintf(":%v", port))
	if err != nil {
		log.Fatal(err)
	}

	return l
}
