# Delete object state

Delete an object state

**URL** : `/v1/api/types/:type_name/statuses/:status/states/:state`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name.
  * `status=[string]` The status name.
  * `state=[string]` The state name.

**Method** : `DELETE`

**Auth required** : YES

**Permissions required** : User is Administrator

**Data** : `{}`

## Success Response

**Condition** : If the object state exists.

**Code** : `204 NO CONTENT`

**Content** : `{}`

## Error Responses

**Condition** : If there was no object state with given name to delete.

**Code** : `404 NOT FOUND`

**Content** : `{}`

### Or

**Condition** : Authorized User is not Administrator.

**Code** : `403 FORBIDDEN`

**Content** : `{}`


## Notes

* Will remove object state
