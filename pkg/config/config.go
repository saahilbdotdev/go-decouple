package config

import (
	"os"
	"path/filepath"
	"reflect"
	"runtime"

	"github.com/saahilbdotdev/go-decouple/pkg/repository"
	"github.com/saahilbdotdev/go-decouple/pkg/utils"
)

type Config struct {
	Repository repository.Repository
}

func (c *Config) Init(repository repository.Repository) {
	c.Repository = repository
}

func (c *Config) CastBoolean(value string) bool {
	if value == "" {
		return false
	}

	boolValue, err := utils.StringToBool(value)
	if err != nil {
		panic(err)
	}

	return boolValue
}

func (c *Config) Get(option string, defaultValue any, cast interface{}) any {
	var value any

	envMap := utils.EnvToMap(os.Environ())

	if envValue, ok := envMap[option]; ok {
		value = envValue
	} else if c.Repository.HasItem(option) {
		value = c.Repository.GetItem(option)
	} else {
		value = defaultValue
	}

	if cast == nil {
		return value
	} else if cast == reflect.TypeOf(true).Name() {
		return c.CastBoolean(value.(string))
	}

	return value
}

type OrderedMap struct {
	Map      map[string]reflect.Type
	MapOrder []string
}

type AutoConfig struct {
	Supported  OrderedMap
	SearchPath string
	Config     Config
}

func (ac *AutoConfig) Init(searchPath string) {
	ac.Supported = OrderedMap{
		Map: map[string]reflect.Type{
			"settings.ini": reflect.TypeOf(repository.RepositoryINI{}),
			".env":         reflect.TypeOf(repository.RepositoryEnv{}),
		},
		MapOrder: []string{"settings.ini", ".env"},
	}
	ac.SearchPath = searchPath
}

func (ac *AutoConfig) FindFile(path string) string {
	for _, configFile := range ac.Supported.MapOrder {
		filename := filepath.Join(path, configFile)

		if _, err := os.Stat(filename); err == nil {
			return filename
		}
	}

	return ""
}

func (ac *AutoConfig) Load(path string) error {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	filename := ac.FindFile(absPath)

	repositoryType := ac.Supported.Map[filepath.Base(filename)]

	repository := reflect.New(repositoryType).Interface().(repository.Repository)

	err = repository.Init(filename)
	if err != nil {
		return err
	}

	ac.Config.Repository = repository

	return nil
}

func (ac *AutoConfig) CallerPath() string {
	_, filename, _, ok := runtime.Caller(2)
	if !ok {
		return ""
	}

	return filepath.Dir(filename)
}
