## My process

Run TimescaleDB as a Docker container:

```bash
docker run -d --name timescaledb -p 5432:5432  -v ./data:/pgdata -e PGDATA=/pgdata -e POSTGRES_PASSWORD=password timescale/timescaledb-ha:pg17
```

Might need to change permissions for the `./data` directory.

```bash
mkdir data
sudo chwon 1000 data
```

Create the *cpu_usage* hypertable in the *homework* database.

> What is a hypertable?

```bash
 psql -d "postgres://postgres:password@localhost/postgres" < cpu_usage.sql
```

Populate the hypertable with data

```bash
psql -d "postgres://postgres:password@localhost/homework" -c "\COPY cpu_usage FROM cpu_usage.csv CSV HEADER"
```
