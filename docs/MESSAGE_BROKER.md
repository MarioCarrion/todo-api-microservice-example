# Message Broker

## RabbitMQ

Used as repository for publishing messages to be consumed by other processes.

Please review the **services** in [compose.rabbitmq.yml](../compose.rabbitmq.yml), the code to publish 
and consume is in the [rabbitmq](../internal/rabbitmq) package.

Then open [http://localhost:15672](http://localhost:15672). Use `guest` as the value for both the username and the password.
