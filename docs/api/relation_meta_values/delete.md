# Delete object meta value

Delete an object meta

**URL** : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation/:to_object_name/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The from object type name.
  * `object_name=[string]` The from object name.
  * `to_type_name=[string]` The to object type name.
  * `relation=[string]` The relation type name.
  * `to_object_name=[string]` The to object name.
  * `meta_name=[string]` The meta name.

**Method** : `DELETE`

**Auth required** : YES

**Permissions required** : User is Administrator

**Data** : `{}`

## Success Response

**Condition** : If the object meta exists.

**Code** : `204 NO CONTENT`

**Content** : `{}`

## Error Responses

**Condition** : If there was no object meta with given name to delete.

**Code** : `404 NOT FOUND`

**Content** : `{}`

### Or

**Condition** : Authorized User is not Administrator.

**Code** : `403 FORBIDDEN`

**Content** : `{}`

## Notes

* Will remove object meta
