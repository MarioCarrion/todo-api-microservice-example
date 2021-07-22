# Caching

## Memcached

Used as repository for caching search values.

```
docker run \
  -d \
  -p 11211:11211 \
  memcached:1.6.9-alpine
```
