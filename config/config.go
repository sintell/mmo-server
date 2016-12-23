package config

import (
	"fmt"
	"io/ioutil"
)

func ReadConfig(path string) []byte {
	fmt.Printf("Reading file %s\n", path)
	str, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err.Error())
	}
	return str
}
