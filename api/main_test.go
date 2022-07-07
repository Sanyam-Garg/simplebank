package api

import (
	"os"
	"testing"

	"github.com/gin-gonic/gin"
)

func TestMain(m *testing.M) {
	gin.SetMode(gin.TestMode)
	// m.Run() starts the main test and exits returning an exit code. The os.Exit then performs operations according to the code.
	os.Exit(m.Run())
}
