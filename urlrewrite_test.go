package urlrewrite

import (
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUrlRewrite(t *testing.T) {
	router := makeRouter()
	res := performRequest("/api/hello_world", router)

	if res.Code != http.StatusOK {
		t.Fatalf("returned %v. should return 200", res.Code)
	}
}

func makeRouter() *gin.Engine {
	router := gin.New()

	extraConfig := make(map[string]interface{})
	rules := make([]Rule, 2)
	rules[0] = Rule{Pattern: "^/api/hello_(.*)$", Replace: "/api/hello/test/$1"}
	rules[1] = Rule{Pattern: "^/api/hello/test/(.*)$", Replace: "/api/hello/$1"}
	extraConfig[Namespace] = rules

	router.Use(HandlerFunc(router, extraConfig))
	router.GET("/api/hello/world", func(c *gin.Context) { c.String(http.StatusOK, "OK") })
	return router
}

func performRequest(path string, router *gin.Engine) *httptest.ResponseRecorder {
	r := httptest.NewRequest(http.MethodGet, path, new(bytes.Buffer))
	w := httptest.NewRecorder()
	router.ServeHTTP(w, r)
	return w
}
