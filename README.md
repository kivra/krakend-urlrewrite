# URL Rewrite

This package contains custom middleware to rewrite URLs of incoming requests.
Rewrite rules are applied to all incoming requests and are evaluated before the
request is handled by KrakenD's request router.

## Installation

To install `urlrewrite` from GitHub:

    go get -u github.com/kivra/krakend-urlrewrite@<commit hash>

Then add `urlrewrite` to the KrakenD [`router_engine`](https://github.com/krakendio/krakend-ce/blob/master/router_engine.go):

```go
func NewEngine(cfg config.ServiceConfig, opt luragin.EngineOptions) *gin.Engine {
  engine := luragin.NewEngine(cfg, opt)
  engine.Use(urlrewrite.HandlerFunc(engine, cfg.ExtraConfig))
  ...
```

## Usage

For example, using the following global `extra_config`, requests to
`/api/hello_world` are routed to `/api/hello/world`:

```json
"kivra/url-rewrite": [
  {
    "pattern": "^/api/hello_(.*)$",
    "replace": "/api/hello/$1"
  }
]
```

Rewrite rules are applied consecutively in the order they are defined. For more
information on the supported Regex syntax, see the Golang [`regexp`](https://pkg.go.dev/regexp)
package.
