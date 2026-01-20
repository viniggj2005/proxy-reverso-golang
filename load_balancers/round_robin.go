package loadbalancers

import (
	"proxy-reverso-golang/structs"
	"sync/atomic"
)

type LoadBalancer interface {
	Next(servers []structs.ServerConfigStruct) *structs.Redirects
}

type RoundRobinBalancer struct {
	current uint64
}
type WeightedRoundRobinBalancer struct {
	targets []string
	current uint64
}

func NewRoundRobinBalancer() *RoundRobinBalancer {
	return &RoundRobinBalancer{}
}

func (receiver *RoundRobinBalancer) Next(servers []structs.ServerConfigStruct) *structs.Redirects {
	if len(servers) == 0 {
		return nil
	}
	number := atomic.AddUint64(&receiver.current, 1)
	for index := 0; index < len(servers); index++ {
		serverIndex := int((number + uint64(index) - 1) % uint64(len(servers)))
		if servers[serverIndex].Available {
			redirect := structs.Redirects{}
			redirect.Url = servers[serverIndex].Url
			return &redirect
		}
	}
	return nil
}

func NewWeightedRoundRobinBalancer(servers []structs.ServerConfigStruct) *WeightedRoundRobinBalancer {
	var targets []string
	for _, server := range servers {
		if server.Available == false {
			continue
		}
		weight := server.Weight
		if weight <= 0 {
			weight = 1
		}
		for index := 0; index < weight; index++ {
			targets = append(targets, server.Url)
		}
	}
	return &WeightedRoundRobinBalancer{targets: targets}
}

func (receiver *WeightedRoundRobinBalancer) Next(servers []structs.ServerConfigStruct) *structs.Redirects {
	if len(receiver.targets) == 0 {
		return nil
	}
	number := atomic.AddUint64(&receiver.current, 1)
	var redirect = structs.Redirects{}
	redirect.Url = receiver.targets[int((number-1)%uint64(len(receiver.targets)))]
	return &redirect
}
