package internal

import (
	"encoding/json"
	"io/ioutil"
	"os"

	v2 "github.com/containerd/cgroups/v2"
)

type Settings struct {
	Proc      string
	Resources v2.Resources
}

// 从 /etc/aproc/settings.json 获取数据
func GetSettings() ([]Settings, error) {
	if _, err := os.Stat("/etc/aproc"); err != nil {
		if err = os.MkdirAll("/etc/aproc", os.ModePerm); err != nil {
			return nil, err
		}
	}
	// 读取 settingsFile 文件
	settingsFile, err := os.OpenFile("/etc/aproc/settings.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return nil, err
	}

	defer settingsFile.Close()

	settingsData, err := ioutil.ReadAll(settingsFile)
	if err != nil {
		return nil, err
	}

	var settings []Settings

	if err = json.Unmarshal(settingsData, &settings); err != nil {
		return nil, err
	}

	return settings, nil
}

func CreateEmptySettings() error {
	// create file
	if _, err := os.Stat("/etc/aproc"); err != nil {
		if err = os.MkdirAll("/etc/aproc", os.ModePerm); err != nil {
			return err
		}
	}

	settingsFile, err := os.OpenFile("/etc/aproc/settings.json", os.O_RDWR|os.O_CREATE, os.ModePerm)
	if err != nil {
		return err
	}

	settingsFile.Close()

	settings := []Settings{
		{
			Proc: "systemd",
		},
	}

	data, err := json.Marshal(settings)
	if err != nil {
		return err
	}

	return ioutil.WriteFile("/etc/aproc/settings.json", data, os.ModePerm)
}
