# timescale-benchmarks

A Go CLI application with the goal of benchmarking TimescaleDB *SELECT* performance.

## Running the benchmarker

The easiest way to see the benchmarker in action:

```bash
docker-compose up
```

Running this command relies on *docker* and *docker-compose*. Make sure you have both properly set up before running it.

This command will:
- Bring up *TimescaleDB* in a container. This includes running the migrations and adding the Timescale extension, which is done automatically.
- Build a Docker image for the benchmarker application if you don't have one already.
- Launch the benchmarker application in its own container. By default, it runs with WORKERS=2, FILE=./db/query_params.csv and DEBUG=false. You can customize that by setting up environment variables when running the command. For instance:

```bash
WORKERS=7 FILE=./db/query_params_long.csv DEBUG=true docker-compose up
```

Alternatively, if you do not want to use Docker, simply build the application:

```bash
go build -o benchmarker ./cmd/benchmarker
```

And then interact with the binary you built, remembering to pass the proper DB_URL environment variable. Example:

```bash
DB_ENV="postgres://postgres:password@localhost:5432/homework" ./benchmarker --workers=3 --file=./db/query_params.csv
```

When interacting with the binary, you can also choose to pipe the CSV file:

```bash
cat ./db/query_params.csv | DB_ENV="postgres://postgres:password@localhost:5432/homework" ./benchmarker --workers=3
```

## Benchmarker architecture

The following diagram illustrates the data flow and architecture of the benchmarker application:

![Architecture Diagram](assets/architecture-diagram.png)

The flow starts with parsing the CSV file which contains the list of query parameters. This is done by a dedicated Goroutine running the *Stream* function, defined in the *csvreader* package.

The name *Stream* reveals something about the behavior of this function. It doesn't consume the entire CSV file at once, instead processing it in a line by line manner. This makes sure the *benchmarker* tool can handle large files with no issues, as the full contents are not loaded into memory at once. Each line in the file is parsed as a QueryData struct, containing *Hostname*, *StartTime* and *EndTime* fields. These structs are then pushed into different channels. The program creates one channel for each Worker, and decides to which channel a QueryData should be pushed to based on the *Hostname* field, which is hashed using FNV-1a and modulo operation to distribute work evenly across Workers. This also guarantees that queries with the same Hostname are only consumed by a fixed Worker.

Once inserted in these channels, the QueryData structs are ready to be consumed by a Worker from the WorkerPool. The size of the pool is defined by the *--workers* flag and defaults to 2 (Note that when running the program via *docker-compose up*, an ENV should be passed instead of a flag). Each Workers listens for items in a given channel and processes them in order.

In more detail, the work done by a Worker every time it picks up a QueryData object is essentially acquiring a connection to the database via a Connection Pool which is shared by all Workers (configured with MaxConns = Workers + 2, MinConns = Workers/2 + 1) and then leveraging that connection to run the following query against the database (using the parameters from the QueryData, of course):

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

While the query runs, the Worker processes the time it takes for the database to actually process the query. This time is the result we are actually interested in, so the Worker appends it to a *Results* array in a concurrency-safe manner.

> You can set *--debug=true* (or DEBUG=true if running via docker-compose) to see individual logs for each processed query. This is optional and, as the flag suggests, for debugging purposes.

Once all Workers finish all of their Jobs, the WorkerPool entity will calculate the final query execution metrics from the Results array and those will be printed to STDOUT.

## Testing

The application includes unit tests. Database interactions are tested via mocks, which are made possible by using the Repository pattern (abstracting away the storage interactions to make they are not tightly coupled to our actual Database). More comprehensive tests would be achieved using Integration Tests but that was beyond the initial scope.
