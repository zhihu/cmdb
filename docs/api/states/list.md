# List object states

Get the details of all object states

**URL** : `/v1/api/types/:type_name/statuses/:status/states`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. Wildcard `-` can be used to list states of all types.
  * `status=[string]` The status name. Wildcard `-` can be used to list states of all statuses.

**Query String Parameters**

* **Optional:**
  * `show_deleted=[bool]` if deleted object types should be included in the response
  * `page_token=[string]` The pagination token
  * `page_size=[int32]` The pagination size

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 Ok`

**Content examples**

```json
{
    "kind": "cmdb#stateList",
    "next_page_token": "xxx",
    "states": [
        {
            "kind": "cmdb#state",
            "status": "ANY",
            "state": "NEW",
            "description": "New",
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#state",
            "status": "ANY",
            "state": "STARTING",
            "description": "Starting",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
```

## Notes

* Deleted object states will not be included unless explicitly specified
