package decouple

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"reflect"
	"slices"
	"strings"

	"gopkg.in/ini.v1"
)

var TRUE_VALUES = []string{"y", "yes", "t", "true", "on", "1"}
var FALSE_VALUES = []string{"n", "no", "f", "false", "off", "0"}

func StringToBool(value any) (bool, error) {
	switch value.(type) {
	case bool:
		return value.(bool), nil
	case string:
		valueStr := strings.ToLower(value.(string))

		if slices.Contains(TRUE_VALUES, valueStr) {
			return true, nil
		}

		if slices.Contains(FALSE_VALUES, valueStr) {
			return false, nil
		}
	}

	return false, errors.New(fmt.Sprintf("Invalid truth value: %v", value))
}

type Undefined struct{}

func UndefinedValueError(message string) error {
	return errors.New(message)
}

func EnvToMap(envVariables []string) map[string]string {
	envMap := make(map[string]string)

	for _, env := range envVariables {
		key, value, found := strings.Cut(env, "=")

		if found {
			envMap[key] = value
		}
	}

	return envMap
}

type Config struct {
	Repository RepositoryEmpty
}

func (c *Config) Init(repository RepositoryEmpty) {
	c.Repository = repository
}

func (c *Config) CastBoolean(value string) bool {
	if value == "" {
		return false
	}

	boolValue, err := StringToBool(value)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	return boolValue
}

func (c *Config) Get(option string, defaultValue any, cast interface{}) (any, error) {
	var value any

	envMap := EnvToMap(os.Environ())

	if envValue, ok := envMap[option]; ok {
		value = envValue
	} else if c.Repository.Contains(option) {
		value = c.Repository.GetItem(option)
	} else {
		if _, ok := defaultValue.(Undefined); ok {
			return nil, UndefinedValueError(fmt.Sprintf("%s not found. Declare it as envvar or define a default value.", option))
		}

		value = defaultValue
	}

	if _, ok := cast.(Undefined); ok {
		return value, nil
	} else if cast == reflect.TypeOf(true).Name() {
		return c.CastBoolean(value.(string)), nil
	}

	return value, nil
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
			"settings.ini": reflect.TypeOf(RepositoryIni{}),
			".env":         reflect.TypeOf(RepositoryEnv{}),
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

func (ac *AutoConfig) Load(path string) {
	absPath, err := filepath.Abs(path)
	if err != nil {
		absPath = path
	}

	filename := ac.FindFile(absPath)

	repositoryType := ac.Supported.Map[filepath.Base(filename)]

	repository := reflect.New(repositoryType).Interface().(RepositoryEmpty)
	repository.Init(filename)

	ac.Config.Repository = repository
}

type RepositoryEmpty interface {
	Init(string)
	Contains(string) bool
	GetItem(string) any
}

type RepositoryIni struct {
	Section *ini.Section
}

func (repo *RepositoryIni) Init(source string) {
	cfg, err := ini.Load(source)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}

	repo.Section = cfg.Section("settings")
}

func (repo *RepositoryIni) Contains(key string) bool {
	envMap := EnvToMap(os.Environ())

	if _, ok := envMap[key]; ok {
		return true
	}

	return repo.Section.HasKey(key)
}

func (repo *RepositoryIni) GetItem(key string) any {
	return repo.Section.Key(key).Value()
}

type RepositoryEnv struct {
	Data map[string]interface{}
}

func (repo *RepositoryEnv) Init(source string) {
	repo.Data = make(map[string]interface{})

	file, err := os.Open(source)
	if err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
	defer file.Close()

	scanner := bufio.NewScanner(file)

	for scanner.Scan() {
		line := scanner.Text()
		line = strings.TrimSpace(line)

		if line == "" || strings.HasPrefix(line, "#") || !strings.Contains(line, "=") {
			continue
		}

		key, value, _ := strings.Cut(line, "=")

		key = strings.TrimSpace(key)
		value = strings.TrimSpace(value)

		valueHasSingleQuotes := strings.HasPrefix(value, "'") && strings.HasSuffix(value, "'")
		valueHasDoubleQuotes := strings.HasPrefix(value, "\"") && strings.HasSuffix(value, "\"")

		if len(value) >= 2 && (valueHasSingleQuotes || valueHasDoubleQuotes) {
			value = strings.TrimPrefix(value, "'")
			value = strings.TrimSuffix(value, "'")
			value = strings.TrimPrefix(value, "\"")
			value = strings.TrimSuffix(value, "\"")
		}

		repo.Data[key] = value
	}

	if err := scanner.Err(); err != nil {
		fmt.Println(err)

		os.Exit(1)
	}
}

func (repo *RepositoryEnv) Contains(key string) bool {
	envMap := EnvToMap(os.Environ())

	if _, ok := envMap[key]; ok {
		return true
	}

	if _, ok := repo.Data[key]; ok {
		return true
	}

	return false
}

func (repo *RepositoryEnv) GetItem(key string) any {
	if _, ok := repo.Data[key]; ok {
		return repo.Data[key]
	}

	return nil
}

func New(path string) *Config {
	config := new(AutoConfig)
	config.Init(path)
	config.Load(path)

	return &config.Config
}
