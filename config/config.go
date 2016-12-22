package config

import (
	"fmt"
	"io/ioutil"
)

func ReadConfig(path string) []byte {
	str, err := ioutil.ReadFile(path)
	if err != nil {
		fmt.Errorf("can't read file %s: %s", path, err.Error())
	}
	return str
}
