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

Define a SQL query that "returns the max cpu usage and min cpu usage of the given hostname for every minute in the time range specified by the start time and end time". This is an example using the data from the first line in *query_params.csv*.

The native *date_bin* function could potentially be used but I preferred TimescaleDB's own *time_bucket*.

```SQL
SELECT
  time_bucket('1 minute', ts, '2017-01-01 08:59:22'::timestamp) AS minute,
  MAX(usage) AS max_usage,
  MIN(usage) AS min_usage
FROM cpu_usage
WHERE host = 'host_000008'
  AND ts >= '2017-01-01 08:59:22'
  AND ts < '2017-01-01 10:59:22'
GROUP BY minute
ORDER BY minute;
```

81.93 |     56.61