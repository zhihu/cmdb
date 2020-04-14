# List object meta values

Get the details of all object meta values

**URL** : `/v1/api/objects/:type_name/:object_name/relations/:to_type_name/:relation/:to_object_name/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The from object type name. Wildcard `-` can be used to list metas of all from types.
  * `object_name=[string]` The from object name. Wildcard `-` can be used to list metas of all from objects.
  * `to_type_name=[string]` The to object type name. Wildcard `-` can be used to list metas of all to types.
  * `relation=[string]` The relation type name. Wildcard `-` can be used to list metas of all relation types.
  * `to_object_name=[string]` The to object name. Wildcard `-` can be used to list metas of all to objects.

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
    "kind": "cmdb#relationMetaValueList",
    "next_page_token": "xxx",
    "values": [
        {
            "kind": "cmdb#relationMetaValue",
            "from_type_name": "SERVER_NODE",
            "from_object_name": "xxx.xxx.xxx",
            "relation": "OWNER",
            "to_type_name": "BUSINESS_LINE",
            "to_object_name": "xxx",
            "version": 123457,
            "meta_name": "PURCHASE_TIME",
            "value": "XXX",
            "create_time": "2017-01-15T01:30:15.01Z"
        }
    ]
}
```

## Notes

* Deleted object metas will not be included unless explicitly specified
