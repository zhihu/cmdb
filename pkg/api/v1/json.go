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
	"bytes"
	"encoding/json"
	"errors"
	"strconv"
	"strings"

	"github.com/golang/protobuf/jsonpb"
)

func (x *ObjectMetaValue) MarshalJSONPB(_ *jsonpb.Marshaler) ([]byte, error) {
	switch x.ValueType {
	case ValueType_STRING:
		return json.Marshal(x.Value)
	case ValueType_BOOLEAN:
		b, _ := strconv.ParseBool(x.Value)
		return json.Marshal(b)
	case ValueType_DOUBLE:
		b, _ := strconv.ParseFloat(x.Value, 64)
		return json.Marshal(b)
	case ValueType_INTEGER:
		v, _ := strconv.ParseInt(x.Value, 10, 64)
		return json.Marshal(v)
	}
	return json.Marshal(x.Value)
}

func (x *ObjectMetaValue) UnmarshalJSONPB(_ *jsonpb.Unmarshaler, data []byte) (err error) {
	buf := bytes.NewBuffer(data)
	decoder := json.NewDecoder(buf)
	decoder.UseNumber()
	var i interface{}
	err = decoder.Decode(&i)
	if err != nil {
		return
	}
	switch c := i.(type) {
	case bool:
		x.ValueType = ValueType_BOOLEAN
		x.Value = strconv.FormatBool(c)
		return nil
	case json.Number:
		if strings.Contains(string(c), ".") {
			x.ValueType = ValueType_DOUBLE
			x.Value = string(c)
		} else {
			x.ValueType = ValueType_INTEGER
			x.Value = string(c)
		}
	case string:
		x.ValueType = ValueType_STRING
		x.Value = c
	}
	return ErrUnknownValueType
}

var ErrUnknownValueType = errors.New("unknown value type")
