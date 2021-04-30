# Message Broker

## RabbitMQ

Used as repository for publishing messages to be consumed by other processes.

```
docker run \
  -d \
  -p 5672:5672 \
  -p 15672:15672 \
  rabbitmq:3-management
```

Then open http://localhost:15672 . To log in use `guest` as the value for both the username and password.
