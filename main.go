package main

import (
	"aproc/lib"
	"log"
	"os"
	"os/signal"
	"syscall"

	v2 "github.com/containerd/cgroups/v2"
)

var (
	rootManager *v2.Manager
	settings    []lib.Settings
	err         error
)

// 根据
func createCGroupForPID(pid uint64) error {
	rootManager, err = lib.GetInstanceOfRootManager()
	if err != nil {
		return err
	}

	procName, err := lib.GetProgressNameByPID(pid)
	if err != nil {
		return err
	}

	for _, v := range settings {
		if v.Proc == procName {
			// 这里要检查时候已经存在了
			manager, err := rootManager.NewChild(procName, &v.Resources)
			if err != nil {
				return err
			}

			err = manager.AddProc(pid)
		}
	}

	return err
}

func main() {
	// 检查是否是超级用户
	if os.Getenv("HOME") != "/root" {
		log.Fatalln("请使用超级用户运行此程序")
	}

	settings, err = lib.GetSettings()
	if err != nil {
		log.Fatalln(err)
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
					createCGroupForPID(uint64(event.PID))
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
