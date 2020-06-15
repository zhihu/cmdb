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
	"context"
	"time"

	_ "github.com/go-sql-driver/mysql"
	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/model/typetables"
	"github.com/zhihu/cmdb/pkg/query"
	"github.com/zhihu/cmdb/pkg/storage"
	"github.com/zhihu/cmdb/pkg/tools/sqly"
	"google.golang.org/grpc/codes"
	errors "google.golang.org/grpc/status"
)

func (s *Storage) CreateObject(ctx context.Context, obj *v1.Object) (n *v1.Object, err error) {
	now := time.Now()
	ts, err := s.GetTS(ctx)
	if err != nil {
		return nil, err
	}
	var typID, statusID, stateID int
	var metas = map[string]model.ObjectMeta{}

	typID, statusID, stateID, metas, err = s.getType(obj)

	if err != nil {
		return nil, err
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()

	result, err := tx.ExecContext(ctx, "insert into object (type_id, name, version, relation_version, description, status_id, state_id, create_time, update_time, delete_time) VALUES (?,?,?,?,?,?,?,?,?,?)",
		typID, obj.Name, ts, ts, obj.Description, statusID, stateID, now, now, nil,
	)
	if err != nil {
		return nil, err
	}

	objectID, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	obj.CreateTime, err = ptypes.TimestampProto(now)
	if err != nil {
		return
	}
	obj.Version = ts
	for name, meta := range obj.Metas {
		m, ok := metas[name]
		if !ok {
			return nil, ErrUnknownMeta
		}
		_, err = tx.Exec("insert into object_meta_value (object_id, meta_id, value, create_time, update_time, delete_time) VALUES (?,?,?,?,?,?)",
			objectID, m.ID, meta.Value, now, now, nil)
		if err != nil {
			return
		}
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return obj, err
}

func (s *Storage) getObject(ctx context.Context, tx *sqlx.Tx, typeID int, name string) (obj *model.Object, metas []*model.ObjectMetaValue, err error) {
	obj = &model.Object{}
	err = tx.GetContext(ctx, obj, "select * from object where type_id = ? and name = ? limit 1", typeID, name)
	if err != nil {
		return nil, nil, err
	}
	err = tx.SelectContext(ctx, &metas, "select * from object_meta_value where object_id = ? and delete_time is null", obj.ID)
	if err != nil {
		return nil, nil, err
	}
	return obj, metas, nil
}

func (s *Storage) getApiObject(ctx context.Context, tx *sqlx.Tx, typ, name string) (n *v1.Object, object *model.Object, err error) {
	var t = model.ObjectType{}
	s.cache.TypeCache(func(d *typetables.Database) {
		row, ok := d.ObjectTypeTable.GetByName(typ)
		if !ok {
			return
		}
		t = row.ObjectType
	})
	if t.ID == 0 {
		return nil, nil, errors.Newf(codes.NotFound, "type not found: %s", typ).Err()
	}
	dest, metas, err := s.getObject(ctx, tx, t.ID, name)
	if err != nil {
		return nil, nil, err
	}
	n = &v1.Object{
		Type:        t.Name,
		Name:        dest.Name,
		Description: dest.Description,
		Version:     dest.Version,
		Metas:       map[string]*v1.ObjectMetaValue{},
	}
	n.CreateTime, _ = ptypes.TimestampProto(dest.CreateTime)
	s.cache.TypeCache(func(d *typetables.Database) {
		row, ok := d.ObjectTypeTable.GetByName(typ)
		if !ok {
			err = errors.Newf(codes.NotFound, "type:%s not found", typ).Err()
			return
		}
		status := row.ObjectStatus[typetables.IndexObjectStatusID{ID: dest.StatusID}]
		if status != nil {
			n.Status = status.Name
			state := status.ObjectState[typetables.IndexObjectStateID{ID: dest.StateID}]
			if state != nil {
				n.State = state.Name
			}
		}
		for _, meta := range metas {
			if m, ok := row.ObjectMeta[typetables.IndexObjectMetaID{ID: meta.MetaID}]; ok {
				n.Metas[m.Name] = &v1.ObjectMetaValue{
					ValueType: v1.ValueType(m.ValueType),
					Value:     meta.Value,
				}
			}
		}
	})
	return n, dest, err
}

func (s *Storage) GetObject(ctx context.Context, typ, name string) (n *v1.Object, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	n, _, err = s.getApiObject(ctx, tx, typ, name)
	return n, err
}

func (s *Storage) DeleteObject(ctx context.Context, typ, name string) (n *v1.Object, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	n, dest, err := s.getApiObject(ctx, tx, typ, name)
	if err != nil {
		return nil, err
	}
	execer := sqly.Execer{
		Ctx: ctx,
		Tx:  tx,
		Err: nil,
	}
	// move and delete object row
	execer.Exec("insert into deleted_object select * from object where object.id = ?", dest.ID)
	execer.Exec("update deleted_object set delete_time = now() where id = ?", dest.ID)
	execer.Exec("delete from object where id = ?", dest.ID)

	execer.Exec("insert into deleted_object_log select * from object_log where object_log.object_id = ?", dest.ID)
	execer.Exec("delete from object_log where object_log.object_id = ?", dest.ID)

	execer.Exec("insert into deleted_object_meta_value select * from object_meta_value where object_meta_value.object_id = ?", dest.ID)
	execer.Exec("delete from  object_meta_value where object_meta_value.object_id = ?", dest.ID)

	execer.Exec("insert into deleted_object_relation select * from object_relation where object_relation.from_object_id = ? or to_object_id = ?", dest.ID, dest.ID)
	execer.Exec("delete from object_relation where object_relation.from_object_id = ? or to_object_id = ?", dest.ID, dest.ID)

	execer.Exec("insert into deleted_object_relation_meta_value select * from object_relation_meta_value where object_relation_meta_value.to_object_id = ? or object_relation_meta_value.from_object_id = ?", dest.ID, dest.ID)
	execer.Exec("delete from object_relation_meta_value where object_relation_meta_value.to_object_id = ? or object_relation_meta_value.from_object_id = ?", dest.ID, dest.ID)
	if execer.Err != nil {
		return nil, execer.Err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return
}

func (s *Storage) UpdateObject(ctx context.Context, option storage.ObjectUpdateOption, obj *v1.Object) (n *v1.Object, err error) {
	ts, err := s.GetTS(ctx)
	if err != nil {
		return nil, err
	}
	var typeID, statusID, stateID int
	var nameMetas = map[string]model.ObjectMeta{}
	typeID, statusID, stateID, nameMetas, err = s.getType(obj)
	if err == ErrUnknownStatus && !option.SetStatus {
		err = nil
	}
	if err == ErrUnknownState && !option.SetState {
		err = nil
	}
	if err != nil {
		return nil, err
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var o = &model.Object{}
	err = tx.GetContext(ctx, o, `select * from object where type_id = ? and name = ?  for update`, typeID, obj.Name)
	if err != nil {
		return
	}
	if option.MatchVersion {
		if o.Version != obj.Version {
			return nil, ErrVersionMatchFailed
		}

	}
	var currentStatus, currentState string
	s.cache.TypeCache(func(d *typetables.Database) {
		status, ok := d.ObjectStatusTable.GetByID(o.StatusID)
		if ok {
			currentStatus = status.Name
		}
		state, ok := d.ObjectStateTable.GetByID(o.StateID)
		if ok {
			currentState = state.Name
		}
	})

	var args = []interface{}{ts}
	var sets string
	if option.SetStatus {
		args = append(args, statusID)
		sets += ", status_id = ? "
	} else {
		obj.Status = currentStatus
	}
	if option.SetState {
		args = append(args, stateID)
		sets += ", state_id = ? "
	} else {
		obj.State = currentState
	}
	if option.SetDescription {
		args = append(args, obj.Description)
		sets += ", description = ? "
	} else {
		obj.Description = o.Description
	}

	obj.Version = ts
	obj.CreateTime, _ = ptypes.TimestampProto(o.CreateTime)
	args = append(args, o.ID)
	_, err = tx.ExecContext(ctx, `update object set version = ? , update_time = now() `+sets+` where id = ?`, args...)
	if err != nil {
		return nil, err
	}
	var metas []*model.ObjectMetaValue
	err = tx.SelectContext(ctx, &metas, `select * from object_meta_value where object_id = ? for update `, o.ID)
	if err != nil {
		return nil, err
	}
	var originMetas = map[int]*model.ObjectMetaValue{}
	for _, meta := range metas {
		originMetas[meta.MetaID] = meta
	}
	for name, metaType := range nameMetas {
		updateMeta, exist := obj.Metas[name]
		originMeta := originMetas[metaType.ID]
		if updateMeta == nil {
			if !option.SetAllMeta && !exist {
				if originMeta == nil {
					continue
				}
				obj.Metas[name] = &v1.ObjectMetaValue{Value: originMeta.Value, ValueType: v1.ValueType(metaType.ValueType)}
				continue
			}
			if originMeta == nil {
				continue
			}
			_, err := tx.ExecContext(ctx, `update object_meta_value set delete_time = now() where meta_id = ? and object_id = ? `, metaType.ID, o.ID)
			if err != nil {
				return nil, err
			}
			continue
		}
		if originMeta == nil {
			_, err = tx.ExecContext(ctx, `insert into object_meta_value (meta_id,object_id,value) value (?,?,?)`, metaType.ID, o.ID, updateMeta.Value)
			if err != nil {
				return nil, err
			}
			continue
		}
		if originMeta.Value == updateMeta.Value && originMeta.DeleteTime == nil {
			continue
		}
		_, err := tx.ExecContext(ctx, `update object_meta_value set value = ?, delete_time = null where meta_id = ? and object_id = ? `, updateMeta.Value, metaType.ID, o.ID)
		if err != nil {
			return nil, err
		}
	}

	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return obj, nil
}

func (s *Storage) ListObjects(ctx context.Context, request *v1.ObjectListRequest) (resp *v1.ObjectListResponse, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	var objectType = model.ObjectType{}
	var nameMetas = map[string]*model.ObjectMeta{}
	s.cache.TypeCache(func(d *typetables.Database) {
		typ, ok := d.ObjectTypeTable.GetByName(request.Type)
		if !ok {
			return
		}
		objectType = typ.ObjectType
		for _, meta := range typ.ObjectMeta {
			var m = meta.ObjectMeta
			nameMetas[meta.Name] = &m
		}
	})
	if objectType.ID == 0 {
		return nil, ErrNoSuchType
	}
	var objects []*model.Object
	if request.Query != "" {
		selector, err := query.Parse(request.Query)
		if err != nil {
			return nil, err
		}
		sql, args, err := selector.QuerySQL(nameMetas)
		sql = "select * from object where id in (" + sql + ")"
		if !request.ShowDeleted {
			sql += " and delete_time is null"
		}
		err = tx.SelectContext(ctx, &objects, sql+" order by id desc", args...)
		if err != nil {
			return nil, err
		}
	} else {
		sql := "select * from object where type_id = ?"
		if !request.ShowDeleted {
			sql += " and delete_time is null"
		}
		err = tx.SelectContext(ctx, &objects, sql+" order by id desc", objectType.ID)
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
	s.cache.TypeCache(func(cache *typetables.Database) {
		for _, o := range objects {
			idList = append(idList, o.ID)
			var object = &v1.Object{
				Type:        request.Type,
				Name:        o.Name,
				Description: o.Description,
				Metas:       nil,
				Version:     o.Version,
			}
			object.CreateTime, _ = ptypes.TimestampProto(o.CreateTime)
			status, ok := cache.ObjectStatusTable.ID[typetables.IndexObjectStatusID{ID: o.StatusID}]
			if ok {
				object.Status = status.Name
				state, ok := status.ObjectState[typetables.IndexObjectStateID{ID: o.StateID}]
				if ok {
					object.State = state.Name
				}
			}
			objectMap[o.ID] = object
			resp.Objects = append(resp.Objects, object)
		}
	})

	switch request.View {
	case v1.ObjectView_BASIC:
		return resp, nil
	case v1.ObjectView_NORMAL:
		var metas []model.ObjectMetaValue
		sql, args, _ := sqlx.In("select * from object_meta_value where object_id in (?) and delete_time is null", idList)
		err := tx.SelectContext(ctx, &metas, sql, args...)
		if err != nil {
			return nil, err
		}
		s.cache.TypeCache(func(d *typetables.Database) {
			for _, meta := range metas {
				object, ok := objectMap[meta.ObjectID]
				if !ok {
					continue
				}
				if object.Metas == nil {
					object.Metas = map[string]*v1.ObjectMetaValue{}
				}
				m, ok := d.ObjectMetaTable.GetByID(meta.MetaID)
				if !ok {
					continue
				}
				object.Metas[m.Name] = &v1.ObjectMetaValue{
					ValueType: v1.ValueType(m.ValueType),
					Value:     meta.Value,
				}
			}
		})
	}
	return
}

func (s *Storage) StopWatchObjects(ctx context.Context, typ string, f storage.FilterWatcher) error {
	objects, err := s.cache.GetObjectsCache(ctx, typ)
	if err != nil {
		return err
	}
	objects.RemoveFilterWatcher(f)
	return nil
}

func (s *Storage) WatchObjects(ctx context.Context, typ string, f storage.FilterWatcher) error {
	objects, err := s.cache.GetObjectsCache(ctx, typ)
	if err != nil {
		return err
	}
	objects.AddFilterWatcher(f)
	<-ctx.Done()
	objects.RemoveFilterWatcher(f)
	return nil
}
