package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"reflect"
	"strings"

	"gopkg.in/yaml.v2"
)

func main() {
	args := os.Args

	if len(args) < 3 {
		panic("invalid parameter")
	}

	var master map[string]interface{}

	destPath, _ := filepath.Abs(os.Args[1])

	masterPath, _ := filepath.Abs(os.Args[2])

	bs, err := ioutil.ReadFile(masterPath)
	if err != nil {
		panic(err)
	}

	if err := yaml.Unmarshal(bs, &master); err != nil {
		panic(err)
	}

	for _, arg := range args[3:] {
		var fragment map[string]interface{}

		path, _ := filepath.Abs(arg)

		bs, err = ioutil.ReadFile(path)
		if err != nil {
			panic(err)
		}

		if err := yaml.Unmarshal(bs, &fragment); err != nil {
			panic(err)
		}

		for k, v := range fragment {
			master[k] = v
		}
	}

	bs, err = yaml.Marshal(deleteUnEncrypted(master))
	if err != nil {
		panic(err)
	}

	if err := ioutil.WriteFile(destPath, bs, 0644); err != nil {
		panic(err)
	}
}

func deleteUnEncrypted(from map[string]interface{}) map[string]interface{} {
	result := make(map[string]interface{})

	for k, v := range from {
		key := k

		ref := reflect.ValueOf(v)

		if ref.Kind().String() == "map" {
			newMap := make(map[string]interface{})
			for _, key := range ref.MapKeys() {
				index := ref.MapIndex(key)
				newMap[key.Interface().(string)] = index.Interface()
			}
			result[key] = deleteUnEncrypted(newMap)
		} else {
			key = strings.Replace(key, "_unencrypted", "", 1)
			result[key] = v
		}
	}

	return result
}
