# Project notification-service

Notification-Service is a microservice, containing logic for managing notifications.
It receives notifications through kafka, by subscribing on corresponding topics. 
On new notification, service saves it to the database, checks for active websocket connections,
and if any connection associated with received notification(whether it is associated or not is decided based on message key)
is found, service writes a message to it containing relevant notification data.
Also, it has 2 HTTP endpoints for getting all notifications, and changing status of notification
from NEW to REVIEWED

## Getting Started

To run this service, just clone it, and start it 
using either [Make Run](#run) or [Make Docker Run](#run-in-docker). 
However, to run properly, it requires a separate microservice that will public messages to kafka
In scope of the [Feed Project](https://github.com/Alieksieiev0/feed-templ)
microservice called [feed-service](https://github.com/Alieksieiev0/feed-service) was used.

## MakeFile

### Build
```bash
make build
```

### Run
```bash
make run
```

### Run in docker
```bash
make docker-run
```

### Run and rebuild in docker
```bash
make docker-build-n-run
```

### Shutdown docker
```bash
make docker-down
```

### Test
```bash
make test
```

### Clean
```bash
make clean
```

### Proto
```bash
make proto
```

### Live Reload
```bash
make watch
```

## Flags
This application supports startup flags, 
that can be passed to change server and kafka urls. 
However, be careful changing notification-service server url 
if you are running it using docker-compose, because by default
only port 3002 is exposed

### Server
- Name: websocket-server
- Default: 3002

### Kafka
- Name: kafka
- Default: kafka:9094
