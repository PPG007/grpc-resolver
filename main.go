package main

import (
	"log"
	"nameresolver/server/manager"
)

func main() {
	addresses, err := manager.GetSelfIPv4Addresses()
	if err != nil {
		panic(err.Error())
	}
	for _, address := range addresses {
		log.Println(address)
	}
}
