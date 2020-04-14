# Get object status

Get the details of a given object status

**URL** : `/v1/api/types/:type_name/statuses/:status`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. 
  * `status=[string]` The status name. 

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#status",
    "type_name": "BUDGET",
    "status": "CREATED",
    "description": "Budget is created",
    "create_time": "2017-01-15T01:30:15.01Z",
    "delete_time": "2017-01-15T01:30:15.01Z"
}
```
