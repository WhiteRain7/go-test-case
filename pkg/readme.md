# Project/pkg/

Infrastructure and utils used shared between services.

# How to run

## Docker run

No additional steps required.

## Local run

No additional steps required.

# Project structure

- `infrastructure/`
  - `config/` - reads service .env with viper, expected variables described explicitly in service `config/`
  - `broker/` - kafka broker
    - `consumer.go` - reads messages from broker
    - `producer.go` - writes messages to broker
  - `logger/` - async logger to file and console
  - `postgres/` - psql client
  - `protoclient/` - gRPC client (server provided by gw-exchanger)
  - `redis/` - redis cache
- `utils/`
  - `coded-errors/` - error wrap to provide HTTP codes and messages to handlers

# environment

Each infrastructure provides Config struct, describing necessary env variables.
If service use some infrastructure, it also imports config in own `config/` directory
to describe what to parse from .env.

## broker

The only one exports two configs:

### ConConfig

```dotenv
host=kafka # docker container name provides broker, set localhost for local run
port=9092 # broker port
topic=abc # source of messages, topic provided by producer
group=abc # messages distributed between several consumers with the same group
```

### ProdConfig

```dotenv
host=kafka # docker container name provides broker, set localhost for local run
port=9092 # broker port
topic=abc # where to place messages
```

## logger

```dotenv
path=/logs/main.log # file to append logging messages
console=1 # 0/1 to turn off/on additional output to console
```

## postgres

```dotenv
user=postgres # psql user to connect
password=postgres # psql user password to connect
host=wallets_db # docker container name or localhost
port=5432 # psql port
database=wallets # name of db
SSLMode=0 # 0/1 to turn off/on SSL mode
```

## protoclient

```dotenv
host=exchanger # docker container name or localhost
port=9001 # grpc port
```

## redis

```dotenv
host=redis # docker container name or localhost
port=6379 # redis port
```