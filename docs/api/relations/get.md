# Get object state

Get the details of a given object state

**URL** : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation/:to_object_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The from object type name.
  * `object_name=[string]` The to object name.
  * `to_type_name=[string]` The to object type name.
  * `relation=[string]` The relation type name.
  * `to_object_name=[string]` The to object name.

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#relation",
    "from_type_name": "SERVER_NODE",
    "from_object_name": "xxx.xxx",
    "to_type_name": "BUSINESS_LINE",
    "to_object_name": "xxx",
    "relation": "BELONGS",
    "create_time": "2017-01-15T01:30:15.01Z",
    "delete_time": "2017-01-15T01:30:15.01Z"
}
```
