# Search Engine (using ElasticSearch)

* [Official Elasticsearch Guide](https://www.elastic.co/guide/en/elasticsearch/reference/7.12/index.html)

Used as repository for persisting searchable data.

```
docker run \
  -d \
  -p 9200:9200 \
  -p 9300:9300 \
  -e "discovery.type=single-node" \
  docker.elastic.co/elasticsearch/elasticsearch:7.12.0
```
