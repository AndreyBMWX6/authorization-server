# authorization-server

Authorization server used in OAuth 2.0 standard.

## run:
Install binary dependencies using:
```shell
make bin-deps
```

Create database and apply migrations.

You can start db locally or in docker container.

To start db locally run:
```shell
make db-reset
```
To start db in docker container run:
```shell
make env-up
make db-up
```
NOTE: make sure you have installed docker-compose and started your docker soket.

After db start you can run authorization server:
```shell
make run
```
