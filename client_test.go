package secretmanagerenv

import (
	"context"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/sue445/gcp-secretmanagerenv/mock_secretmanagerenv"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"os"
	"testing"
)

func setupSecretManagerMock(ctx context.Context, t *testing.T) secretManagerClient {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/test/secrets/SECRET_MANAGER_KEY/versions/latest",
	}

	resp := &secretmanagerpb.AccessSecretVersionResponse{
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte("secret_value"),
		},
	}

	m := mock_secretmanagerenv.NewMocksecretManagerClient(ctrl)
	m.
		EXPECT().
		AccessSecretVersion(ctx, req).
		Return(resp, nil)

	return m
}

func TestClient_GetSecretManagerValue(t *testing.T) {
	ctx := context.Background()
	m := setupSecretManagerMock(ctx, t)

	c := &Client{projectID: "test", ctx: ctx, client: m}

	got, err := c.GetSecretManagerValue("SECRET_MANAGER_KEY", "latest")
	if assert.NoError(t, err) {
		assert.Equal(t, "secret_value", got)
	}
}

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
