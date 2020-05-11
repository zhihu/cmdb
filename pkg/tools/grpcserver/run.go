// Copyright 2020 Zhizhesihai (Beijing) Technology Limited.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// See the License for the specific language governing permissions and
// limitations under the License.

package grpcserver

import (
	"context"
	"net"

	"github.com/juju/loggo"
	"google.golang.org/grpc"
)

var log = loggo.GetLogger("grpcserver")

func Run(ctx context.Context, listener net.Listener , grpcServer *grpc.Server) error {
	errCh := make(chan error, 1)
	go func() {
		errCh <- grpcServer.Serve(listener)
	}()
	select {
	case err := <-errCh:
		return err
	case <-ctx.Done():
		grpcServer.Stop()
	}
	return nil
}
