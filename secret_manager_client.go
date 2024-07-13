package secretmanagerenv

import (
	"cloud.google.com/go/secretmanager/apiv1/secretmanagerpb"
	"context"
	"github.com/googleapis/gax-go/v2"
)

type secretManagerClient interface {
	AccessSecretVersion(ctx context.Context, req *secretmanagerpb.AccessSecretVersionRequest, opts ...gax.CallOption) (*secretmanagerpb.AccessSecretVersionResponse, error)
}
