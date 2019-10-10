package config

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/pelletier/go-toml"
	"gopkg.in/yaml.v2"
	"io"
	"os"
	"strconv"
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
		convertMapKey(&config.configs)
	case "json":
		err := json.Unmarshal(buf.Bytes(), &config.configs)
		if err != nil {
			return fmt.Errorf("failed to parse the json file, error detail: %v", err)
		}
	case "toml":
		tree, err := toml.LoadReader(buf)
		if err != nil {
			return fmt.Errorf("failed to parse the toml file, error detail: %v", err)
		}
		c.configs = tree.ToMap()
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
	if len(keyList) == 0 {
		return nil
	}

	m := config.configs

	for _, val := range keyList[:len(keyList)-1] {
		face := m[val]
		if nm, ok := face.(map[string]interface{}); ok {
			m = nm
		} else {
			return nil
		}
	}

	return m[keyList[len(keyList)-1]]
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

func GetBool(key string) bool {
	return c.getBool(key)
}
func (config *Config) getBool(key string) bool {
	i := config.get(key)
	b, ok := i.(bool)
	if !ok {
		return false
	}

	return b
}

func RegisterEnv(key string) {
	c.registerEnv(key)
}
func (config *Config) registerEnv(key string) {
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

func GetEnvInt(key string) int {
	return c.getEnvInt(key)
}
func (config *Config) getEnvInt(key string) int {
	str := config.getEnv(key)
	num, err := strconv.Atoi(str)
	if err != nil {
		return 0
	}

	return num
}

func convertMapKey(m *map[string]interface{}) {
	for key, val := range *m {
		if _, ok := val.(map[interface{}]interface{}); ok {
			nm := convertOneMapKey(val.(map[interface{}]interface{}))
			(*m)[key] = nm
			convertMapKey(&nm)
		} else {
			continue
		}
	}
}

func convertOneMapKey(m map[interface{}]interface{}) map[string]interface{} {
	nm := make(map[string]interface{}, 0)
	for key, val := range m {
		nm[key.(string)] = val
	}

	return nm
}
