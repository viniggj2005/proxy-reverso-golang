package global

import (
	"proxy-reverso-golang/structs"
	"sync"
)

var ProxiesConfig structs.AllProxiesConfigStruct

var ProxyMutex sync.RWMutex
