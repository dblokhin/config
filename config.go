package config

import (
    "encoding/json"
    "strings"
    "strconv"
    "os"
	"io"
)

type Config map[string]interface{}

// GetString return string config value
func (conf Config) GetString(path string) string {

    result := conf.Get(path)
    if result == nil {
        return ""
    }

    // Если строка, то это результат
    switch val := result.(type) {
    case string: return val

    default:
        return ""
    }
}

// GetArray return array config value
func (conf Config) GetArray(path string) []interface{} {

    result := conf.Get(path)
    if result == nil {
        return []interface{}{}
    }

    switch val := result.(type) {
    case []interface{}: return val

    default:
        return []interface{}{}
    }
}

// GetBool return bool config value
func (conf Config) GetBool(path string) bool {

    result := conf.Get(path)
    if result == nil {
        return false
    }

    switch val := result.(type) {
    case bool: return val

    default:
        return false
    }
}

// GetInt return int64 config value. It may be in hex & oct variants
func (conf Config) GetInt(path string) int64 {

	result := conf.Get(path)
	if result == nil {
		return 0
	}

	switch val := result.(type) {
	case int: return int64(val)
	case int64: return val
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

// GetFloat64 return float64 config value
func (conf Config) GetFloat64(path string) float64 {

    result := conf.Get(path)
    if result == nil {
        return 0
    }

    switch val := result.(type) {
    case float64: return val
    case int: return float64(val)
    case int64: return float64(val)
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

// Get return config value by dotted path in json tree. Path should be like this "root.option.item"
func (conf Config) Get(path string) interface{} {
    items := strings.Split(path, ".")

    idx := 0
    value := map[string]interface{}(conf)

    // Перебор до предпоследнего элемента
    for idx < len(items) - 1 {
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

// New create config from file
func New(filename string) Config {
    file, err := os.Open(filename)
	if err != nil {
		panic(err)
	}

    return NewFromIO(file)
}

// NewFromIO create config from io.Reader
func NewFromIO(input io.Reader) Config {
	decoder := json.NewDecoder(input)
	decoder.UseNumber()

	res := make(Config)
	decoder.Decode(&res)

	if err := decoder.Decode(&res); err != nil {
		panic(err)
	}

	return res
}
