# Package to simple handle configuration in Golang (json)
Each config key value is requested by a key path. It's dotted json path. See example and it will be clear.

## How to use
JSON config content:
```json
{
  "db": {
    "enabled": true,
    "port": "33096",
    "user": "root",
    "pass": "***"
  }
  // and others options
}
```
```go
import (
	"github.com/dblokhin/config"
	// and others packages ...
)

func main()  {
	config := config.New("config.json")
	
	if config.GetBool("db.enabled") {
		port := config.GetInt("db.port")
		user := config.GetString("db.user")
		pass := config.GetString("db.pass")
		
		// connect to database...
	}
	// ... 
}
```
