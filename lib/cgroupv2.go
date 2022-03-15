package lib

import (
	"container/list"

	v2 "github.com/containerd/cgroups/v2"
	"github.com/sirupsen/logrus"
)

var (
	rootManager *v2.Manager
	managers    list.List
)

func getRootManager() (*v2.Manager, error) {
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

func CreateManager(pid uint64, name string, resources *v2.Resources) (*v2.Manager, error) {
	rootManager, err := getRootManager()
	if err != nil {
		return nil, err
	}

	subManager, err := rootManager.NewChild(name, resources)
	if err != nil {
		return nil, err
	}

	managers.PushBack(subManager)

	return subManager, nil
}

func CleanManager() {
	for elment := managers.Front(); elment != nil; elment = elment.Next() {
		if manager, ok := elment.Value.(*v2.Manager); ok {
			res, err := manager.Procs(true)
			if err != nil {
				continue
			}

			if len(res) == 0 {
				if err := manager.Delete(); err != nil {
					logrus.Errorln(err)
				}

				managers.Remove(elment)
			}
		}
	}
}
