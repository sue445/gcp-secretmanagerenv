package secretmanagerenv

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"fmt"
	"github.com/cockroachdb/errors"
	"os"
	"strings"
)

// Client represents secretmanagerenv client
type Client struct {
	projectID string
	ctx       context.Context
	client    secretManagerClient
}

// NewClient creates a new Client instance
func NewClient(ctx context.Context, projectID string) (*Client, error) {
	c := &Client{
		projectID: projectID,
		ctx:       ctx,
	}

	if projectID != "" {
		client, err := secretmanager.NewClient(ctx)
		if err != nil {
			return nil, errors.WithStack(err)
		}
		c.client = client
	}

	return c, nil
}

// GetValueFromEnvOrSecretManager returns value from environment variable or Secret Manager
func (c *Client) GetValueFromEnvOrSecretManager(key string, required bool) (string, error) {
	if os.Getenv(key) != "" {
		return strings.TrimSpace(os.Getenv(key)), nil
	}

	if c.projectID == "" {
		if required {
			return "", fmt.Errorf("%s is required in environment variable", key)
		}

		return "", nil
	}

	ret, err := c.GetSecretManagerValue(key, "latest")
	if err != nil {
		if !required && strings.Contains(err.Error(), "code = NotFound") {
			return "", nil
		}
		return "", err
	}

	return ret, nil
}

// GetSecretManagerValue returns value from SecretManager
func (c *Client) GetSecretManagerValue(key string, version string) (string, error) {
	req := &secretmanagerpb.AccessSecretVersionRequest{
		Name: fmt.Sprintf("projects/%s/secrets/%s/versions/%s", c.projectID, key, version),
	}

	resp, err := c.client.AccessSecretVersion(c.ctx, req)
	if err != nil {
		return "", errors.WithStack(err)
	}

	str := string(resp.Payload.Data)

	return strings.TrimSpace(str), nil
}
