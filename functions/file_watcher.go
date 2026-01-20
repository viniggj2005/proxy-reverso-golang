package functions

import (
	"fmt"
	"os"
	"time"

	"github.com/fsnotify/fsnotify"
)

func WatchConfigs() {
	directory, err := os.UserConfigDir()
	if err != nil {
		fmt.Println("\033[31mErro ao obter diretório de configuração:\033[0m", err)
		return
	}

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		fmt.Println("Erro ao criar watcher:", err)
		return
	}
	defer watcher.Close()

	watcher.Add(fmt.Sprintf("%s/teste-proxy/proxy-config", directory))

	var timer *time.Timer

	for {
		select {
		case event, ok := <-watcher.Events:
			if !ok {
				return
			}
			if event.Op&fsnotify.Write == fsnotify.Write {
				if timer != nil {
					timer.Stop()
				}
				timer = time.AfterFunc(100*time.Millisecond, func() {
					fmt.Println("\033[34mAlteração de arquivos detectada! Recarregando...\033[0m")
					GetConfig()
				})
			}
		case err, ok := <-watcher.Errors:
			if !ok {
				return
			}
			fmt.Println("Erro no watcher:", err)
		}
	}
}
