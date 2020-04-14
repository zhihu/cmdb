# Delete object relation meta

Delete an object relation meta

**URL** : `/v1/api/relation_types/:from_type_name/:to_type_name/:relation/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `from_type_name=[string]` The from object type name.
  * `to_type_name=[string]` The to object type name.
  * `relation=[string]` The to object type name.
  * `meta_name=[string]` The relation meta name.

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
