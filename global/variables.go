package global

import (
	loadbalancers "proxy-reverso-golang/load_balancers"
	"proxy-reverso-golang/structs"
	"sync"
)

var ProxiesConfig structs.AllProxiesConfigStruct

var ProxyMutex sync.RWMutex

var LoadBalancers = make(map[string]loadbalancers.LoadBalancer)
var BalancerMutex sync.Mutex
