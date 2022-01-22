#!/bin/bash
docker run -d \
    --rm -it \
    -p 5432:5432 \
    --name postgres \
    -e POSTGRES_PASSWORD=password \
    -e PGDATA=/var/lib/postgresql/data \
    -v $(pwd)/mntdata:/var/lib/postgresql/data \
    postgres:13.1
