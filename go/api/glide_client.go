// Copyright Valkey GLIDE Project Contributors - SPDX Identifier: Apache-2.0

package api

// #cgo LDFLAGS: -lglide_rs
// #cgo !windows LDFLAGS: -lm
// #cgo darwin LDFLAGS: -framework Security
// #cgo linux, amd64 LDFLAGS: -L${SRCDIR}/../rustbin/x86_64-unknown-linux-gnu
// #cgo linux, arm64 LDFLAGS: -L${SRCDIR}/../rustbin/aarch64-unknown-linux-gnu
// #cgo darwin, arm64 LDFLAGS: -L${SRCDIR}/../rustbin/aarch64-apple-darwin
// #cgo LDFLAGS: -L../target/release
// #include "../lib.h"
import "C"

import (
	"github.com/jamesx-improving/valkey-glide/go/utils"
)

// GlideClient interface compliance check.
var _ GlideClient = (*glideClient)(nil)

// GlideClient is a client used for connection in Standalone mode.
type GlideClient interface {
	BaseClient
	GenericCommands
	ServerManagementCommands
}

// glideClient implements standalone mode operations by extending baseClient functionality.
type glideClient struct {
	*baseClient
}

// NewGlideClient creates a [GlideClient] in standalone mode using the given [GlideClientConfiguration].
func NewGlideClient(config *GlideClientConfiguration) (GlideClient, error) {
	client, err := createClient(config)
	if err != nil {
		return nil, err
	}

	return &glideClient{client}, nil
}

func (client *glideClient) CustomCommand(args []string) (interface{}, error) {
	res, err := client.executeCommand(C.CustomCommand, args)
	if err != nil {
		return nil, err
	}
	return handleInterfaceResponse(res)
}

func (client *glideClient) ConfigSet(parameters map[string]string) (Result[string], error) {
	result, err := client.executeCommand(C.ConfigSet, utils.MapToString(parameters))
	if err != nil {
		return CreateNilStringResult(), err
	}
	return handleStringResponse(result)
}

func (client *glideClient) ConfigGet(args []string) (map[Result[string]]Result[string], error) {
	res, err := client.executeCommand(C.ConfigGet, args)
	if err != nil {
		return nil, err
	}
	return handleStringToStringMapResponse(res)
}

func (client *glideClient) Select(index int64) (Result[string], error) {
	result, err := client.executeCommand(C.Select, []string{utils.IntToString(index)})
	if err != nil {
		return CreateNilStringResult(), err
	}

	return handleStringResponse(result)
}
