package secretmanagerenv

import (
	secretmanager "cloud.google.com/go/secretmanager/apiv1"
	"context"
	"fmt"
	secretmanagerpb "google.golang.org/genproto/googleapis/cloud/secretmanager/v1"
)

// Client represents secretmanagerenv client
type Client struct {
	projectID string
	ctx       context.Context
	client    *secretmanager.Client
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

	return str, nil
}
