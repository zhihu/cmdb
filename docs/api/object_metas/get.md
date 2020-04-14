# Get object meta

Get the details of a given object meta

**URL** : `/v1/api/types/:type_name/metas/:meta_name`

**URL Parameters**

* **Required:**
  * `type_name=[string]` The object type name.
  * `meta_name=[string]` The object meta name.

**Method** : `GET`

**Auth required** : YES

**Permissions required** : None

## Success Response

**Code** : `200 OK`

**Content examples**

```json
{
    "kind": "cmdb#objectMeta",
    "type_name": "SERVER_NODE",
    "meta_name": "CPU_THREADS",
    "value_type": "INTEGER",
    "description": "Number of threads per CPU core",
    "create_time": "2017-01-15T01:30:15.01Z",
    "delete_time": "2017-01-15T01:30:15.01Z"
}
```
