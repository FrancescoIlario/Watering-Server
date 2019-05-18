#! /bin/bash

### Network
# docker network create watering

### Database
# docker run -d --rm -p 9042:9042 -v watering-cassandra-vol:/var/lib/cassandra/data --name watering-cassandra --net watering cassandra:latest

docker run -d --rm -p 5432:5432 --name watering-postgres -v postgres-default-data:/var/lib/postgresql/data postgres

### RabbitMQ
#### Management

docker run --rm -d --network host --name watering-rabbitmq-mgt -e RABBITMQ_DEFAULT_USER=guest -e RABBITMQ_DEFAULT_PASS=guest  rabbitmq:3-management

