package urlrewrite

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
	"github.com/luraproject/lura/v2/config"
)

func TestUrlRewriteWithExtraCfg(t *testing.T) {
	rules := make([]Rule, 2)
	rules[0] = Rule{Pattern: "^/api/hello_(.*)$", Replace: "/api/hello/test/$1"}
	rules[1] = Rule{Pattern: "^/api/hello/test/(.*)$", Replace: "/api/hello/$1"}

	router := makeRouter(extraConfig(rules), Config{})
	res := performRequest("/api/hello_world", router)

	if res.Code != http.StatusOK {
		t.Fatalf("returned %v. should return 200", res.Code)
	}
}

func TestUrlRewriteWithDefaultCfg(t *testing.T) {
	defaultCfg := append(
		Config{},
		Rule{Pattern: "^/api/hello_(.*)$", Replace: "/api/hello/test/$1"},
		Rule{Pattern: "^/api/hello/test/(.*)$", Replace: "/api/hello/$1"},
	)

	router := makeRouter(extraConfig([]Rule{}), defaultCfg)
	res := performRequest("/api/hello_world", router)

	if res.Code != http.StatusOK {
		t.Fatalf("returned %v. should return 200", res.Code)
	}
}

func TestUrlRewriteWithDefaultAndExtraCfg(t *testing.T) {
	rules := make([]Rule, 1)
	rules[0] = Rule{Pattern: "^/api/hello_(.*)$", Replace: "/api/hello/test/$1"}

	defaultCfg := append(
		Config{},
		Rule{Pattern: "^/api/hello/test/(.*)$", Replace: "/api/hello/$1"},
	)

	router := makeRouter(extraConfig(rules), defaultCfg)
	res := performRequest("/api/hello_world", router)

	if res.Code != http.StatusOK {
		t.Fatalf("returned %v. should return 200", res.Code)
	}
}

func makeRouter(extraCfg config.ExtraConfig, defaultCfg Config) *gin.Engine {
	router := gin.New()

	router.Use(HandlerFunc(router, extraCfg, defaultCfg))
	router.GET("/api/hello/world", func(c *gin.Context) { c.String(http.StatusOK, "OK") })
	return router
}

func extraConfig(rules []Rule) map[string]interface{} {
	extraConfig := make(map[string]interface{})
	extraConfig[Namespace] = rules
	return extraConfig
}

func performRequest(path string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, path, new(bytes.Buffer))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
