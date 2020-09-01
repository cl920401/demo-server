package router

import (
	"net/http/httptest"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestUrl(t *testing.T) {
	_, r := gin.CreateTestContext(httptest.NewRecorder())
	API(r)
}
