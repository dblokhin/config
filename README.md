[![License MIT](https://img.shields.io/badge/License-MIT-brightgreen.svg)](https://img.shields.io/badge/License-MIT-brightgreen.svg)

# Simple configuration in Golang (JSON)
Each config key value is requested by a key path. Path is dotted JSON path. See example and it will be clear:

## How to use JSON config
Example content `config.json`:
```json
{
  "db": {
    "enabled": true,
    "port": "33096",
    "user": "root",
    "pass": "***"
  }, 
  "some": "... and other options"
}
```
Load config from `config.json`:
```go
import (
	"github.com/dblokhin/config"
)

func main()  {
	// load from file
	conf := config.New("config.json")
	
	// easy access to variables
	if conf.GetBool("db.enabled") {
		port := conf.GetInt("db.port")
		user := conf.GetString("db.user")
		pass := conf.GetString("db.pass")
		
		// connect to database...
	}
	// ... 
}
```
Load & store it in your `context`:
```go
    // load & store it in your context.Context
	app := config.NewContext(context.Background(), "config.json")
	conf := config.Config(app)
```

## Why does it panic?
**The package doesn't return any errors, but it does panic.** In my opinion it's good *error handling* way that allows easy coding & good concentrating on that.

Lets me describe this point. If you use function that returns an errors, you have to (must) check every time annoying `if err != nil {}`. A good way of easy coding, in my opinion, is that functions could return only the useful values or only errors, like this:

`func SomeFunc() int | error`

and handles values and errors separately. Directly Golang doesn't allow it, but `defer` & `panic` allow us it. Just few examples:

#### Before. Consider some `initial` function:
```go
	// create app instance & load config
	app, err := webapp.New()
	if err != nil {
		return nil, err
	}

	app, err = config.NewContext(app, "config.json")
	if err != nil {
		return nil, err
	}

	conf, err := config.Config(app)
	if err != nil {
		return nil, err
	}
    // ...
```
The some caller is: 
```go
	// initiate the app
	app, err := someInitial()
	if err != nil {
		log.Println(err)
		os.Exit(someErrCode)
	}
```
In most cases `error` just means return function. And caller can checks the `err` again & return it too...
#### After. New nice code:
```go
	app := webapp.New()
	app = config.NewContext(app, "config.json")
	conf := config.Config(app)
```
#### But How to handle Errors?
Error handler in caller (or may be caller of caller):
```go
	// error handler
	defer func() {
		if err := recover(); err != nil {
			log.Println(err)
            os.Exit(someErrCode)
		}
	}()
    
    // initiate the app
    app := someInitial()
```
```go
	// error handler
	defer func() {
		if ret := recover(); ret != nil {
			cli.Status(http.StatusBadRequest)

			switch err := ret.(type) {
			case error:
				json.NewEncoder(cli).Encode(jsonError{
					Msg: err.Error(),
				})
			case string:
				json.NewEncoder(cli).Encode(jsonError{
					Msg: err,
				})

			default:
				log.Println("unkown error panic")
			}
		}
	}()
```

We can manipulate error values in `recover()`, we can place error handlers package.