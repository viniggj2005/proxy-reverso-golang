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

func (r *RoundRobinBalancer) Next(servers []structs.ServerConfigStruct) *structs.Redirects {
	if len(servers) == 0 {
		return nil
	}
	n := atomic.AddUint64(&r.current, 1)
	for i := 0; i < len(servers); i++ {
		idx := int((n + uint64(i) - 1) % uint64(len(servers)))
		if servers[idx].Available {
			var redirect = structs.Redirects{}
			redirect.Url = servers[idx].Url
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
		for i := 0; i < weight; i++ {
			targets = append(targets, server.Url)
		}
	}
	return &WeightedRoundRobinBalancer{targets: targets}
}

func (w *WeightedRoundRobinBalancer) Next(servers []structs.ServerConfigStruct) *structs.Redirects {
	if len(w.targets) == 0 {
		return nil
	}
	n := atomic.AddUint64(&w.current, 1)
	var redirect = structs.Redirects{}
	redirect.Url = w.targets[int((n-1)%uint64(len(w.targets)))]
	return &redirect
}
