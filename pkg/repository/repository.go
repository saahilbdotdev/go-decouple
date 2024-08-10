package repository

import (
	"bufio"
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/saahilbdotdev/go-decouple/pkg/utils"
	"gopkg.in/ini.v1"
)

type Repository interface {
	Init(string) error
	HasItem(string) bool
	GetItem(string) string
}

type RepositoryINI struct {
	Section *ini.Section
}

func (repo *RepositoryINI) Init(source string) error {
	cfg, err := ini.Load(source)
	if err != nil {
		return errors.New(fmt.Sprintf("Error loading .ini file: %v", err))
	}

	repo.Section = cfg.Section("settings")

	return nil
}

func (repo *RepositoryINI) HasItem(key string) bool {
	envMap := utils.EnvToMap(os.Environ())

	if _, ok := envMap[key]; ok {
		return true
	}

	return repo.Section.HasKey(key)
}

func (repo *RepositoryINI) GetItem(key string) string {
	return repo.Section.Key(key).Value()
}

type RepositoryEnv struct {
	Data map[string]string
}

func (repo *RepositoryEnv) Init(source string) error {
	repo.Data = make(map[string]string)

	file, err := os.Open(source)
	if err != nil {
		return errors.New(fmt.Sprintf("Error opening .env file: %v", err))
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
		return errors.New(fmt.Sprintf("Error reading .env file: %v", err))
	}

	return nil
}

func (repo *RepositoryEnv) HasItem(key string) bool {
	if _, ok := utils.EnvToMap(os.Environ())[key]; ok {
		return true
	}

	_, exists := repo.Data[key]

	return exists
}

func (repo *RepositoryEnv) GetItem(key string) string {
	if repo.HasItem(key) {
		return repo.Data[key]
	}

	return ""
}

func NewEnv(source string) *RepositoryEnv {
	repository := new(RepositoryEnv)

	err := repository.Init(source)
	if err != nil {
		panic(err)
	}

	return repository
}

func NewINI(source string) *RepositoryINI {
	repository := new(RepositoryINI)

	err := repository.Init(source)
	if err != nil {
		panic(err)
	}

	return repository
}
