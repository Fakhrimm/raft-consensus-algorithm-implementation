package node

import (
	"log"
	"net"
	"strings"
)

func parseAddresses(clusterAddresses string) []net.TCPAddr {
	var clusterAddressesTcp []net.TCPAddr

	if clusterAddresses != "" {

		addresses := strings.Split(clusterAddresses, ",")

		for _, address := range addresses {
			a, err := net.ResolveTCPAddr("tcp", address)
			if err != nil {
				log.Fatalf("Invalid address: %v, error: %v", address, err.Error())
			}

			clusterAddressesTcp = append(clusterAddressesTcp, *a)
		}
	}

	return clusterAddressesTcp
}
