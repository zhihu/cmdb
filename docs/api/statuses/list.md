# List object statuses

Get the details of object statuses

**URL** : `/v1/api/types/:type_name/statuses`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. Wildcard `-` can be used to list statuses of all types.

**Query String Parameters**

* **Optional:**
  * `show_deleted=[bool]` if deleted object types should be included in the response
  * `page_token=[string]` The pagination token
  * `page_size=[int32]` The pagination size

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

**Parameters**


## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#statusList",
    "next_page_token": "xxx",
    "statuses": [
        {
            "kind": "cmdb#status",
            "type_name": "ANY",
            "status": "ANY",
            "description": "Placeholder that represents any stateus",
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#status",
            "type_name": "BUDGET",
            "status": "CREATED",
            "description": "Budget is created",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted object statuses will not be included unless explicitly specified
