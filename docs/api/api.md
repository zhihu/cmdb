### Object types related
URI : `/v1/api/types/type_name`
* [List object types](types/list.md) : `GET /v1/api/types`
* [Get an object type](types/get.md) : `GET /v1/api/types/:type_name`
* [Create a new object type](types/create.md) : `POST /v1/api/types`
* [Delete an existing object type](types/delete.md) : `DELETE /v1/api/types/:type_name`

### Object statuses related
URI : `/v1/api/types/:type_name/statuses/:status_name`
* [List all statuses](statuses/list.md) : `GET /v1/api/types/-/statuses`
* [List statuses for a given object type](statuses/list.md) : `GET /v1/api/types/:type_name/statuses`
* [Get an object status](statuses/get.md) : `GET /v1/api/types/:type_name/statuses/:status_name`
* [Create object status for a given object type](statuses/create.md) : `POST /v1/api/types/:type_name/statuses`
* [Delete status for a given object type](statuses/delete.md) : `DELETE /v1/api/types/:type_name/statuses/:status_name`

### Object states related
URI : `/v1/api/types/:type_name/statuses/:status_name/states/:state_name`
* [List all object states](states/list.md) : `GET /v1/api/types/-/statuses/-/states`
* [List states for a given object type](states/list.md) : `GET /v1/api/types/:type_name/statuses/-/states`
* [List states for a given status](states/list.md) : `GET /v1/api/types/:type_name/statuses/:status_name/states`
* [Get object state](states/get.md) : `GET /v1/api/types/:type_name/statuses/:status_name/states/:state_name`
* [Create object state for a given status](states/create.md) : `POST /v1/api/types/:type_name/statuses/:status_name/states`
* [Delete object state](states/delete.md) : `DELETE /v1/api/types/:type_name/statuses/:status_name/states/:state_name`

### Object metas related
URI : `/v1/api/types/:type_name/metas/:meta_name`
* [List all object metas](object_metas/list.md) : `GET /v1/api/types/-/metas`
* [List metas for a given object type](object_metas/list.md) : `GET /v1/api/types/:type_name/metas`
* [Get object meta](object_metas/get.md) : `GET /v1/api/types/:type_name/metas/:meta_name`
* [Create object metas for given object type](object_metas/create.md) : `POST /v1/api/types/:type_name/metas`
* [Delete object meta](object_metas/delete.md) : `DELETE /v1/api/types/:type_name/metas/:meta_name`

### Relation type related
URI : `/v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name`
* [List all relation types](relation_types/list.md) : `GET /v1/api/relation_types/-/-`
* [List relation types for given from object type](relation_types/list.md) : `GET /v1/api/relation_types/:from_type_name/-`
* [List relation types for given from and to object types](relation_types/list.md) : `GET /v1/api/relation_types/:from_type_name/:to_type_name`
* [Get relation type](relation_types/get.md) : `GET /v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name`
* [Create relation type for given from and to object types](relation_types/create.md) : `POST /v1/api/relation_types/:from_type_name/:to_type_name`
* [Delete relation type](relation_types/delete.md) : `DELETE /v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name`

### Object relation meta related
URI : `/v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name/metas/:meta_name`
* [List all object relation metas](relation_metas/list.md) : `GET /v1/api/relation_types/-/-/-/metas`
* [List object relation metas for given from object type](relation_metas/list.md) : `GET /v1/api/relation_types/:from_type_name/-/-/metas`
* [List object relation metas for given from and to object type](relation_metas/list.md) : `GET /v1/api/relation_types/:from_type_name/:to_type_name/-/metas`
* [List object relation metas for given from and to object type as well as relation type](relation_metas/list.md) : `GET /v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name/metas`
* [Get object relation meta](relation_metas/get.md) : `GET /v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name/metas/:meta_name`
* [Create object meta for given from and to object type as well as relation type](relation_metas/create.md) : `POST /v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name/metas`
* [Delete object meta](relation_metas/delete.md) : `DELETE /v1/api/relation_types/:from_type_name/:to_type_name/:relation_type_name/metas/:meta_name`

### Object related
URI : `/v1/api/objects/:type_name/:object_name`
* [List all objects](objects/list.md) : `GET /v1/api/objects/-`
* [List objects of given type](objects/list.md) : `GET /v1/api/objects/:type_name`
* [Get object](objects/get.md) : `GET /v1/api/objects/:type_name/:object_name`
* [Create object of given type](objects/create.md) : `POST /v1/api/objects/:type_name`
* [Delete an object](objects/delete.md) : `DELETE /v1/api/objects/:type_name/:object_name`

### Object meta values related
URI : `/v1/api/objects/:type_name/:object_name/metas/:meta_name`
* [Get all meta values of all objects](object_meta_values/list.md) : `GET /v1/api/objects/-/-/metas`
* [Get all meta values of a given object type](object_meta_values/list.md) : `GET /v1/api/objects/:type_name/-/metas`
* [Get all meta values of a given object](object_meta_values/list.md) : `GET /v1/api/objects/:type_name/:object_name/metas`
* [Get a meta value of a given object](object_meta_values/get.md) : `GET /v1/api/objects/:type_name/:object_name/metas/:meta_name`
* [Update meta values of a particular object](object_meta_values/update.md) : `PUT /v1/api/objects/:type_name/:object_name/metas/:meta_name`
* [Delete a meta value of a given object](object_meta_values/delete.md) : `DELETE /v1/api/objects/:type_name/:object_name/metas/:meta_name`

### Object log related
URI : `/v1/api/objects/:type_name/:object_name/logs/:id`
* [List all logs](logs/list.md) : `GET /v1/api/objects/-/-/logs`
* [List logs for a given object type](logs/list.md) : `GET /v1/api/objects/:type_name/-/logs`
* [List logs for a particular object](logs/list.md) : `GET /v1/api/objects/:type_name/:object_name/logs`
* [Create log for a particular object](logs/create.md) : `POST /v1/api/objects/:type_name/:object_name/logs`

### Object relation related
URI : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name`
* [List relations of a given object](relations/list.md) : `GET /v1/api/objects/:type_name/:object_name/relations/-/-`
* [List relations of a given object of a particular relation type](relations/list.md) : `GET /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name`
* [Get a relation between two objects](relations/get.md) : `GET /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name`
* [Create a relation between two objects](relations/create.md) : `POST /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name`
* [Delete a relation between two objevts](relations/delete.md) : `DELETE /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name`

### Object relations meta values related
URI : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name/metas/:meta_name`
* [Get all metas of all relations](relation_meta_values/list.md) : `GET /v1/api/objects/-/-/relations/-/-/-/metas`
* [Get all metas of a given relation](relation_meta_values/list.md) : `GET /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name/metas`
* [Get a meta of a given relation](relation_meta_values/list.md) : `GET /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name/metas/:meta_name`
* [Update metas of a particular relation](relation_meta_values/update.md) : `PUT /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name/metas/:meta_name`
* [Delete a meta value of a given relation](relation_meta_values/delete.md) : `DELETE /v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation_type_name/:to_object_name/metas/:meta_name`

### Atomic batch execution
* [Execute multiple requests atomically](batch.md) : `/v1/api/batch`
