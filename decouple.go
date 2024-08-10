package decouple

import (
	"github.com/saahilbdotdev/go-decouple/pkg/config"
	"github.com/saahilbdotdev/go-decouple/pkg/repository"
)

func Default(path ...string) *config.Config {
	config := new(config.AutoConfig)

	if len(path) > 0 {
		config.Init(path[0])
	} else {
		config.Init("")
	}

	var p string

	if config.SearchPath != "" {
		p = config.SearchPath
	} else {
		p = config.CallerPath()
	}

	if err := config.Load(p); err != nil {
		panic(err)
	}

	return &config.Config
}

func New(repository repository.Repository) *config.Config {
	config := new(config.Config)
	config.Init(repository)

	return config
}
