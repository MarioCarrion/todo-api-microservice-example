# Search Engine (using ElasticSearch)

* [Official Elasticsearch Guide](https://www.elastic.co/guide/en/elasticsearch/reference/7.12/index.html)

Please review the `elasticsearch*` **services** in [compose.yml](../compose.yml), the code to index,
and get records is in the [elasticsearch](../internal/elasticsearch) package.

In order to correctly map results you must run the following **before creating new records**:

```
curl -X PUT -H 'Content-Type: application/json' "http://localhost:9200/tasks" -d '
{
  "mappings": {
    "properties": {
      "id": {
        "type": "keyword"
      },
      "description": {
        "type": "text"
      }
    }
  }
}'
```

This configuration is set automatically when using `docker compose up`, as part of the [`compose.yml`](../compose.yml)
file, in the `elasticsearch_setup` service.
