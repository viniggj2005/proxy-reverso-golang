package loadbalancers

import (
	"fmt"
	"math/rand"
	"proxy-reverso-golang/structs"
	"time"
)

type RandomBalancer struct {
	current uint64
}

func NewRandomBalancer() *RandomBalancer {
	rand.Seed(time.Now().UnixNano())
	return &RandomBalancer{}
}

func (r *RandomBalancer) Next(servers []structs.ServerConfigStruct) *structs.Redirects {
	if len(servers) == 0 {
		return nil
	}
	number := uint64(rand.Intn(len(servers)))
	for i := 0; i < len(servers); i++ {
		idx := int((number + uint64(i)) % uint64(len(servers)))
		if servers[idx].Available {
			var redirect = structs.Redirects{}
			redirect.Url = servers[idx].Url
			fmt.Println("RandomBalancer: ", servers[idx].Url)
			return &redirect
		}
	}
	return nil
}
