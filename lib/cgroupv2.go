package lib

import (
	"log"

	v2 "github.com/containerd/cgroups/v2"
)

func DeleteManager(manager *v2.Manager) {
	if err := manager.Delete(); err != nil {
		log.Fatalln(err)
	}
}

var (
	rootManager *v2.Manager
	period      uint64
)

func GetInstanceOfRootManager() (*v2.Manager, error) {
	if rootManager != nil {
		return rootManager, nil
	}

	// v2.NewManager 的第二个参数必须以 / 开头
	rootManager, err := v2.NewManager("/sys/fs/cgroup", "/aproc", &v2.Resources{})
	if err != nil {
		return nil, err
	}

	return rootManager, nil
}

func init() {
	period = 100 * 1000
}

func CreateManager(pid int64, resources *v2.Resources) (*v2.Manager, error) {
	rootManager, err := GetInstanceOfRootManager()
	if err != nil {
		return nil, err
	}

	// TODO: get a name by pid
	var name string

	subManager, err := rootManager.NewChild(name, resources)
	if err != nil {
		return nil, err
	}

	return subManager, nil
}
