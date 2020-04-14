# List object logs

Get the details of object logs

**URL** : `/v1/api/objects/:type_name/:object_name/logs`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. Wildcard `-` can be used to list metas of all types.
  * `object_name=[string]` The object name. Wildcard `-` can be used to list metas of all objects.

**Query String Parameters**

* **Optional:**
  * `show_deleted=[bool]` if deleted object types should be included in the response
  * `page_token=[string]` The pagination token
  * `page_size=[int32]` The pagination size

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#objectLogList",
    "next_page_token": "xxx",
    "logs": [
        {
            "kind": "cmdb#objectLog",
            "type_name": "xxx.xxx.xxx",
            "object_name": "xxx.xxx.xxx",
            "level": "EMERGENCY",
            "format": "TEXT",
            "source": "USER",
            "message": "log",
            "created_by": "administrator",
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#objectLog",
            "type_name": "xxx.xxx.xxx",
            "object_name": "xxx.xxx.xxx",
            "level": "EMERGENCY",
            "format": "TEXT",
            "source": "USER",
            "message": "log",
            "created_by": "administrator",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted object logs will not be included unless explicitly specified
