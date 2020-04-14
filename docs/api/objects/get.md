# Get object

Get the details of a given object

**URL** : `/v1/api/objects/:type_name/:object_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name.
  * `object_name=[string]` The object name.

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#object",
    "type_name": "SERVER_NODE",
    "object_name": "abc.def.com",
    "description": "server",
    "status": "NEW",
    "state": "STARTING",
    "create_time": "2017-01-15T01:30:15.01Z"
}
```
