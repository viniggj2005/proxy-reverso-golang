package providers

import (
	loadbalancers "proxy-reverso-golang/load_balancers"
	"proxy-reverso-golang/structs"
	"sync"
)

var proxyMutex sync.RWMutex
var balancerMutex sync.Mutex
var proxiesConfig structs.AllProxiesConfigStruct
var loadBalancers = make(map[string]loadbalancers.LoadBalancer)

func DeleteLoadBalancer(prefix string) {
	balancerMutex.Lock()
	defer balancerMutex.Unlock()
	delete(loadBalancers, prefix)
}

func GetProxies() []structs.ProxyConfigStruct {
	proxyMutex.RLock()
	defer proxyMutex.RUnlock()
	return proxiesConfig.Proxies
}

func SetProxies(proxies []structs.ProxyConfigStruct) {
	proxyMutex.Lock()
	defer proxyMutex.Unlock()
	proxiesConfig.Proxies = proxies
}

func SetProxyAvailability(availability bool, proxyIndex int, serverIndex int) bool {
	proxyMutex.Lock()
	defer proxyMutex.Unlock()
	if proxyIndex >= 0 && proxyIndex < len(proxiesConfig.Proxies) {
		if serverIndex >= 0 && serverIndex < len(proxiesConfig.Proxies[proxyIndex].Servers) {
			proxiesConfig.Proxies[proxyIndex].Servers[serverIndex].Available = availability
			return true
		}
	}
	return false
}

func GetOrCreateBalancer(prefix string, factory func() loadbalancers.LoadBalancer) loadbalancers.LoadBalancer {
	balancerMutex.Lock()
	defer balancerMutex.Unlock()

	if balancer, exists := loadBalancers[prefix]; exists {
		return balancer
	}

	balancer := factory()
	loadBalancers[prefix] = balancer
	return balancer
}
