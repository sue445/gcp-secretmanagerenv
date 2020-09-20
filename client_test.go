package secretmanagerenv

import (
	"context"
	"github.com/stretchr/testify/assert"
	"os"
	"testing"
)

func TestClient_GetSecretManagerValue_IntegrationTest(t *testing.T) {
	if os.Getenv("INTEGRATION_TEST_PROJECT_ID") == "" || os.Getenv("INTEGRATION_TEST_WANT") == "" {
		return
	}

	ctx := context.Background()
	c, err := NewClient(ctx, os.Getenv("INTEGRATION_TEST_PROJECT_ID"))
	if !assert.NoError(t, err) {
		return
	}

	got, err := c.GetSecretManagerValue("SECRET_MANAGER_KEY", "latest")
	if assert.NoError(t, err) {
		assert.Equal(t, os.Getenv("INTEGRATION_TEST_WANT"), got)
	}
}
