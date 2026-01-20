package loadbalancers

import (
	"math/rand"
	"proxy-reverso-golang/structs"
)

type RandomBalancer struct {
}

func NewRandomBalancer() *RandomBalancer {
	return &RandomBalancer{}
}

func (receiver *RandomBalancer) Next(servers []structs.ServerConfigStruct) *structs.Redirects {
	if len(servers) == 0 {
		return nil
	}
	number := uint64(rand.Intn(len(servers)))
	for index := 0; index < len(servers); index++ {
		serverIndex := int((number + uint64(index)) % uint64(len(servers)))
		if servers[serverIndex].Available {
			redirect := structs.Redirects{}
			redirect.Url = servers[serverIndex].Url
			return &redirect
		}
	}
	return nil
}
