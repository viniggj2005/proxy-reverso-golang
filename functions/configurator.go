package functions

import (
	"encoding/json"
	"fmt"
	"os"
	"proxy-reverso-golang/providers"
	"proxy-reverso-golang/structs"
)

func GetConfig() {
	directory, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("\033[31mErro ao obter diretório de configuração:\033[0m", err)
		return
	}
	info, err := os.ReadDir(fmt.Sprintf("%s/teste-proxy", directory))
	if err != nil {
		fmt.Println("\033[31mErro ao ler diretório de configuração:\033[0m", err)
		return
	}
	for _, entry := range info {
		if entry.IsDir() {
			getProxiesconfigs(entry.Name(), directory)
		}
	}
}

func GetMainConfig(fileName string, directory string) (ConfigStruct, error) {
	if directory == "" {
		var err error
		directory, err = os.UserConfigDir()
		if err != nil {
			fmt.Println("\033[31mErro ao obter diretório de configuração:\033[0m", err)
			return ConfigStruct{}, err
		}
	}
	var config ConfigStruct
	err := openFileAndGetContent(fmt.Sprintf("%s/teste-proxy/%s", directory, fileName), &config)
	if err != nil {
		fmt.Println("\033[31mErro ao ler config:\033[0m", err)
		return ConfigStruct{}, err
	}
	return config, nil
}

func getProxiesconfigs(fileName string, directory string) {
	var tempProxies []structs.ProxyConfigStruct
	files, _ := os.ReadDir(fmt.Sprintf("%s/teste-proxy/%s", directory, fileName))
	for _, file := range files {
		var config structs.ProxyConfigStruct

		err := openFileAndGetContent(fmt.Sprintf("%s/teste-proxy/%s/%s", directory, fileName, file.Name()), &config)
		if err != nil {
			fmt.Println("\033[31mErro ao ler config:\033[0m", err)
		}
		tempProxies = append(tempProxies, config)
	}

	providers.SetProxies(tempProxies)

	for _, proxy := range tempProxies {
		providers.DeleteLoadBalancer(proxy.Prefix)
	}

}

func openFileAndGetContent(filePath string, target interface{}) error {
	content, err := os.ReadFile(filePath)
	if err != nil {
		fmt.Println("\033[31mErro ao abrir arquivo:\033[0m", err)
		return err
	}
	return json.Unmarshal(content, target)
}
