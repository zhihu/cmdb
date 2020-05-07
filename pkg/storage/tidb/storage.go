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
	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes/timestamp"
	"github.com/jmoiron/sqlx"
	"github.com/juju/loggo"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/query"
)

var log = loggo.GetLogger("storage")

type DSN string

func NewStorage(dsn DSN) (*Storage, error) {
	db, err := sqlx.Open("mysql", string(dsn))
	if err != nil {
		return nil, err
	}
	return &Storage{db: db}, nil
}

type Storage struct {
	db *sqlx.DB
}

func (s *Storage) ListObjects(request *v1.ObjectListRequest) (resp *v1.ObjectListResponse, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	typeInfo, err := LoadType(tx, request.Type)
	if err != nil {
		return nil, err
	}
	var objects []*model.Object
	if request.Query != "" {
		selector, err := query.Parse(request.Query)
		if err != nil {
			return nil, err
		}
		sql, args, err := selector.QuerySQL(typeInfo.NameMetas)
		sql = "select * from object where id in (" + sql + ")"
		if !request.ShowDeleted {
			sql += " and delete_time is null"
		}
		err = tx.Select(&objects, sql+" order by id desc", args...)
		if err != nil {
			return nil, err
		}
	} else {
		sql := "select * from object where type_id = ?"
		if !request.ShowDeleted {
			sql += " and delete_time is null"
		}
		err = tx.Select(&objects, sql+" order by id desc", typeInfo.ID)
		if err != nil {
			return nil, err
		}
	}
	var idList = make([]int, 0, len(objects))
	resp = &v1.ObjectListResponse{
		Kind: "cmdb#objectList",
	}
	var objectMap = map[int]*v1.Object{}
	resp.Objects = make([]*v1.Object, 0, len(objects))
	for _, o := range objects {
		idList = append(idList, o.ID)
		var object = &v1.Object{
			Type:        request.Type,
			Name:        o.Name,
			Description: o.Description,
			CreateTime: &timestamp.Timestamp{
				Seconds: o.CreateTime.Unix(),
				Nanos:   int32(o.CreateTime.Nanosecond()),
			},
			Metas: nil,
		}
		status, ok := typeInfo.Statuses[o.StatusID]
		if ok {
			object.Status = status.Name
			state, ok := typeInfo.States[o.StateID]
			if ok {
				object.State = state.Name
			}
		}
		objectMap[o.ID] = object
		resp.Objects = append(resp.Objects, object)
	}
	switch request.View {
	case v1.ObjectView_BASIC:
		return resp, nil
	case v1.ObjectView_NORMAL:
		var metas []model.ObjectMetaValue
		sql, args, _ := sqlx.In("select * from object_meta_value where object_id in (?) and delete_time is null", idList)
		err := tx.Select(&metas, sql, args...)
		if err != nil {
			return nil, err
		}
		for _, meta := range metas {
			object, ok := objectMap[meta.ObjectID]
			if !ok {
				continue
			}
			if object.Metas == nil {
				object.Metas = map[string]*v1.ObjectMeta{}
			}
			m, ok := typeInfo.Metas[meta.MetaID]
			if !ok {
				continue
			}
			object.Metas[m.Name] = &v1.ObjectMeta{
				ValueType: v1.ValueType(m.ValueType),
				Value:     meta.Value,
			}
		}
	}
	return
}

type Type struct {
	model.ObjectType
	Statuses     map[int]*model.ObjectStatus
	States       map[int]*model.ObjectState
	NameStatuses map[string]*model.ObjectStatus
	NameStates   map[string]*model.ObjectState
	Metas        map[int]*model.ObjectMeta
	NameMetas    map[string]*model.ObjectMeta
}

func LoadType(db DB, typ string) (typeInfo *Type, err error) {
	objectType := &model.ObjectType{}
	err = db.Get(objectType, "select * from object_type where name = ?", typ)
	if err != nil {
		return
	}
	var metas []*model.ObjectMeta
	err = db.Select(&metas, "select * from object_meta where type_id = ? and delete_time is null", objectType.ID)
	if err != nil {
		return
	}
	var statuses []*model.ObjectStatus
	var states []*model.ObjectState

	err = db.Select(&statuses, "select * from object_status where type_id = ? and delete_time is null", objectType.ID)
	if err != nil {
		return
	}
	if len(statuses) != 0 {
		var statusesIDList = make([]int, 0, len(statuses))
		for _, status := range statuses {
			statusesIDList = append(statusesIDList, status.ID)
		}
		sql, args, _ := sqlx.In("select * from object_state where  status_id in (?) and delete_time is null", statusesIDList)
		err = db.Select(&states, sql, args...)
		if err != nil {
			return
		}
	}
	typeInfo = &Type{
		ObjectType:   *objectType,
		Metas:        make(map[int]*model.ObjectMeta, len(metas)),
		NameMetas:    make(map[string]*model.ObjectMeta, len(metas)),
		Statuses:     make(map[int]*model.ObjectStatus, len(statuses)),
		NameStatuses: make(map[string]*model.ObjectStatus, len(statuses)),
		States:       make(map[int]*model.ObjectState, len(states)),
		NameStates:   make(map[string]*model.ObjectState, len(states)),
	}
	for _, meta := range metas {
		typeInfo.Metas[meta.ID] = meta
		typeInfo.NameMetas[meta.Name] = meta
	}
	for _, status := range statuses {
		typeInfo.Statuses[status.ID] = status
		typeInfo.NameStatuses[status.Name] = status
	}
	for _, state := range states {
		typeInfo.States[state.ID] = state
		status, ok := typeInfo.Statuses[state.StatusID]
		if ok {
			if status.States == nil {
				status.States = map[string]*model.ObjectState{}
			}
			status.States[state.Name] = state
		}
	}
	return typeInfo, nil
}

type DB interface {
	Select(dest interface{}, sql string, args ...interface{}) error
	Get(dest interface{}, sql string, args ...interface{}) error
}
