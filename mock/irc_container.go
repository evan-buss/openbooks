package mock

import (
	"context"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

type ircContainer struct {
	testcontainers.Container
	URI string
}

// docker run --name ergo -d -p 6667:6667 -p 6697:6697 ghcr.io/ergochat/ergo:stable
func StartIrcServer(ctx context.Context) (*ircContainer, error) {

	req := testcontainers.ContainerRequest{
		Image:        "ghcr.io/ergochat/ergo:stable",
		ExposedPorts: []string{"6667/tcp", "6697/tcp"},
		WaitingFor:   wait.ForLog("Server running"),
	}
	container, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: req,
		Started:          true,
	})
	if err != nil {
		return nil, err
	}

	defer container.Terminate(ctx)

	uri, err := container.Endpoint(ctx, "6667")
	if err != nil {
		return nil, err
	}

	return &ircContainer{container, uri}, nil
}
