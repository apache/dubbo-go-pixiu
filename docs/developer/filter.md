# Filter


filter provider request handle abstraction. user can combinate many filter together into filter-chain.

when receive request from listener, filter will handle it in the order at pre or post phase.

because pixiu want offer network protocol transform function, so the filter contains network filter and the filter below network filter such as http filter

the request handle process like below
```
client -> listner -> network filter such as httpconnectionmanager -> http filter chain

```

### build-in filter

[reference](../user/response.md#timeout)

### Response filter

Response result or err.

### Default filter

### Optional filter 

call downstream service, only for call, not process the response.

#### field

> level: mockLevel

- 0:open Global mock is open, api need config `mock=true` in `api_config.yaml` will mock response. If some api need mock, you need use this.
- 1:close Not mock anyone.
- 2:all Global mock setting, all api auto mock.

result
```json
{
    "message": "success"
}
```

#### Pixiu log
```bash
2020-11-17T11:31:05.718+0800    DEBUG   remote/call.go:92       [dubbo-go-pixiu] client call resp:map[age:88 iD:3213 name:tiecheng time:<nil>]
```

### Timeout filter

api timeout control, independent config for each interface.


#### Common result

[reference](../sample/dubbo-body.md)

### How to define custom http filter
#### step one
there are two abstraction interface plugin and filter

```go
// HttpFilter describe http filter
type HttpFilter interface {
    // PrepareFilterChain add filter into chain
    PrepareFilterChain(ctx *http.HttpContext) error

    // Handle filter hook function
    Handle(ctx *http.HttpContext)

    Apply() error

    // Config get config for filter
    Config() interface{}
}

// HttpFilter describe http filter
type HttpFilter interface {
    // PrepareFilterChain add filter into chain
    PrepareFilterChain(ctx *http.HttpContext) error

    // Handle filter hook function
    Handle(ctx *http.HttpContext)

    Apply() error

    // Config get config for filter
    Config() interface{}
}
```

you should define yourself plugin and filter and implement its function

there several notice in defining

you should define filter-own config

```go
// Config describe the config of Filter
	Config struct {
		AllowOrigin []string `yaml:"allow_origin" json:"allow_origin" mapstructure:"allow_origin"`
		// AllowMethods access-control-allow-methods
		AllowMethods string `yaml:"allow_methods" json:"allow_methods" mapstructure:"allow_methods"`
		// AllowHeaders access-control-allow-headers
		AllowHeaders string `yaml:"allow_headers" json:"allow_headers" mapstructure:"allow_headers"`
		// ExposeHeaders access-control-expose-headers
		ExposeHeaders string `yaml:"expose_headers" json:"expose_headers" mapstructure:"expose_headers"`
		// MaxAge access-control-max-age
		MaxAge           string `yaml:"max_age" json:"max_age" mapstructure:"max_age"`
		AllowCredentials bool   `yaml:"allow_credentials" json:"allow_credentials" mapstructure:"allow_credentials"`
	}
```

you can initialize config instance when plugin CreateFilter, and return it at config function

```go

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
     return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}
```
then pixiu will fill it value using the value in config yaml

after fill config value, pixiu will call Apply function, you should prepare filter, such as fetch remote info etc.

when request comes, pixiu will call PrepareFilterChain function to allow filter add itself into request-filter-chain.

```go
func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}
```
if not use AppendFilterFunc to add filter, then filter will not handle the request this time

finally pixiu will call Handle function 

```go
func (f *Filter) Handle(ctx *http.HttpContext) {
	f.handleCors(ctx)
	ctx.Next()
}

```

there are two phase during request handle, pre and post. before calling ctx.Next function, there is pre phase,otherwise threre is post phase



#### step two

register plugin into management in init function

```go
const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPCorsFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}
```

#### step three

add unnamed import in pkg/pluginregistry/registry.go file to make init function invoking

```go
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/cors"
```


#### step four

add filter config in yaml file

```go
                http_filters:
                  - name: dgp.filter.http.httpproxy
                    config:
                  - name: dgp.filter.http.cors
                    config:
                      allow_origin:
                        - api.dubbo.com
                      allow_methods: ""
                      allow_headers: ""
                      expose_headers: ""
                      max_age: ""
                      allow_credentials: false
```


#### example

there is a simple filter located in pkg/filter/cors which provider http cors function

```go
type (
	// Plugin is http filter plugin.
	Plugin struct {
	}

	// Filter is http filter instance
	Filter struct {
		cfg *Config
	}

	// Config describe the config of Filter
	Config struct {
		AllowOrigin []string `yaml:"allow_origin" json:"allow_origin" mapstructure:"allow_origin"`
		// AllowMethods access-control-allow-methods
		AllowMethods string `yaml:"allow_methods" json:"allow_methods" mapstructure:"allow_methods"`
		// AllowHeaders access-control-allow-headers
		AllowHeaders string `yaml:"allow_headers" json:"allow_headers" mapstructure:"allow_headers"`
		// ExposeHeaders access-control-expose-headers
		ExposeHeaders string `yaml:"expose_headers" json:"expose_headers" mapstructure:"expose_headers"`
		// MaxAge access-control-max-age
		MaxAge           string `yaml:"max_age" json:"max_age" mapstructure:"max_age"`
		AllowCredentials bool   `yaml:"allow_credentials" json:"allow_credentials" mapstructure:"allow_credentials"`
	}
)

```

do not forget add unamed import in pkg/pluginregistry/registry
 