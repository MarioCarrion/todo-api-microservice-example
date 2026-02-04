# In-Memory Data Store

## Redis

Used as repository for publishing messages to be consumed by other processes.

Please review the **services** in [compose.redis.yml](../compose.redis.yml), the code to publish 
and consume is in the [redis](../internal/redis) package.
