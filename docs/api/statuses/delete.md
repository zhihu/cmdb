# Delete object status

Delete an object status

**URL** : `/v1/api/types/:type_name/statuses/:status`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. 
  * `status=[string]` The status name. 

**Method** : `DELETE`

**Auth required** : YES

**Permissions required** : User is Administrator

**Data** : `{}`

## Success Response

**Condition** : If the object status exists.

**Code** : `204 NO CONTENT`

**Content** : `{}`

## Error Responses

**Condition** : If there was no object status with given name to delete.

**Code** : `404 NOT FOUND`

**Content** : `{}`

### Or

**Condition** : Authorized User is not Administrator.

**Code** : `403 FORBIDDEN`

**Content** : `{}`


## Notes

* Will remove object status
