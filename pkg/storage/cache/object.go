package cache

import (
	"context"
	"sync"

	"github.com/golang/protobuf/ptypes"
	"github.com/jmoiron/sqlx"
	"github.com/juju/loggo"
	v1 "github.com/zhihu/cmdb/pkg/api/v1"
	"github.com/zhihu/cmdb/pkg/model"
	"github.com/zhihu/cmdb/pkg/model/objects"
	"github.com/zhihu/cmdb/pkg/model/typetables"
	"github.com/zhihu/cmdb/pkg/storage"
	"github.com/zhihu/cmdb/pkg/storage/cdc"
)

var log = loggo.GetLogger("cache")

type Objects struct {
	name           string
	id             int
	loaded         bool
	memory         objects.Database
	handlers       []storage.FilterWatcher
	typ            *Types
	objects        map[int]*v1.Object
	mutex          sync.RWMutex
	bufferedEvents [][]cdc.Event
}

func NewObjects(typ *Types, id int, name string) *Objects {
	return &Objects{
		name:    name,
		id:      id,
		typ:     typ,
		objects: map[int]*v1.Object{},
		mutex:   sync.RWMutex{},
	}
}

func (t *Objects) ResetBuffer() {
	t.bufferedEvents = nil
}

func (t *Objects) LoadData(ctx context.Context, db *sqlx.DB) error {
	tx, err := db.Beginx()
	if err != nil {
		return err
	}
	defer tx.Rollback()
	var list []*model.Object
	err = tx.SelectContext(ctx, &list, `select * from object where type_id = ? and delete_time is null`, t.id)
	if err != nil {
		return err
	}
	var metas []*model.ObjectMetaValue
	err = tx.SelectContext(ctx, &metas, `select m.* from object_meta_value m left join object o on m.object_id = o.id where o.type_id = ? and o.delete_time is null and m.delete_time is null `, t.id)
	if err != nil {
		return err
	}
	_ = tx.Commit()
	t.mutex.Lock()
	t.memory.Init()
	for _, object := range list {
		t.memory.InsertObject(object)
	}
	for _, meta := range metas {
		t.memory.InsertObjectMetaValue(meta)
	}
	t.loaded = true
	t.objects = map[int]*v1.Object{}
	for id, object := range t.memory.ObjectTable.ID {
		t.objects[id.ID] = t.convert(object)
	}
	t.mutex.Unlock()

	return nil
}

func (t *Objects) RemoveFilterWatcher(f storage.FilterWatcher) {
	t.mutex.Lock()
	var nHandlers = make([]storage.FilterWatcher, 0, len(t.handlers)-1)
	for _, handler := range t.handlers {
		if handler != f {
			nHandlers = append(nHandlers, handler)
		}
	}
	t.handlers = nHandlers
	t.mutex.Unlock()
}

func (t *Objects) AddFilterWatcher(f storage.FilterWatcher) {
	t.mutex.Lock()
	t.handlers = append(t.handlers, f)
	t.mutex.Unlock()
	f.OnInit(t.Filter(f))
}

func (t *Objects) Filter(f storage.FilterWatcher) (list []*v1.Object) {
	t.mutex.RLock()
	for _, object := range t.objects {
		if f.Filter(object) {
			list = append(list, object)
		}
	}
	t.mutex.RUnlock()
	return list
}

func (t *Objects) convert(object *objects.RichObject) *v1.Object {
	var obj = &v1.Object{}
	obj.Name = object.Name
	obj.Version = object.Version
	obj.Description = object.Description
	obj.CreateTime, _ = ptypes.TimestampProto(object.CreateTime)
	var loadTypeFailed = false
	t.typ.Read(func(d *typetables.Database) {
		var typ, ok = d.ObjectTypeTable.GetByID(object.TypeID)
		if !ok {
			// TODO: fix this situation
			loadTypeFailed = true
			return
		}
		obj.Type = typ.Name
		var statusTyp = typ.ObjectStatus[typetables.IndexObjectStatusID{ID: object.StatusID}]
		obj.Status = statusTyp.Name
		obj.State = statusTyp.ObjectState[typetables.IndexObjectStateID{ID: object.StateID}].Name
		obj.Metas = map[string]*v1.ObjectMetaValue{}
		for _, value := range object.ObjectMetaValue {
			var metaID = value.MetaID
			var metaTyp = typ.ObjectMeta[typetables.IndexObjectMetaID{ID: metaID}]
			obj.Metas[metaTyp.Name] = &v1.ObjectMetaValue{
				ValueType: v1.ValueType(metaTyp.ValueType),
				Value:     value.Value,
			}
		}
	})
	if loadTypeFailed {
		log.Errorf("fail to load type: %s", t.name)
	}
	return obj
}

func (t *Objects) OnEvents(events []cdc.Event) {
	var addObjects []int
	var deleteObjects []int
	var updateObjects []int
	var needEvents []cdc.Event
	for _, event := range events {
		switch row := event.Row.(type) {
		case *model.Object:
			if row.TypeID != t.id {
				return
			}
			needEvents = append(needEvents, event)
			switch event.Type {
			case cdc.Create:
				addObjects = append(addObjects, row.ID)
			case cdc.Update:
				if row.DeleteTime != nil {
					origin, ok := t.memory.ObjectTable.GetByID(row.ID)
					if ok && origin.DeleteTime == nil {
						deleteObjects = append(deleteObjects, row.ID)
						// trigger delete
					}
					continue
				}
				updateObjects = append(updateObjects, row.ID)
			}
		case *model.ObjectMetaValue:
			needEvents = append(needEvents, event)
		}
	}
	if len(needEvents) == 0 {
		return
	}
	t.mutex.Lock()
	if !t.loaded {
		t.bufferedEvents = append(t.bufferedEvents, needEvents)
		t.mutex.Unlock()
		return
	}
	t.mutex.Unlock()
	t.memory.OnEvents(needEvents)
	for _, id := range addObjects {
		obj, ok := t.memory.ObjectTable.GetByID(id)
		if !ok {
			continue
		}
		v1Obj := t.convert(obj)
		t.mutex.Lock()
		t.objects[id] = v1Obj
		t.mutex.Unlock()

		var evt = storage.ObjectEvent{
			Object: v1Obj,
			Event:  cdc.Create,
		}
		for _, handler := range t.handlers {
			if handler.Filter(v1Obj) {
				handler.OnEvent(evt)
			}
		}
	}
	for _, id := range deleteObjects {
		t.mutex.Lock()
		var obj, ok = t.objects[id]
		if !ok {
			t.mutex.Unlock()
			continue
		}
		delete(t.objects, id)
		t.mutex.Unlock()
		var evt = storage.ObjectEvent{
			Object: obj,
			Event:  cdc.Delete,
		}
		for _, handler := range t.handlers {
			if handler.Filter(obj) {
				handler.OnEvent(evt)
			}
		}
	}
	for _, id := range updateObjects {
		obj, ok := t.memory.ObjectTable.GetByID(id)
		if !ok {
			continue
		}
		var evtType = cdc.Update
		v1Obj := t.convert(obj)
		t.mutex.Lock()
		_, ok = t.objects[id]
		if !ok {
			evtType = cdc.Create
		}
		t.objects[id] = v1Obj
		t.mutex.Unlock()

		var evt = storage.ObjectEvent{
			Object: v1Obj,
			Event:  evtType,
		}
		for _, handler := range t.handlers {
			if handler.Filter(v1Obj) {
				handler.OnEvent(evt)
			}
		}
	}
}
