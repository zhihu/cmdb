# List object meta values

Get the details of object meta values

**URL** : `/v1/api/objects/:type_name/:object_name/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name. Wildcard `-` can be used to list metas of all types.
  * `object_name=[string]` The object name. Wildcard `-` can be used to list metas of all objects.

* **Optional:**
  * `meta_name=[string]` The meta name.

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
    "kind": "cmdb#objectMetaValueList",
    "next_page_token": "xxx",
    "values": [
        {
            "kind": "cmdb#objectMetaValue",
            "type_name": "SERVER_NODE",
            "object_name": "xxx.xxx.xxx",
            "meta_name": "CPU_DESCRIPTION",
            "version": 123457,
            "value": "XXX",
            "create_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted object metas will not be included unless explicitly specified
