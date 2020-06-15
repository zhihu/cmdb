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

package server

import (
	"context"

	"github.com/grpc-ecosystem/grpc-gateway/runtime"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"google.golang.org/grpc"
)

type Server struct {
	Objects     *Objects
	ObjectTypes *ObjectTypes
}

func (s *Server) Register(server *grpc.Server, mux *runtime.ServeMux) {
	v1.RegisterObjectsServer(server, s.Objects)
	v1.RegisterObjectTypesServer(server, s.ObjectTypes)
	_ = v1.RegisterObjectsHandlerServer(context.Background(), mux, s.Objects)
	_ = v1.RegisterObjectTypesHandlerServer(context.Background(), mux, s.ObjectTypes)
}
