package urlrewrite

import (
	"bytes"
	"encoding/json"
	"regexp"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
)

const Namespace = "kivra/url-rewrite"

type Rule struct {
	Pattern string `json:"pattern"`
	Replace string `json:"replace"`
}

type Config []Rule

type PathTransformer func(c *gin.Context) string

func ConfigGetter(e config.ExtraConfig) *Config {
	cfg := new(Config)

	tmp, ok := e[Namespace]
	if !ok {
		return cfg
	}

	buf := new(bytes.Buffer)
	if err := json.NewEncoder(buf).Encode(tmp); err != nil {
		return cfg
	}
	if err := json.NewDecoder(buf).Decode(cfg); err != nil {
		return cfg
	}

	return cfg
}

func HandlerFunc(engine *gin.Engine, extraCfg config.ExtraConfig) gin.HandlerFunc {
	cfg := ConfigGetter(extraCfg)
	if len(*cfg) == 0 {
		return func(c *gin.Context) {
			c.Next()
		}
	}

	transformers := make([]PathTransformer, len(*cfg))
	for i, rule := range *cfg {
		rex, err := regexp.Compile(rule.Pattern)
		if err != nil {
			panic("urlrewrite: Error: invalid regex pattern: " + rule.Pattern)
		}
		replace := rule.Replace // unsafe to use loop var in func closure
		transformers[i] = func(c *gin.Context) string {
			return rex.ReplaceAllString(c.Request.URL.Path, replace)
		}
	}

	return func(c *gin.Context) {
		path := c.Request.URL.Path
		for i := range transformers {
			c.Request.URL.Path = transformers[i](c)
		}
		if path != c.Request.URL.Path {
			c.Abort()
			engine.HandleContext(c)
		} else {
			c.Next()
		}
	}
}
