package secretmanagerenv

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
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
	client, err := secretmanager.NewClient(ctx)
	if err != nil {
		return nil, err
	}

	c := &Client{
		projectID: projectID,
		ctx:       ctx,
		client:    client,
	}
	return c, nil
}

// GetValueFromEnvOrSecretManager returns value from environment variable or Secret Manager
func (c *Client) GetValueFromEnvOrSecretManager(key string, required bool) (string, error) {
	if os.Getenv(key) != "" {
		return strings.TrimSpace(os.Getenv(key)), nil
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
		return "", err
	}

	str := string(resp.Payload.Data)

	return strings.TrimSpace(str), nil
}
