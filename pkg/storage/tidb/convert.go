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

package tidb

import (
	"github.com/golang/protobuf/ptypes"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/model/typetables"
)

func convertObjectType(objectType *model.ObjectType) *v1.ObjectType {
	var ot = &v1.ObjectType{}
	ot.Name = objectType.Name
	ot.Description = objectType.Description
	ot.CreateTime, _ = ptypes.TimestampProto(objectType.CreateTime)
	if objectType.DeleteTime != nil {
		ot.DeleteTime, _ = ptypes.TimestampProto(*objectType.DeleteTime)
	}
	return ot
}

func convertFromCache(objectType *typetables.RichObjectType) *v1.ObjectType {
	var ot = convertObjectType(&objectType.ObjectType)
	ot.Metas = map[string]*v1.ObjectMeta{}
	for _, meta := range objectType.ObjectMeta {
		var m = &v1.ObjectMeta{
			Name:        meta.Name,
			Description: meta.Description,
			ValueType:   v1.ValueType(meta.ValueType),
		}
		ot.Metas[m.Name] = m
	}
	ot.Statuses = map[string]*v1.ObjectStatus{}
	for _, status := range objectType.ObjectStatus {
		var ss = &v1.ObjectStatus{
			Name:        status.Name,
			Description: status.Description,
			States:      map[string]*v1.ObjectState{},
		}
		ot.Statuses[status.Name] = ss
		for _, state := range status.ObjectState {
			var s = &v1.ObjectState{
				Name:        state.Name,
				Description: state.Description,
			}
			ss.States[s.Name] = s
		}
	}
	return ot
}

func convertStatus(status *model.ObjectStatus) *v1.ObjectStatus {
	var states = make(map[string]*v1.ObjectState, len(status.States))
	for _, state := range status.States {
		states[state.Name] = &v1.ObjectState{
			Name:        state.Name,
			Description: state.Description,
		}
	}
	return &v1.ObjectStatus{
		Name:        status.Name,
		Description: status.Description,
		States:      states,
	}
}

func convertMeta(meta *model.ObjectMeta) *v1.ObjectMeta {
	return &v1.ObjectMeta{
		Name:        meta.Name,
		Description: meta.Description,
		ValueType:   v1.ValueType(meta.ValueType),
	}
}
