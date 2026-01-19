package loadbalancers

import (
	"proxy-reverso-golang/structs"
	"sync/atomic"
)

type RoundRobinBalancer struct {
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
	var redirect = structs.Redirects{}
	redirect.Url = servers[(int(n)-1)%len(servers)].Url
	return &redirect
}
