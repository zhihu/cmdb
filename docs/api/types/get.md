# Get object type

Get the detail of object type with a given name

**URL** : `/v1/api/types/:type_name`

**URL Parameters**

* **Required:**
  * type_name=[string] The object type name

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#type",
    "name": "ANY",
    "description": "Special object type that represents anything",
    "create_time": "2017-01-15T01:30:15.01Z",
    "delete_time": "2017-01-15T01:30:15.01Z"
}
```
