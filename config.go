package config

import (
	"encoding/json"
	"strings"
	"strconv"
	"os"
	"io"
	"context"
	"sync"
)

// Config is struct where stored configuration loaded from file or io stream
type Config struct {
	mu sync.Mutex

	// config data
	data map[string]interface{}
}

// GetString returns string config value
func (cnf *Config) GetString(path string) string {

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
func (cnf *Config) GetArray(path string) []interface{} {

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
func (cnf *Config) GetBool(path string) bool {

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
func (cnf *Config) GetInt(path string) int64 {

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
func (cnf *Config) GetFloat64(path string) float64 {

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
func (cnf *Config) Get(path string) interface{} {
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
func New(file string) (*Config, error) {
	f, err := os.Open(file)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return NewFromIO(f)
}

// NewFromIO creates config from io.Reader
func NewFromIO(input io.Reader) (*Config, error) {
	decoder := json.NewDecoder(input)
	decoder.UseNumber()

	res := &Config{
		mu: sync.Mutex{},
		data: make(map[string]interface{}),
	}
	if err := decoder.Decode(&res.data); err != nil {
		return nil, err
	} else {
		return res, nil
	}
}

type key int
const keyConfig key = iota

// NewContext creates new context with config
func NewContext(ctx context.Context, filename string) (context.Context, error) {
	if conf, err := New(filename); err != nil {
		return ctx, err
	} else {
		return context.WithValue(ctx, keyConfig, conf), nil
	}
}

// FromContext returns config from context
func FromContext(ctx context.Context) (*Config, bool) {
	value, ok := ctx.Value(keyConfig).(*Config)
	return value, ok
}
