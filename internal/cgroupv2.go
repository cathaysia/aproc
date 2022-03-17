package internal

import (
	"container/list"
	"unsafe"

	v2 "github.com/containerd/cgroups/v2"
	"github.com/sirupsen/logrus"
)

//go:linkname setResources github.com/containerd/cgroups/v2.setResources
func setResources(path string, resources *v2.Resources) error

type Manager struct { // v2.Manager for access private path
	_    string
	path string
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

func ReloadManager(settings []Settings) error {
	for element := managers.Front(); element != nil; element = element.Next() {
		manager, ok := element.Value.(*v2.Manager)
		if !ok {
			continue
		}

		procs, err := manager.Procs(false)
		if err != nil || len(procs) == 0 {
			continue
		}

		name, err := GetProgressNameByPID(procs[0])
		if err != nil {
			continue
		}

		for _, v := range settings {
			if v.Proc == name {
				p := *(*Manager)(unsafe.Pointer(manager))
				if err := setResources(p.path, &v.Resources); err != nil {
					if err := manager.Delete(); err != nil {
						return err
					}

					return err
				}
			}
		}
	}

	return nil
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
