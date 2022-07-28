package manager

import (
	"net"
	"strings"
)

type ServerManager interface {
	RegisterToCenter(serviceName string) error
}

func GetSelfIPv4Addresses() ([]string, error) {
	result := []string{}
	addresses, err := net.InterfaceAddrs()
	if err != nil {
		return result, err
	}
	for _, address := range addresses {
		host := strings.Split(address.String(), "/")[0]
		i := net.ParseIP(host).To4()
		if i != nil && !i.IsLoopback() {
			result = append(result, i.String())
		}
	}
	return result, nil
}
