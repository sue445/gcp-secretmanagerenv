package secretmanagerenv

import (
	"context"
	"fmt"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"github.com/sue445/gcp-secretmanagerenv/mock_secretmanagerenv"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
	"os"
	"testing"
)

func setupSecretManagerMock(ctx context.Context, t *testing.T) *mock_secretmanagerenv.MocksecretManagerClient {
	ctrl := gomock.NewController(t)
	t.Cleanup(func() {
		ctrl.Finish()
	})

	return mock_secretmanagerenv.NewMocksecretManagerClient(ctrl)
}

func stubAccessSecretVersionWithValidResponse(ctx context.Context, m *mock_secretmanagerenv.MocksecretManagerClient) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/test/secrets/SECRET_MANAGER_KEY/versions/latest",
	}

	resp := &secretmanagerpb.AccessSecretVersionResponse{
		Payload: &secretmanagerpb.SecretPayload{
			Data: []byte("secret_value"),
		},
	}

	m.
		EXPECT().
		AccessSecretVersion(ctx, req).
		Return(resp, nil).
		AnyTimes()
}

func stubAccessSecretVersionWithInvalidResponse(ctx context.Context, m *mock_secretmanagerenv.MocksecretManagerClient) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: "projects/test/secrets/INVALID_KEY/versions/latest",
	}

	m.
		EXPECT().
		AccessSecretVersion(ctx, req).
		Return(nil, fmt.Errorf("rpc error: code = NotFound desc = Secret [projects/000000000000/secrets/INVALID_KEY] not found or has no versions")).
		AnyTimes()
}

func TestClient_GetValueFromEnvOrSecretManager_WithValidKey(t *testing.T) {
	ctx := context.Background()
	m := setupSecretManagerMock(ctx, t)
	stubAccessSecretVersionWithValidResponse(ctx, m)

	c := &Client{projectID: "test", ctx: ctx, client: m}

	os.Setenv("ENV_KEY", "env_value")
	t.Cleanup(func() {
		os.Unsetenv("ENV_KEY")
	})

	type args struct {
		key      string
		required bool
	}
	tests := []struct {
		name    string
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "Get from env",
			args: args{
				key:      "ENV_KEY",
				required: true,
			},
			want: "env_value",
		},
		{
			name: "Get from Secret Manager",
			args: args{
				key:      "SECRET_MANAGER_KEY",
				required: true,
			},
			want: "secret_value",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := c.GetValueFromEnvOrSecretManager(tt.args.key, tt.args.required)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func TestClient_GetValueFromEnvOrSecretManager_WithInvalidKey(t *testing.T) {
	ctx := context.Background()
	m := setupSecretManagerMock(ctx, t)
	stubAccessSecretVersionWithInvalidResponse(ctx, m)

	type fields struct {
		projectID string
	}
	type args struct {
		key      string
		required bool
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		want    string
		wantErr bool
	}{
		{
			name: "NotFound in Secret manager (required)",
			fields: fields{
				projectID: "test",
			},
			args: args{
				key:      "INVALID_KEY",
				required: true,
			},
			wantErr: true,
		},
		{
			name: "NotFound in Secret manager (optional)",
			fields: fields{
				projectID: "test",
			},
			args: args{
				key:      "INVALID_KEY",
				required: false,
			},
			want: "",
		},
		{
			name: "When projectID is empty, don't check Secret Manager (required)",
			fields: fields{
				projectID: "",
			},
			args: args{
				key:      "THE_KEY_WHICH_MUST_NOT_TO_BE_CALLED",
				required: true,
			},
			wantErr: true,
		},
		{
			name: "When projectID is empty, don't check Secret Manager (optional)",
			fields: fields{
				projectID: "",
			},
			args: args{
				key:      "THE_KEY_WHICH_MUST_NOT_TO_BE_CALLED",
				required: false,
			},
			want: "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &Client{projectID: tt.fields.projectID, ctx: ctx, client: m}

			got, err := c.GetValueFromEnvOrSecretManager(tt.args.key, tt.args.required)

			if tt.wantErr {
				assert.Error(t, err)
			} else {
				if assert.NoError(t, err) {
					assert.Equal(t, tt.want, got)
				}
			}
		})
	}
}

func TestClient_GetSecretManagerValue(t *testing.T) {
	ctx := context.Background()
	m := setupSecretManagerMock(ctx, t)
	stubAccessSecretVersionWithValidResponse(ctx, m)

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
	require.NoError(t, err)

	got1, err := c.GetSecretManagerValue("SECRET_MANAGER_KEY", "latest")
	if assert.NoError(t, err) {
		assert.Equal(t, os.Getenv("INTEGRATION_TEST_WANT"), got1)
	}

	_, err = c.GetSecretManagerValue("INVALID_KEY", "latest")
	assert.Error(t, err)
}
