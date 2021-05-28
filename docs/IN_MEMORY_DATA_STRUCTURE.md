# In-Memory Data Store

## Redis

Used as repository for publishing messages to be consumed by other processes.

```
docker run \
  --rm \
  -d \
  -p 6379:6379 \
  redis:6.2.3-alpine3.13
```
