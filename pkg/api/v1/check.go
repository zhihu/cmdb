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

package v1

import "strings"

func CheckObjectType(o *ObjectType) *ObjectType {
	var metas = make(map[string]*ObjectMeta, len(o.Metas))
	o.Name = strings.ToLower(o.Name)
	for name, meta := range o.Metas {
		meta.Name = strings.ToLower(name)
		metas[meta.Name] = meta
	}
	o.Metas = metas
	var statuses = make(map[string]*ObjectStatus, len(o.Statuses))
	for name, status := range o.Statuses {
		status.Name = strings.ToUpper(name)
		var states = make(map[string]*ObjectState, len(status.States))
		for name, state := range status.States {
			state.Name = strings.ToUpper(name)
			states[state.Name] = state
		}
		status.States = states
		statuses[status.Name] = status
	}
	o.Statuses = statuses
	return o
}

func CheckObject(o *Object) *Object {
	o.Type = strings.ToLower(o.Type)
	o.State = strings.ToUpper(o.State)
	o.Status = strings.ToUpper(o.Status)
	var metas = make(map[string]*ObjectMetaValue, len(o.Metas))
	for name, value := range o.Metas {
		metas[strings.ToLower(name)] = value
	}
	o.Metas = metas
	return o
}
