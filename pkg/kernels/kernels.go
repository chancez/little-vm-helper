package kernels

import (
	"encoding/json"
	"fmt"
	"os"
	"path"
)

var ConfigFname = "conf.json"

// Initalizes a new directory (it will create it if it does not exist).
// the provided conf will be saved in the directory.
// if conf is nil, an empty configuration will be used.
func InitDir(dir string, conf *Conf) error {
	err := os.MkdirAll(dir, 0755)
	if err != nil && !os.IsExist(err) {
		return fmt.Errorf("failed to create directory '%s': %w", dir, err)
	}

	confFname := path.Join(dir, ConfigFname)
	if _, err := os.Stat(confFname); err == nil {
		return fmt.Errorf("config file `%s` already exists", dir)
	}

	if conf == nil {
		conf = &Conf{
			Kernels: make([]KernelConf, 0),
		}
	}

	return conf.SaveTo(dir)
}

// Load configuration from a directory
func LoadDir(dir string) (*KernelsDir, error) {
	data, err := os.ReadFile(path.Join(dir, ConfigFname))
	if err != nil {
		return nil, err
	}

	ks := KernelsDir{Dir: dir}
	err = json.Unmarshal(data, &ks.Conf)
	if err != nil {
		return nil, err
	}
	return &ks, nil
}

func AddKernel(dir string, cnf *KernelConf) error {
	kd, err := LoadDir(dir)
	if err != nil {
		return err
	}

	if kd.KernelConfig(cnf.Name) != nil {
		return fmt.Errorf("kernel `%s` already exist", cnf.Name)
	}

	kd.Conf.Kernels = append(kd.Conf.Kernels, *cnf)
	return kd.Conf.SaveTo(dir)
}

func BuildKernel(dir string, name string) error {
	return nil
}
