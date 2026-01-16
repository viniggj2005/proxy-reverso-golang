package functions

import (
	"encoding/json"
	"fmt"
	"os"
	"proxy-reverso-golang/global"
	"proxy-reverso-golang/structs"
)

func GetConfig() {
	fmt.Println("entrei na função")
	dir, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("\033[31mErro ao obter diretório de configuração:\033[0m", err)
		return
	}
	info, err := os.ReadDir(fmt.Sprintf("%s/teste-proxy", dir))
	if err != nil {
		fmt.Println("\033[31mErro ao ler diretório de configuração:\033[0m", err)
		return
	}
	for _, entry := range info {
		if entry.IsDir() {
			fmt.Println("\033[32mDiretório: \033[0m", entry.Name())
			getProxiesconfigs(entry.Name(), dir)
		} else {
			GetMainConfig(entry.Name(), dir)
		}
	}
}

func GetMainConfig(fileName string, dir string) (structs.ConfigStruct, error) {
	var config structs.ConfigStruct
	err := openFileAndGetContent(fmt.Sprintf("%s/teste-proxy/%s", dir, fileName), &config)
	if err != nil {
		fmt.Println("\033[31mErro ao ler config:\033[0m", err)
	}
	fmt.Println("\033[32mConfig lida main config:\033[0m", config)
	return config, nil
}

func getProxiesconfigs(fileName string, dir string) {
	files, _ := os.ReadDir(fmt.Sprintf("%s/teste-proxy/%s", dir, fileName))
	for _, file := range files {
		var config structs.ProxyConfigStruct

		err := openFileAndGetContent(fmt.Sprintf("%s/teste-proxy/%s/%s", dir, fileName, file.Name()), &config)
		if err != nil {
			fmt.Println("\033[31mErro ao ler config:\033[0m", err)
		}
		fmt.Println("\033[32mConfig lida proxies config:\033[0m", config)
		appendProxyConfigIfNotExists(config)
	}
}

func openFileAndGetContent(filePath string, target interface{}) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("\033[31mErro ao abrir arquivo:\033[0m", err)
	}
	return json.Unmarshal(content, target)
}

func appendProxyConfigIfNotExists(config structs.ProxyConfigStruct) {
	global.ProxyMutex.Lock()
	defer global.ProxyMutex.Unlock()

	for _, proxy := range global.ProxiesConfig.Proxies {
		if proxy.Prefix == config.Prefix {
			return
		}
	}
	global.ProxiesConfig.Proxies = append(global.ProxiesConfig.Proxies, config)
}
