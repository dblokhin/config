package config

import (
	"encoding/json"
	"strings"
	"strconv"
	"os"
	"io"
	"context"
	"sync"
	"errors"
)

var errNoConfigData = errors.New("config: no config data in context")

// ConfigData is struct where stored configuration loaded from file or io stream
type ConfigData struct {
	mu sync.Mutex

	// config data
	data map[string]interface{}
}

// GetString returns string config value
func (cnf *ConfigData) GetString(path string) string {

	result := cnf.Get(path)
	if result == nil {
		return ""
	}

	// Если строка, то это результат
	switch val := result.(type) {
	case string:
		return val

	default:
		return ""
	}
}

// GetArray returns array config value
func (cnf *ConfigData) GetArray(path string) []interface{} {

	result := cnf.Get(path)
	if result == nil {
		return []interface{}{}
	}

	switch val := result.(type) {
	case []interface{}:
		return val

	default:
		return []interface{}{}
	}
}

// GetBool returns bool config value
func (cnf *ConfigData) GetBool(path string) bool {

	result := cnf.Get(path)
	if result == nil {
		return false
	}

	switch val := result.(type) {
	case bool:
		return val

	default:
		return false
	}
}

// GetInt returns int64 config value. It may be in hex & oct variants
func (cnf *ConfigData) GetInt(path string) int64 {

	result := cnf.Get(path)
	if result == nil {
		return 0
	}

	switch val := result.(type) {
	case int:
		return int64(val)
	case int64:
		return val
	case json.Number:
		if res, err := strconv.ParseInt(string(val), 0, 64); err != nil {
			return 0
		} else {
			return res
		}

	default:
		return 0
	}
}

// GetFloat64 returns float64 config value
func (cnf *ConfigData) GetFloat64(path string) float64 {

	result := cnf.Get(path)
	if result == nil {
		return 0
	}

	switch val := result.(type) {
	case float64:
		return val
	case int:
		return float64(val)
	case int64:
		return float64(val)
	case json.Number:
		if res, err := strconv.ParseFloat(string(val), 64); err != nil {
			return 0
		} else {
			return res
		}

	default:
		return 0
	}
}

// Get returns config value by dotted json path. Path should be like this "root.option.item"
func (cnf *ConfigData) Get(path string) interface{} {
	items := strings.Split(path, ".")

	// lock concurrent access
	cnf.mu.Lock()
	defer cnf.mu.Unlock()

	idx := 0
	value := map[string]interface{}(cnf.data)

	// Перебор до предпоследнего элемента
	for idx < len(items)-1 {
		tmp, ok := value[items[idx]]
		if !ok {
			return ""
		}

		value = tmp.(map[string]interface{})
		idx++
	}

	// Последний элемент
	return value[items[idx]]
}

// New creates config from file
func New(file string) *ConfigData {
	f, err := os.Open(file)
	if err != nil {
		panic(err)
	}
	defer f.Close()

	return NewFromReader(f)
}

// NewFromReader creates config from io.Reader
func NewFromReader(input io.Reader) *ConfigData {
	decoder := json.NewDecoder(input)
	decoder.UseNumber()

	data := make(map[string]interface{})
	if err := decoder.Decode(&data); err != nil {
		panic(err)
	}

	return &ConfigData{
		mu:   sync.Mutex{},
		data: data,
	}
}

type key int

const keyConfig key = iota

// NewContext creates new context with config
func NewContext(ctx context.Context, filename string) context.Context {
	return context.WithValue(ctx, keyConfig, New(filename))
}

// Config returns config from context
func Config(ctx context.Context) *ConfigData {
	value, ok := ctx.Value(keyConfig).(*ConfigData)
	if !ok {
		panic(errNoConfigData)
	}

	return value
}
