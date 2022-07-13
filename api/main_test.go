package api

import (
	"os"
	"testing"
	"time"

	db "github.com/Sanyam-Garg/simplebankgo/db/sqlc"
	"github.com/Sanyam-Garg/simplebankgo/util"
	"github.com/gin-gonic/gin"
	"github.com/stretchr/testify/require"
)

func newTestServer(t *testing.T, store db.Store)*Server{
	config := util.Config{
		TokenSymmetricKey: util.RandomString(32),
		AccessTokenDuration: time.Minute,
	}

	server, err := NewServer(config, store)
	require.NoError(t, err)
	require.NotEmpty(t, server)

	return server
}

func TestMain(m *testing.M) {
	// to get proper std_out
	gin.SetMode(gin.TestMode)
	// m.Run() starts the main test and exits returning an exit code. The os.Exit then performs operations according to the code.
	os.Exit(m.Run())
}
