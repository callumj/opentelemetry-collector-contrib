// Copyright The OpenTelemetry Authors
// SPDX-License-Identifier: Apache-2.0

package skywalkingreceiver // import "github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver"

// This file implements factory for skywalking receiver.

import (
	"context"
	"fmt"
	"net"
	"strconv"

	"go.opentelemetry.io/collector/component"
	"go.opentelemetry.io/collector/config/configgrpc"
	"go.opentelemetry.io/collector/config/confighttp"
	"go.opentelemetry.io/collector/config/confignet"
	"go.opentelemetry.io/collector/consumer"
	"go.opentelemetry.io/collector/receiver"

	"github.com/open-telemetry/opentelemetry-collector-contrib/receiver/skywalkingreceiver/internal/metadata"
)

const (
	// Protocol values.
	protoGRPC = "grpc"
	protoHTTP = "http"

	// Default endpoints to bind to.
	defaultGRPCBindEndpoint = "0.0.0.0:11800"
	defaultHTTPBindEndpoint = "0.0.0.0:12800"
)

// NewFactory creates a new Skywalking receiver factory.
func NewFactory() receiver.Factory {
	return receiver.NewFactory(
		metadata.Type,
		createDefaultConfig,
		receiver.WithTraces(createTracesReceiver, metadata.TracesStability))
}

// CreateDefaultConfig creates the default configuration for Skywalking receiver.
func createDefaultConfig() component.Config {
	return &Config{
		Protocols: Protocols{
			GRPC: &configgrpc.GRPCServerSettings{
				NetAddr: confignet.NetAddr{
					Endpoint:  defaultGRPCBindEndpoint,
					Transport: "tcp",
				},
			},
			HTTP: &confighttp.HTTPServerSettings{
				Endpoint: defaultHTTPBindEndpoint,
			},
		},
	}
}

// createTracesReceiver creates a trace receiver based on provided config.
func createTracesReceiver(
	_ context.Context,
	set receiver.CreateSettings,
	cfg component.Config,
	nextConsumer consumer.Traces,
) (receiver.Traces, error) {

	// Convert settings in the source c to configuration struct
	// that Skywalking receiver understands.
	rCfg := cfg.(*Config)

	var err error
	var c configuration
	// Set ports
	if rCfg.Protocols.GRPC != nil {
		c.CollectorGRPCServerSettings = *rCfg.Protocols.GRPC
		if c.CollectorGRPCPort, err = extractPortFromEndpoint(rCfg.Protocols.GRPC.NetAddr.Endpoint); err != nil {
			return nil, fmt.Errorf("unable to extract port for the gRPC endpoint: %w", err)
		}
	}

	if rCfg.Protocols.HTTP != nil {
		c.CollectorHTTPSettings = *rCfg.Protocols.HTTP
		if c.CollectorHTTPPort, err = extractPortFromEndpoint(rCfg.Protocols.HTTP.Endpoint); err != nil {
			return nil, fmt.Errorf("unable to extract port for the HTTP endpoint: %w", err)
		}
	}

	// Create the receiver.
	return newSkywalkingReceiver(&c, nextConsumer, set)
}

// extract the port number from string in "address:port" format. If the
// port number cannot be extracted returns an error.
func extractPortFromEndpoint(endpoint string) (int, error) {
	_, portStr, err := net.SplitHostPort(endpoint)
	if err != nil {
		return 0, fmt.Errorf("endpoint is not formatted correctly: %w", err)
	}
	port, err := strconv.ParseInt(portStr, 10, 0)
	if err != nil {
		return 0, fmt.Errorf("endpoint port is not a number: %w", err)
	}
	if port < 1 || port > 65535 {
		return 0, fmt.Errorf("port number must be between 1 and 65535")
	}
	return int(port), nil
}
