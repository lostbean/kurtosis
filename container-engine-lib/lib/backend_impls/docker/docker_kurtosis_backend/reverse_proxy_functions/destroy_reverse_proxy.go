package reverse_proxy_functions

import (
	"context"
	"github.com/kurtosis-tech/kurtosis/container-engine-lib/lib/backend_impls/docker/docker_manager"
	"github.com/kurtosis-tech/stacktrace"
	"github.com/sirupsen/logrus"
	"time"
)

const (
	stopReverseProxyContainerTimeout = 2 * time.Second
)

// Destroys reverse proxy idempotently, returns nil if no reverse proxy reverse proxy container was found
func DestroyReverseProxy(ctx context.Context, dockerManager *docker_manager.DockerManager) error {
	_, maybeReverseProxyContainerId, err := getReverseProxyObjectAndContainerId(ctx, dockerManager)
	if err != nil {
		logrus.Warnf("Attempted to destroy reverse proxy but no reverse proxy container was found.")
		return nil
	}

	if maybeReverseProxyContainerId == "" {
		return nil
	}

	if err := dockerManager.StopContainer(ctx, maybeReverseProxyContainerId, stopReverseProxyContainerTimeout); err != nil {
		return stacktrace.Propagate(err, "An error occurred stopping the reverse proxy container with ID '%v'", maybeReverseProxyContainerId)
	}

	if err := dockerManager.RemoveContainer(ctx, maybeReverseProxyContainerId); err != nil {
		return stacktrace.Propagate(err, "An error occurred removing the reverse proxy container with ID '%v'", maybeReverseProxyContainerId)
	}

	return nil
}
