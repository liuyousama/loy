package config

import (
	"bytes"
	"fmt"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strings"
)

type Config struct {
	configFilePath string
	configFileType string
	configPrefix   string
	envPrefix      string
	keySep         string

	envConfigs map[string]string
	configs    map[string]interface{}
}

var c = &Config{
	configFilePath: "config.yaml",
	configFileType: "yaml",
	keySep:         ".",
}

func SetFilePath(path string) {
	c.setFilePath(path)
}
func (config *Config) setFilePath(path string) {
	config.configFilePath = path
}

func SetFileType(t string) {
	c.setFileType(t)
}
func (config *Config) setFileType(t string) {
	config.configFileType = strings.ToLower(t)
}

func SetEnvPrefix(prefix string) {
	c.setEnvPrefix(prefix)
}
func (config *Config) setEnvPrefix(prefix string) {
	config.envPrefix = prefix
}

func SetKeySeparator(sep string) {
	c.setKeySeparator(sep)
}
func (config *Config) setKeySeparator(sep string) {
	config.keySep = sep
}

func ReadConfig() error {
	return c.readConfig()
}
func (config *Config) readConfig() error {
	file, err := os.OpenFile(config.configFilePath, os.O_RDONLY, 0655)
	if err != nil {
		return fmt.Errorf("failed to open config file %s, error detail: %v", c.configFilePath, err)
	}

	err = config.unmarshalReader(file)
	if err != nil {
		return err
	}

	return nil
}

func (config *Config) unmarshalReader(r io.Reader) error {
	buf := new(bytes.Buffer)
	_, err := buf.ReadFrom(r)
	if err != nil {
		return fmt.Errorf("failed to read byte from file to buffer, error detail: %v", err)
	}

	switch config.configFileType {
	case "yaml", "yml":
		err = yaml.Unmarshal(buf.Bytes(), &config.configs)
		if err != nil {
			return fmt.Errorf("failed to parse the yaml file, error detail: %v", err)
		}
	default:
		return fmt.Errorf("the file type is not supported")
	}

	return nil
}

func Get(key string) interface{} {
	return c.get(key)
}
func (config *Config) get(key string) interface{} {
	if key == "" {
		return nil
	}

	keyList := strings.Split(key, config.keySep)
	configs := config.configs

	var ok bool
	for i := 0; i < len(keyList)-1; i++ {
		configs, ok = configs[keyList[i]].(map[string]interface{})
		if !ok {
			return nil
		}
	}

	return configs[keyList[len(keyList)-1]]
}

func GetString(key string) string {
	return c.getString(key)
}
func (config *Config) getString(key string) string {
	i := config.get(key)
	str, ok := i.(string)
	if !ok {
		return ""
	}

	return str
}

func GetInt(key string) int {
	return c.getInt(key)
}
func (config *Config) getInt(key string) int {
	i := config.get(key)
	num, ok := i.(int)
	if !ok {
		return 0
	}

	return num
}

func RegisterEnv(key string)  {
	c.registerEnv(key)
}
func (config *Config) registerEnv(key string)  {
	key = config.envPrefix + key
	value := os.Getenv(key)
	config.envConfigs[key] = value
}

func GetEnv(key string) string {
	return c.getEnv(key)
}
func (config *Config) getEnv(key string) string {
	key = config.envPrefix + key
	value := config.envConfigs[key]
	return value
}