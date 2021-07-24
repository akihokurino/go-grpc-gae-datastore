package yaml

import (
	"io/ioutil"
	"os"

	"gopkg.in/yaml.v2"
)

type env struct {
	Variables map[string]string `yaml:"env_variables"`
}

func MustLoadLocalEnv(path string) {
	if os.Getenv("IS_LOCAL") != "true" {
		return
	}

	buf, err := ioutil.ReadFile(path)
	if err != nil {
		panic(err)
	}

	e := env{}
	if err := yaml.Unmarshal(buf, &e); err != nil {
		panic(err)
	}

	for k, v := range e.Variables {
		if err := os.Setenv(k, v); err != nil {
			panic(err)
		}
	}
}
