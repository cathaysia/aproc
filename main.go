package main

import (
	"aproc/lib"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	settings []lib.Settings
	err      error
)

// 根据
func createCGroupForPID(pid uint64) error {
	procName, err := lib.GetProgressNameByPID(pid)
	if err != nil {
		return err
	}

	for _, v := range settings {
		if v.Proc == procName {
			// 这里要检查是否已经存在了
			manager, err := lib.CreateManager(pid, procName, &v.Resources)
			if err != nil {
				return err
			}

			err = manager.AddProc(pid)
			if err != nil {
				return err
			}
		}
	}

	return nil
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

	sigInt := make(chan os.Signal, 2)
	signal.Notify(sigInt, syscall.SIGINT)
	// 监控 /proc 目录的变动
	watcher := lib.NewProgressWatcher(2)
	watchExit := make(chan bool)
	waitExit := make(chan bool)

	go func() {
		for {
			select {
			case event := <-watcher.Event:
				if event.IsCreate() {
					log.Printf("%v is Created\n", event.PID)

					if err := createCGroupForPID(event.PID); err != nil {
						watcher.Error <- err
					}
				} else if event.IsDelete() {
					log.Printf("%v is Deleted\n", event.PID)

					lib.CleanManager()
				}
			case err := <-watcher.Error:
				log.Println(err)
				signal.Notify(sigInt, syscall.SIGINT) // 请求退出进程
			case <-watchExit: // 等待退出协程
				waitExit <- true

				return
			}

		}
	}()

	if err := watcher.Watch(); err != nil {
		log.Println(err)
	}

	//

	if <-sigInt == syscall.SIGINT {
		watcher.Exit()        // 通知 watcher 退出
		<-watcher.WaitForExit // 等待 watcher 退出
		watchExit <- true     // 通知轮询协程退出
		<-waitExit            // 等待轮询协程退出

		return
	}
}
