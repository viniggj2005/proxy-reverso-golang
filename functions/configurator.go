package functions

import (
	"encoding/json"
	"fmt"
	"os"
	"proxy-reverso-golang/global"
	"proxy-reverso-golang/structs"
)

func GetConfig() {
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
			fmt.Println("\033[32mArquivo: \033[0m", entry.Name())
		}
	}
}

func GetMainConfig(fileName string, dir string) (ConfigStruct, error) {
	if dir == "" {
		var err error
		dir, err = os.UserConfigDir()
		if err != nil {
			fmt.Println("\033[31mErro ao obter diretório de configuração:\033[0m", err)
			return ConfigStruct{}, err
		}
	}
	var config ConfigStruct
	err := openFileAndGetContent(fmt.Sprintf("%s/teste-proxy/%s", dir, fileName), &config)
	if err != nil {
		fmt.Println("\033[31mErro ao ler config:\033[0m", err)
		return ConfigStruct{}, err
	}
	fmt.Printf("\033[32mConfig lida main config:\033[0m %+v\n", config)
	return config, nil
}

func getProxiesconfigs(fileName string, dir string) {
	var tempProxies []structs.ProxyConfigStruct
	files, _ := os.ReadDir(fmt.Sprintf("%s/teste-proxy/%s", dir, fileName))
	for _, file := range files {
		var config structs.ProxyConfigStruct

		err := openFileAndGetContent(fmt.Sprintf("%s/teste-proxy/%s/%s", dir, fileName, file.Name()), &config)
		if err != nil {
			fmt.Println("\033[31mErro ao ler config:\033[0m", err)
		}
		fmt.Printf("\033[32mConfig lida proxies config:\033[0m %+v\n", config)
		tempProxies = append(tempProxies, config)
	}
	global.ProxyMutex.Lock()
	global.ProxiesConfig.Proxies = tempProxies
	global.ProxyMutex.Unlock()
}

func openFileAndGetContent(filePath string, target interface{}) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("\033[31mErro ao abrir arquivo:\033[0m", err)
	}
	return json.Unmarshal(content, target)
}
