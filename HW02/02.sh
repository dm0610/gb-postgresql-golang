#!/bin/bash
docker run -d \
    --rm -it \
    -p 5432:5432 \
    --name postgres \
    -e POSTGRES_PASSWORD=password \
    -e PGDATA=/var/lib/postgresql/data \
    -v $(pwd)/mntdata:/var/lib/postgresql/data \
    -v $(pwd)/workdir:/workdir \
    -w /workdir \
    postgres:13.1