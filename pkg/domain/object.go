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

package domain

type Object struct {
	id          int64
	statusID    int64
	typeID      int64
	stateID     int64
	Type        string           `json:"type"`
	Name        string           `json:"name"`
	Metas       map[string]*Meta `json:"metas"`
	Status      string           `json:"status"`
	State       string           `json:"state"`
	Description string           `json:"description"`
}

type Meta struct {
	metaID   int64
	RawValue string
	Type     int
	Name     string      `json:"name"`
	Value    interface{} `json:"value"`
}
