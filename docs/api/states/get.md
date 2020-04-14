# Get object state

Get the details of a given object state

**URL** : `/v1/api/types/:type_name/statuses/:status/states/:state`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name.
  * `status=[string]` The status name.
  * `state=[string]` The state name.

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#state",
    "status": "ANY",
    "state": "RUNNING",
    "description": "Running",
    "create_time": "2017-01-15T01:30:15.01Z",
    "delete_time": "2017-01-15T01:30:15.01Z"
}
```
