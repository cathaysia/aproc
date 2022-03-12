package main

import (
	"aproc/lib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

func main() {
	// 检查是否是超级用户
	if os.Getenv("HOME") != "/root" {
		log.Fatalln("请使用超级用户运行此程序")
	}

	rootManager, err := lib.GetInstanceOfRootManager()
	if err != nil {
		log.Fatalln(err)
	}

	defer lib.DeleteManager(rootManager)
	// 监控 /proc 目录的变动
	watcher := lib.NewProgressWatcher(2)

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if event.IsCreate() {
					log.Printf("%v is Created\n", event.PID)
				} else if event.IsDelete() {
					log.Printf("%v is Deleted\n", event.PID)
				}
			case err := <-watcher.Error:
				log.Println(err)
				watcher.Exit()

				break
			}
		}
	}()

	if err := watcher.Watch(); err != nil {
		log.Println(err)
	}

	//
	sigInt := make(chan os.Signal, 2)
	signal.Notify(sigInt, syscall.SIGINT)

	if <-sigInt == syscall.SIGINT {
		watcher.Exit()
		<-watcher.WaitForExit

		return
	}
}
