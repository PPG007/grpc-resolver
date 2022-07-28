package tests

import (
	"log"
	"nameresolver/server/manager"
	"net"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGetIp(t *testing.T) {
	addresses, err := net.InterfaceAddrs()
	assert.NoError(t, err)
	for _, address := range addresses {
		log.Println(address.Network(), address.String(), net.ParseIP(strings.Split(address.String(), "/")[0]).IsLoopback())
	}

}

func TestGetOneIPv4(t *testing.T) {
	addresses, err := manager.GetSelfIPv4Addresses()
	assert.NoError(t, err)
	log.Printf("%+v\n", addresses)
}
