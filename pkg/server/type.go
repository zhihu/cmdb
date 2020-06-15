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
	"strings"

	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/storage"
)

type ObjectTypes struct {
	Storage storage.Storage
}

func (o *ObjectTypes) List(ctx context.Context, request *v1.ListObjectTypesRequest) (*v1.ListObjectTypesResponse, error) {
	types, err := o.Storage.ListObjectTypes(ctx, request)
	if err != nil {
		return nil, err
	}
	return &v1.ListObjectTypesResponse{Types: types}, nil
}

func (o *ObjectTypes) Get(ctx context.Context, request *v1.GetObjectTypeRequest) (*v1.ObjectType, error) {
	return o.Storage.GetObjectType(ctx, request.Name, false)
}

func (o *ObjectTypes) Create(ctx context.Context, objectType *v1.ObjectType) (*v1.ObjectType, error) {
	v1.CheckObjectType(objectType)
	return o.Storage.CreateObjectType(ctx, objectType)
}

func (o *ObjectTypes) Update(ctx context.Context, req *v1.ObjectTypeUpdateRequest) (*v1.ObjectType, error) {
	for i, path := range req.UpdateMask.Paths {
		req.UpdateMask.Paths[i] = strings.ToLower(path)
	}
	v1.CheckObjectType(req.Type)
	return o.Storage.UpdateObjectType(ctx, req.UpdateMask.Paths, req.Type)
}
