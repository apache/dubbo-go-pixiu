# Filter


Filter provides request handle abstraction. User can combinate many filters together into filter-chain.

When receiving request from the listener, filter will handle it one by one at its pre or post phase.

Because pixiu want to provide network protocol transform function, so the filter contains network filter and the filter below network filter such as http filter.

the request processing order is as follows.

```
client -> listner -> network filter such as httpconnectionmanager -> http filter chain

```

### Http Filter List

You can find out all filters in `pkg/filter`. Just list some filters as the following.


- ratelimit the filter provides rate limit function using sentinel-go;
- timeout the filter provides timeout function;
- tracing the filter provides tracing function with jaeger;
- grpcproxy the filter provides http to grpc proxy function;
- httpproxy the filter provides http to http proxy function;
- proxywrite the filter provides http request path or header amend function;
- apiconfig the filter provides http to dubbo transform mapping function with remote filter;
- remote the filter provides http to dubbo proxy function with apiconfig filter.


### How to define custom http filter
#### step one

There are two abstraction interface: plugin and filter.

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

You should define yourself plugin and filter, then implement its function.

Otherwise, you maybe would define filter-own config like as below.

```go
// Config describe the config of Filter
type Config struct {		
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

You can initialize filter-own config instance at plugin `CreateFilter` function, and return it at `config` function.

```go

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
     return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}
```
Then pixiu will fill it's value by pixiu config yaml.

After filling config value, pixiu will call `Apply` function, you should prepare filter, such as fetch remote info etc.

when request comes, pixiu will call `PrepareFilterChain` function to allow filter add itself into request-filter-chain.

```go
func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}
```
If not use `AppendFilterFunc` to add self into filter-chain, then filter will not handle the request this time.

Finally pixiu will call `Handle` function.

```go
func (f *Filter) Handle(ctx *http.HttpContext) {
	f.handleCors(ctx)
	ctx.Next()
}

```

There are two phases during handle the request, the pre and the post. Before calling `ctx.Next` function, the phase is pre. And After calling it, the phase is post.




#### step two

register plugin into management in init function.

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

Add unnamed import in `pkg/pluginregistry/registry.go` file to make `init` function invoking.

```go
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter/cors"
```


#### step four

Add filter config in yaml file.

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

There is a simple filter located in `pkg/filter/cors` which provider http `cors` function.


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
 
