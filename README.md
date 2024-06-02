# oauth2

Implementation of OAuth 2.0 standard. Authrization server, simple client and resource server examples.

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

After db start you can run authorization server and client:
```shell
make run
```

You also can run them separately by using:
```shell
make run-auth
```
and 
```shell
make run-client
```
