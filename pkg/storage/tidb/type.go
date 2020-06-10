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
	"database/sql"
	"strings"
	"time"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/model/typetables"
	"github.com/zhihu/cmdb/pkg/tools/sqly"
	"google.golang.org/grpc/codes"
	errors "google.golang.org/grpc/status"
)

// GetObjectType get objectType by name.
//
func (s *Storage) GetObjectType(ctx context.Context, name string, consistent bool) (n *v1.ObjectType, err error) {
	if !consistent {
		s.cache.TypeCache(func(d *typetables.Database) {
			row, ok := d.ObjectTypeTable.GetByName(name)
			if !ok {
				err = errors.Newf(codes.NotFound, "type:%s not found", name).Err()
				return
			}
			n = convertFromCache(row)
		})
		return
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	n, _, _, _, err = s.loadTypeFromDatabase(ctx, tx, name)
	return
}

func (s *Storage) ListObjectTypes(ctx context.Context, request *v1.ListObjectTypesRequest) (list []*v1.ObjectType, err error) {
	if !request.ShowDeleted && !request.Consistent {
		s.cache.TypeCache(func(d *typetables.Database) {
			for _, objectType := range d.ObjectTypeTable.ID {
				list = append(list, convertFromCache(objectType))
			}
		})
		return list, nil
	}
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	var types []*model.ObjectType
	if request.ShowDeleted {
		err = tx.SelectContext(ctx, &types, "select * from object_type;")
	} else {
		err = tx.SelectContext(ctx, &types, "select * from object_type where delete_time is null;")
	}
	if err != nil {
		return nil, err
	}
	list, err = s.loadTypesFromDatabase(ctx, tx, types)
	return list, err
}

func (s *Storage) DeleteObjectType(ctx context.Context, name string) (n *v1.ObjectType, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return nil, err
	}
	defer tx.Rollback()
	n, typ, _, _, err := s.loadTypeFromDatabase(ctx, tx, name)
	if err != nil {
		return nil, err
	}
	_, err = tx.ExecContext(ctx, "update object_type set delete_time = now() where id = ?", typ.ID)
	return nil, err
}

func (s *Storage) ForceDeleteObjectType(ctx context.Context, name string) (n *v1.ObjectType, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	n, objectType, _, statuses, err := s.loadTypeFromDatabase(ctx, tx, name)
	if err != nil {
		return nil, err
	}
	e := &sqly.Execer{Ctx: ctx, Tx: tx}

	for _, status := range statuses {
		e.Exec("delete from object_state where status_id = ?", status.ID)
	}
	e.Exec("delete from object_status where type_id = ?", objectType.ID)
	e.Exec("delete from object_meta where type_id = ?", objectType.ID)
	e.Exec("delete from object_log where object_id in (select id from object where type_id = ?)", objectType.ID)
	e.Exec("delete from object_meta_value where object_id in (select id from object where type_id = ?)", objectType.ID)
	e.Exec("delete from object_relation where relation_type_id in (select id from object_relation_type where object_relation_type.from_type_id = ? or to_type_id = ?)", objectType.ID, objectType.ID)
	e.Exec("delete from object_relation_meta_value where relation_type_id  in (select id from object_relation_type where object_relation_type.from_type_id = ? or to_type_id = ?)", objectType.ID, objectType.ID)
	e.Exec("delete from object_relation_meta where relation_type_id  in (select id from object_relation_type where from_type_id = ? or to_type_id = ?)", objectType.ID, objectType.ID)
	e.Exec("delete from object_relation_type where from_type_id = ? or to_type_id = ?", objectType.ID, objectType.ID)
	e.Exec("delete from object_type where id = ?", objectType.ID)
	err = e.Err
	if err != nil {
		return nil, err
	}
	tx.Commit()
	return n, nil
}

func (s *Storage) CreateObjectType(ctx context.Context, typ *v1.ObjectType) (n *v1.ObjectType, err error) {
	now := time.Now()
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	result, err := tx.ExecContext(ctx,
		"insert into object_type (name, description, create_time, delete_time) VALUES (?,?,?,?)",
		typ.Name, typ.Description,
		now, nil,
	)
	if err != nil {
		return nil, err
	}
	id, err := result.LastInsertId()
	if err != nil {
		return nil, err
	}
	for _, meta := range typ.Metas {
		_, err = tx.ExecContext(ctx, "insert into object_meta (type_id, name, value_type, description, create_time, delete_time) VALUES (?,?,?,?,?,?)",
			id, meta.Name, meta.ValueType, meta.Description, now, nil,
		)
		if err != nil {
			return nil, err
		}
	}
	for _, status := range typ.Statuses {
		err = s.InsertAllObjectTypeStatus(ctx, tx, status, int(id))
		if err != nil {
			return
		}
	}
	typ.CreateTime, _ = ptypes.TimestampProto(now)
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return typ, nil
}

func (s *Storage) InsertAllObjectTypeStatus(ctx context.Context, tx *sqlx.Tx, status *v1.ObjectStatus, id int) (err error) {
	result, err := tx.ExecContext(ctx, "insert into object_status (type_id, name, description, create_time) VALUES (?,?,?,now())",
		id, status.Name, status.Description,
	)
	if err != nil {
		return
	}
	statusID, err := result.LastInsertId()
	if err != nil {
		return
	}
	for _, state := range status.States {
		_, err = tx.ExecContext(ctx, "insert into object_state (status_id, name, description, create_time) VALUES (?,?,?,now())",
			statusID, state.Name, state.Description,
		)
		if err != nil {
			return
		}
	}
	return
}

func PathNotFoundError(path string) error {
	return errors.Newf(codes.InvalidArgument, "path %s not found", path).Err()
}

func (s *Storage) loadTypesFromDatabase(ctx context.Context, tx *sqlx.Tx, types []*model.ObjectType) (list []*v1.ObjectType, err error) {
	var typIDs = make([]int, 0, len(types))
	idTypes := make(map[int]*v1.ObjectType, len(types))
	list = make([]*v1.ObjectType, 0, len(types))
	for _, objectType := range types {
		typ := convertObjectType(objectType)
		idTypes[objectType.ID] = typ
		list = append(list, typ)
		typIDs = append(typIDs, objectType.ID)

	}
	query, args, _ := sqlx.In("select * from object_meta where type_id in(?) and delete_time is null", typIDs)
	var metas []*model.ObjectMeta
	err = tx.SelectContext(ctx, &metas, query, args...)
	if err != nil {
		return nil, err
	}
	for _, meta := range metas {
		typ, ok := idTypes[meta.TypeID]
		if !ok {
			continue
		}
		if typ.Metas == nil {
			typ.Metas = make(map[string]*v1.ObjectMeta)
		}
		typ.Metas[meta.Name] = convertMeta(meta)
	}
	query, args, _ = sqlx.In("select * from object_status where type_id in (?) and delete_time is null", typIDs)
	var statuses []*model.ObjectStatus
	err = tx.SelectContext(ctx, &statuses, query, args...)
	if err != nil {
		return nil, err
	}
	var states []*model.ObjectState
	query, args, _ = sqlx.In("select object_state.* from object_state inner join object_status os on object_state.status_id = os.id where os.type_id in (?) and object_state.delete_time is null", typIDs)
	err = tx.SelectContext(ctx, &states, query, args...)
	if err != nil {
		return nil, err
	}
	for _, state := range states {
		for _, status := range statuses {
			if status.States == nil {
				status.States = map[string]*model.ObjectState{}
			}
			status.States[state.Name] = state
		}
	}
	for _, status := range statuses {
		typ, ok := idTypes[status.TypeID]
		if !ok {
			continue
		}
		if typ.Statuses == nil {
			typ.Statuses = make(map[string]*v1.ObjectStatus)
		}
		typ.Statuses[status.Name] = convertStatus(status)
	}
	return
}

func (s *Storage) loadTypeFromDatabase(ctx context.Context, tx *sqlx.Tx, name string) (*v1.ObjectType, *model.ObjectType, []*model.ObjectMeta, map[string]*model.ObjectStatus, error) {
	var typ = &model.ObjectType{}
	err := tx.GetContext(ctx, typ, "select * from object_type where name = ?", name)
	if err != nil {
		if err == sql.ErrNoRows {
			err = errors.Newf(codes.NotFound, "no such object_type: %s", name).Err()
		}
		return nil, nil, nil, nil, err
	}
	var metas []*model.ObjectMeta
	err = tx.SelectContext(ctx, &metas, "select * from object_meta where type_id = ? and delete_time is null", typ.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	statuses, err := s.loadStatuses(ctx, tx, typ.ID)
	if err != nil {
		return nil, nil, nil, nil, err
	}
	var objectType = convertObjectType(typ)
	objectType.Metas = make(map[string]*v1.ObjectMeta, len(metas))
	for _, meta := range metas {
		objectType.Metas[meta.Name] = convertMeta(meta)
	}
	objectType.Statuses = make(map[string]*v1.ObjectStatus, len(statuses))
	for _, status := range statuses {
		objectType.Statuses[status.Name] = convertStatus(status)
	}
	return objectType, typ, metas, statuses, nil
}

func (s *Storage) UpdateObjectType(ctx context.Context, paths []string, typ *v1.ObjectType) (n *v1.ObjectType, err error) {
	tx, err := s.db.Beginx()
	if err != nil {
		return
	}
	defer tx.Rollback()
	var t model.ObjectType
	err = tx.GetContext(ctx, &t, "select * from object_type where name = ? limit 1 for update", typ.Name)
	if err != nil {
		if err == sql.ErrNoRows {
			return nil, errors.Newf(codes.NotFound, "no such type: %s", typ.Name).Err()
		}
		return nil, err
	}

	for _, path := range paths {
		names := strings.Split(path, ".")
		switch names[0] {
		case "description":
			err = s.updateObjectTypeDescription(ctx, tx, t.ID, typ.Description)
			if err != nil {
				return nil, err
			}
		case "metas":
			if len(names) == 1 {
				err = s.updateObjectTypeMetas(ctx, tx, &t, typ)
				if err != nil {
					return nil, err
				}
				continue
			}
			if len(names) == 2 {
				err = s.updateObjectTypeMeta(ctx, tx, &t, typ.Metas[names[1]])
				if err != nil {
					return nil, err
				}
				continue
			}
			if len(names) == 3 {
				err = s.updateObjectTypeMetaField(ctx, tx, &t, typ.Metas[names[1]], names[2])
				if err != nil {
					return nil, err
				}
				continue
			}
		case "statuses":
			if len(names) == 1 {
				err = s.updateObjectTypeStatues(ctx, tx, t.ID, typ.Statuses)
				if err != nil {
					return nil, err
				}
				continue
			}
			// statuses.{{status_name}}
			if len(names) == 2 {
				var name = names[1]
				origin, err := s.getObjectTypeStatus(ctx, tx, t.ID, name, true)
				if err != nil {
					return nil, err
				}
				status := typ.Statuses[name]
				if status == nil {
					err = s.deleteObjectTypeStatus(ctx, tx, origin)
					if err != nil {
						return nil, err
					}
					continue
				}
				err = s.updateObjectTypeStatus(ctx, tx, origin, status)
				if err != nil {
					return nil, err
				}
				continue
			}
			// statuses.{{status_name}}.{{field}}
			if len(names) == 3 {
				statusName := names[1]
				status, ok := typ.Statuses[statusName]
				if !ok || status == nil {
					return nil, PathNotFoundError(path)
				}
				origin, err := s.getObjectTypeStatus(ctx, tx, t.ID, status.Name, true)
				if err != nil {
					return nil, err
				}
				switch names[2] {
				case "states":
					err := s.updateObjectTypeStatusStates(ctx, tx, origin, status)
					if err != nil {
						return nil, err
					}
				case "description":
					err := s.updateObjectTypeStatusDescription(ctx, tx, origin, status)
					if err != nil {
						return nil, err
					}
				default:
					return nil, PathNotFoundError(path)
				}
				continue
			}
			// statuses.{{status_name}}.states.{{state_name}}
			if len(names) == 4 {
				if names[2] != "states" {
					return nil, PathNotFoundError(path)
				}
				statusName := names[1]
				stateName := names[3]
				origin, err := s.getObjectTypeState(ctx, tx, statusName, stateName)
				if err == sql.ErrNoRows {
					err = nil
				}
				if err != nil {
					return nil, err
				}
				status, ok := typ.Statuses[statusName]
				if !ok || status == nil {
					// delete state
					if origin == nil {
						continue
					}
					err = s.deleteObjectTypeState(ctx, tx, origin)
					if err != nil {
						return nil, err
					}
					continue
				}
				state, ok := status.States[stateName]
				if !ok || state == nil {
					if origin == nil {
						continue
					}
					err = s.deleteObjectTypeState(ctx, tx, origin)
					if err != nil {
						return nil, err
					}
					continue
				}
				if origin == nil {
					originStatus, err := s.getObjectTypeStatus(ctx, tx, t.ID, statusName, false)
					if err == sql.ErrNoRows {
						return nil, errors.Newf(codes.NotFound, "not found such status: %s", statusName).Err()
					}
					if err != nil {
						return nil, err
					}
					err = s.insertObjectTypeState(ctx, tx, originStatus.ID, state)
					if err != nil {
						return nil, err
					}
				}
				err = s.updateObjectTypeState(ctx, tx, origin, state)
				if err != nil {
					return nil, err
				}
				continue
			}
			// statuses.{{status_name}}.states.{{state_name}}.{{field}}
			if len(names) == 5 {
				if names[2] != "states" {
					return nil, PathNotFoundError(path)
				}
				statusName := names[1]
				stateName := names[3]
				field := names[4]
				if field != "description" {
					return nil, PathNotFoundError(path)
				}
				status := typ.Statuses[statusName]
				if status == nil {
					return nil, PathNotFoundError(path)
				}

				state := status.States[stateName]
				if state == nil {
					return nil, PathNotFoundError(path)
				}
				origin, err := s.getObjectTypeState(ctx, tx, statusName, stateName)
				if err != nil {
					if err == sql.ErrNoRows {
						return nil, errors.Newf(codes.NotFound, "state %s not found", stateName).Err()
					}
					return nil, err
				}
				err = s.updateObjectTypeState(ctx, tx, origin, state)
				if err != nil {
					return nil, err
				}
				continue
			}
		}
	}
	n, _, _, _, err = s.loadTypeFromDatabase(ctx, tx, typ.Name)
	if err != nil {
		return nil, err
	}
	err = tx.Commit()
	if err != nil {
		return nil, err
	}
	return n, nil
}

func (s *Storage) updateObjectTypeDescription(ctx context.Context, tx *sqlx.Tx, tID int, description string) (err error) {
	_, err = tx.ExecContext(ctx, "update object_type set description = ? where id = ?", description, tID)
	return err
}

func (s *Storage) updateObjectTypeMetas(ctx context.Context, tx *sqlx.Tx, t *model.ObjectType, typ *v1.ObjectType) (err error) {
	var metas []*model.ObjectMeta
	err = tx.SelectContext(ctx, &metas, "select * from object_meta where type_id = ?", t.ID)
	if err != nil {
		return
	}
	originMetas := make(map[string]*model.ObjectMeta, len(metas))
	for _, meta := range metas {
		m := typ.Metas[meta.Name]
		if m == nil {
			if meta.DeleteTime == nil {
				// delete meta
				_, err = tx.ExecContext(ctx, "update object_meta set delete_time = now() where id = ?", meta.ID)
			}
			continue
		}
		originMetas[meta.Name] = meta
		if meta.Description != m.Description || meta.ValueType != int(m.ValueType) {
			// update meta
			_, err = tx.ExecContext(ctx, "update object_meta set delete_time = null, description = ?, value_type = ? where id = ?",
				m.Description, m.ValueType, meta.ID)
			continue
		}
	}
	for name, meta := range typ.Metas {
		if meta == nil {
			continue
		}
		_, ok := originMetas[name]
		if !ok {
			// create meta
			_, err = tx.ExecContext(ctx, "insert into object_meta(type_id, name, value_type, description, create_time) VALUES (?,?,?,?,now())",
				t.ID, meta.Name, meta.ValueType, meta.Description)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Storage) updateObjectTypeMeta(ctx context.Context, tx *sqlx.Tx, t *model.ObjectType, meta *v1.ObjectMeta) (err error) {
	var m = &model.ObjectMeta{}
	err = tx.GetContext(ctx, m, "select * from object_meta where type_id = ? and name = ?", t.ID, meta.Name)
	if err == sql.ErrNoRows {
		err = nil
	}
	if err != nil {
		return
	}
	if meta == nil {
		if m.ID == 0 {
			// do nothing
			return nil
		}
		_, err = tx.ExecContext(ctx, "update object_meta set delete_time = now() where id = ?", m.ID)
		return
	}
	if m.ID == 0 {
		// create meta
		_, err = tx.ExecContext(ctx, "insert into object_meta(type_id, name, value_type, description, create_time) VALUES (?,?,?,?,now())",
			t.ID, meta.Name, meta.ValueType, meta.Description)
	}
	_, err = tx.ExecContext(ctx, "update object_meta set delete_time = null, description = ?, value_type = ? where id = ?",
		meta.Description, meta.ValueType, m.ID)
	return
}

func (s *Storage) updateObjectTypeMetaField(ctx context.Context, tx *sqlx.Tx, t *model.ObjectType, meta *v1.ObjectMeta, field string) (err error) {
	if meta == nil {
		return
	}
	var m = &model.ObjectMeta{}
	err = tx.GetContext(ctx, m, "select * from object_meta where type_id = ? and name = ?", t.ID, meta.Name)
	if err == sql.ErrNoRows {
		err = nil
		return
	}
	if err != nil {
		return
	}
	switch field {
	case "description":
		_, err = tx.ExecContext(ctx, "update object_meta set description = ? where id = ?",
			meta.Description, m.ID)
	case "value_type":
		_, err = tx.ExecContext(ctx, "update object_meta set description = ?, value_type = ? where id = ?",
			meta.ValueType, m.ID)
	}
	return err
}

func (s *Storage) loadStatuses(ctx context.Context, tx *sqlx.Tx, typeID int) (namedStatus map[string]*model.ObjectStatus, err error) {
	var statuses []*model.ObjectStatus
	err = tx.SelectContext(ctx, &statuses, "select * from object_status where type_id = ? and delete_time is null;", typeID)
	if err != nil {
		return nil, err
	}
	var statusesIDs = make([]int, 0, len(statuses))
	for _, status := range statuses {
		statusesIDs = append(statusesIDs, status.ID)
	}
	query, args, err := sqlx.In("select * from object_state where status_id in (?) and delete_time is null", statusesIDs)
	var states []*model.ObjectState
	err = tx.SelectContext(ctx, &states, query, args...)
	if err != nil {
		return nil, err
	}
	namedStatus = make(map[string]*model.ObjectStatus, len(statuses))
	for _, status := range statuses {
		namedStatus[status.Name] = status
		status.States = map[string]*model.ObjectState{}
		for _, state := range states {
			if state.StatusID == status.ID {
				status.States[state.Name] = state
			}
		}
	}
	return
}

func (s *Storage) updateObjectTypeStatues(ctx context.Context, tx *sqlx.Tx, tID int, updateStatuses map[string]*v1.ObjectStatus) (err error) {
	namedStatus, err := s.loadStatuses(ctx, tx, tID)
	for _, status := range updateStatuses {
		origin := namedStatus[status.Name]
		if origin == nil {
			if status != nil {
				// insert status all
				err = s.InsertAllObjectTypeStatus(ctx, tx, status, tID)
				if err != nil {
					return
				}
			}
			continue
		}

		if status == nil {
			if origin.DeleteTime == nil {
				// delete status
				err = s.deleteObjectTypeStatus(ctx, tx, origin)
				if err != nil {
					return err
				}
			}
			continue
		}
		err = s.updateObjectTypeStatus(ctx, tx, origin, status)
		if err != nil {
			return err
		}
	}
	for _, status := range namedStatus {
		_, ok := updateStatuses[status.Name]
		if !ok {
			// delete status
			err = s.deleteObjectTypeStatus(ctx, tx, status)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (s *Storage) getObjectTypeStatus(ctx context.Context, tx *sqlx.Tx, typID int, name string, loadStates bool) (status *model.ObjectStatus, err error) {
	status = &model.ObjectStatus{}
	err = tx.GetContext(ctx, status, "select * from object_status where type_id = ? and name = ? limit 1", typID, name)
	if err != nil {
		return nil, err
	}
	if !loadStates {
		return status, nil
	}
	var states []*model.ObjectState
	err = tx.SelectContext(ctx, &states, "select * from object_state where status_id = ? and delete_time is null", status.ID)
	if err != nil {
		return nil, err
	}
	status.States = make(map[string]*model.ObjectState, len(states))
	for _, state := range states {
		status.States[state.Name] = state
	}
	return status, nil
}

func (s *Storage) deleteObjectTypeStatus(ctx context.Context, tx *sqlx.Tx, status *model.ObjectStatus) (err error) {
	_, err = tx.ExecContext(ctx, "update object_status set delete_time = now() where id = ?", status.ID)
	return err
}

func (s *Storage) updateObjectTypeStatus(ctx context.Context, tx *sqlx.Tx, origin *model.ObjectStatus, status *v1.ObjectStatus) (err error) {
	err = s.updateObjectTypeStatusDescription(ctx, tx, origin, status)
	if err != nil {
		return err
	}
	return s.updateObjectTypeStatusStates(ctx, tx, origin, status)
}

func (s *Storage) updateObjectTypeStatusDescription(ctx context.Context, tx *sqlx.Tx, origin *model.ObjectStatus, status *v1.ObjectStatus) (err error) {
	if origin.Description != status.Description {
		_, err = tx.ExecContext(ctx, "update object_status set description = ? , delete_time = null where id = ?", origin.ID)
		if err != nil {
			return err
		}
	}
	return
}

func (s *Storage) updateObjectTypeStatusStates(ctx context.Context, tx *sqlx.Tx, origin *model.ObjectStatus, status *v1.ObjectStatus) (err error) {
	for _, origin := range origin.States {
		state, ok := status.States[origin.Name]
		if !ok || state == nil {
			// delete state
			err = s.deleteObjectTypeState(ctx, tx, origin)
			if err != nil {
				return
			}
			continue
		}
		err = s.updateObjectTypeState(ctx, tx, origin, state)
		if err != nil {
			return
		}
	}
	for _, state := range status.States {
		_, ok := origin.States[state.Name]
		if !ok {
			err = s.insertObjectTypeState(ctx, tx, origin.ID, state)
			if err != nil {
				return
			}
		}
	}
	return nil
}

func (s *Storage) deleteObjectTypeStateByName(ctx context.Context, tx *sqlx.Tx, statusName string, stateName string) (err error) {
	state, err := s.getObjectTypeState(ctx, tx, statusName, stateName)
	if err == sql.ErrNoRows {
		return nil
	}
	return s.deleteObjectTypeState(ctx, tx, state)
}

func (s *Storage) deleteObjectTypeState(ctx context.Context, tx *sqlx.Tx, state *model.ObjectState) (err error) {
	if state.DeleteTime != nil {
		return nil
	}
	_, err = tx.ExecContext(ctx, "update object_state set delete_time = now() where id = ?", state.ID)
	return err
}

func (s *Storage) getObjectTypeState(ctx context.Context, tx *sqlx.Tx, statusName string, stateName string) (_ *model.ObjectState, err error) {
	var state model.ObjectState
	err = tx.GetContext(ctx, &state, "select object_state.* from object_state left join object_status os on object_state.status_id = os.id where os.name = ? and object_state.name = ? limit 1", statusName, stateName)
	if err != nil {
		return nil, err
	}
	return &state, nil
}

func (s *Storage) updateObjectTypeState(ctx context.Context, tx *sqlx.Tx, origin *model.ObjectState, state *v1.ObjectState) (err error) {
	if origin.Description != state.Description {
		_, err = tx.ExecContext(ctx, "update object_state set description = ? , delete_time = null  where id = ? ", state.Description, origin.ID)
	}
	return
}

func (s *Storage) insertObjectTypeState(ctx context.Context, tx *sqlx.Tx, statusID int, state *v1.ObjectState) (err error) {
	_, err = tx.ExecContext(ctx, "insert into object_state (status_id, name, description, create_time) VALUES (?,?,?,now())", statusID, state.Name, state.Description)
	return
}
