package config

import (
	"fmt"

	"github.com/knadh/koanf/parsers/yaml"
	"github.com/knadh/koanf/providers/file"
	"github.com/knadh/koanf/v2"
)

// Global koanf instance. Use "." as the key path delimiter. This can be "/" or any character.
var k = koanf.New(".")

func InitConfig() *koanf.Koanf {

	// Load YAML config and merge into the previously loaded config (because we can).
	err := k.Load(file.Provider("./config/config.yaml"), yaml.Parser())
	if err != nil {
		fmt.Println("Error occured", err)
		return nil
	}
	return k

}
