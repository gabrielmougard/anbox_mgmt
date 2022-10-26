# Anbox cloud management


## How to

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

(Check that you have a PostgreSQL Docker container running : `docker ps | grep "anboxcloud_postgresql_server"`)