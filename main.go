package main

import (
	"aproc/lib"
	"log"
	"os"
	// v2 "github.com/containerd/cgroups/v2"
)

// const (
// 	defaultCgroup2Path = "/sys/fs/cgroup"
// )

// func deleteManager(manager *v2.Manager) {
// 	if err := manager.Delete(); err != nil {
// 		log.Fatalln(err)
// 	}
// }

func main() {
	// 检查是否是超级用户
	if os.Getenv("HOME") != "/root" {
		log.Fatalln("请使用超级用户运行此程序")
	}
	// 创建 rootManager cgroup
	// var (
	// 	rootManager *v2.Manager
	// 	err         error
	// )

	// if rootManager, err = v2.NewManager(defaultCgroup2Path, "aproc", &v2.Resources{}); err != nil {
	// 	log.Fatalln(err)
	// }

	// defer deleteManager(rootManager)
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
	// TODO: something else
	// watcher.Exit()
	watcher.WaitForExit()
	// 创建子组
	// var quota int64 = 10 * 1000

	// var period uint64 = 100 * 1000

	// zsh, err := rootManager.NewChild("zsh", &v2.Resources{
	// 	CPU: &v2.CPU{
	// 		Max: v2.NewCPUMax(&quota, &period),
	// 	},
	// })
	// if err != nil {
	// 	log.Fatalln(err)
	// }

	// defer deleteManager(zsh)
}
