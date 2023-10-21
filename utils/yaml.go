package utils

import (
	"os"

	"gopkg.in/yaml.v3"
)

type Admin struct {
	Name       string
	CidrBlocks []string
}

func ReadYaml() {
	file, err := os.ReadFile("../cidrs.yaml")
	Check(err)

	data := make(map[string][]string)

	err = yaml.Unmarshal(file, &data)

	admins := make([]Admin, 0)
	for k, v := range data {
		admin := Admin{
			Name:       k,
			CidrBlocks: v,
		}
		admins = append(admins, admin)
	}
}
