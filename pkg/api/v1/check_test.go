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

import (
	"testing"

	"github.com/gogo/protobuf/jsonpb"
)

func TestCheckObject(t *testing.T) {
	o := &Object{
		Name: "SomeThing",
		Type: "MySQL",
		Metas: map[string]*ObjectMetaValue{
			"Age": {ValueType: ValueType_INTEGER, Value: "1"},
		},
		State:  "state",
		Status: "status",
	}
	CheckObject(o)
	if o.Type != "mysql" || o.Status != "STATUS" || o.State != "STATE" || len(o.Metas) != 1 || o.Metas["age"] == nil || o.Metas["age"].Value != "1" {
		t.Fatalf("format object's field names failed: %v", o)
	}
	t.Log((&jsonpb.Marshaler{Indent: "  "}).MarshalToString(o))
}

func TestCheckObjectType(t *testing.T) {
	ot := &ObjectType{
		Name: "MySQL",
		Metas: map[string]*ObjectMeta{
			"Address": {
				Name:        "addr",
				Description: "mysql server's address",
				ValueType:   ValueType_STRING,
			},
		},
		Statuses: map[string]*ObjectStatus{
			"CREATED": {
				States: map[string]*ObjectState{
					"INIT": {
						Name: "init",
					},
				},
			},
		},
	}
	CheckObjectType(ot)
	if ot.Name != "mysql" {
		t.Fatalf("format object type's name failed: %s", ot.Name)
	}
	if len(ot.Statuses) != 1 || ot.Statuses["CREATED"] == nil || ot.Statuses["CREATED"].Name != "CREATED" {
		t.Fatalf("format object type's statuses failed: %s", ot.Statuses)
	}
	if len(ot.Statuses["CREATED"].States) != 1 || ot.Statuses["CREATED"].States["INIT"] == nil || ot.Statuses["CREATED"].States["INIT"].Name != "INIT" {
		t.Fatalf("format object type's states failed: %s", ot.Statuses["CREATED"].States)
	}
	if len(ot.Metas) != 1 || ot.Metas["address"] == nil || ot.Metas["address"].Name != "address" {
		t.Fatalf("format object type's metas failed: %s", ot.Metas)
	}

	t.Log((&jsonpb.Marshaler{Indent: "  "}).MarshalToString(ot))
}
