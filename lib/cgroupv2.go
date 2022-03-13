package lib

import (
	"container/list"
	"log"
	"runtime"

	v2 "github.com/containerd/cgroups/v2"
)

func DeleteManager(manager *v2.Manager) {
	if err := manager.Delete(); err != nil {
		log.Fatalln(err)
	}
}

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

	runtime.SetFinalizer(rootManager, DeleteManager)

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
	runtime.SetFinalizer(subManager, DeleteManager)

	return subManager, nil
}

func CleanManager() {
	for elment := managers.Front(); elment != nil; elment = elment.Next() {
		if m, ok := elment.Value.(*v2.Manager); ok {
			res, err := m.Procs(true)
			if err != nil {
				continue
			}

			if len(res) == 0 {
				managers.Remove(elment)
			}
		}
	}
}
