# List objects

Get the details of objects

**URL** : `/v1/api/objects/:type_name`

* **Required:**
  * `type_name=[string]` The object type name. Wildcard `-` can be used to list objects of all object types.

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
    "kind": "cmdb#objectList",
    "next_page_token": "xxx",
    "objects": [
        {
            "kind": "cmdb#object",
            "type_name": "SERVER_NODE",
            "object_name": "abc.def.com",
            "description": "server",
            "status": "NEW",
            "state": "STARTING", 
            "create_time": "2017-01-15T01:30:15.01Z"
        },
        {
            "kind": "cmdb#object",
            "type_name": "BUDGET",
            "object_name": "xxx_2020Q2",
            "description": "2020 Q2 budget of xxx department",
            "status": "REVIEWED",
            "state": "NEW",
            "create_time": "2017-01-15T01:30:15.01Z",
            "delete_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted objects will not be included unless explicitly specified
