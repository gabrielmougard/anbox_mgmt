# Anbox cloud management


## Prerequisite : Load seed database

* First, download the tools :

```
cd scripts/postgresql && \
    make psql-add-postgres-ubuntu-client && \
    make psql-add-migrate-ubuntu-client
```

* Then, launch the database :

```
cd scripts/postgresql && \
    make psql-run
```

* Apply the migrations :

```
cd scripts/postgresql && \
    make psql-migrate-up
```

(Check that you have a PostgreSQL Docker container running : `docker ps | grep "anboxcloud_postgresql_server"`)

* (Optional) Seed the database :

```
cd scripts/postgresql && \
    make psql-import-db
```

## Building the server and the CLI

* To build the binaries, simply do :

```
make build
```

* To run the server, do :

```
make run
```

* The CLI client to interact with the server is at `bin/anbox-cli` (**use the full `./bin/anbox-client` path when executing, else some ENV variables won't be declared and the client will panic**):

```
$ ./bin/anbox-cli --help

Managing the Anbox application. We can do CRUD operations on "users" and "games" and also create links between entities

Usage:
  anbox-cli [command]

Available Commands:
  completion  Generate the autocompletion script for the specified shell
  create      Create entities
  delete      Delete entities
  help        Help about any command
  link        Link entities
  list        List entities
  login       Login to a user account
  update      Update entities

Flags:
  -h, --help   help for anbox-cli

Use "anbox-cli [command] --help" for more information about a command.
```

## OpenAPI/Postman integration

* Once you have a running server, you can import the **OpenAPI** specification at `pkg/api/open-api.yml` in Postman to directly interact with the server through the beautiful Postman UI. This solution is ideal to see the API  documentation and use the prefilled Postman queries for each API calls.