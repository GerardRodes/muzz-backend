# Muzz backend test
## Gerard Rodes Vidal

## Run
```sh
docker-compose up
```
My service waits for dependencies to be ready on startup so it will handle
gracefully slow startups of the databases.
On the other side `migrate/migrate` does not, so keep an eye on it and check that
it has automatically applied the migrations, if it has not, just restart it or
apply the migrations manually with:
```sh
bash scripts/migrate.sh up
```

After that feel free to create a bunch of users with something like this:
```sh
seq 1 1000 | xargs -P100 -n1 curl -XPOST http://localhost:8080/user/create
```

If you want the service to run on port `80` you can configure it with the
`MUZZ_HOST_HTTP_PORT` env var on the `.env` file and restart the service or docker container.
But you will need to provide the correct permissions to the running service to
access that port on your machine.

### Run service directly
If you want to make some changes to the src code and restart the service you have
2 options:
1. Tell docker to rebuild the image: `docker-compose up --build --force-recreate --no-deps`
2. Stop the docker `muzz` container and start the service manually, the ports to use on your host network can be configured on the `.env` file `MUZZ_HOST_*` variables.
   1. `docker-compose stop muzz`
   2. `bash scripts/dev.sh`

## Project structure
```
├── cmd
│   └── muzz            // service entrypoint
├── internal            // private packages
│   ├── config          // load service config
│   ├── domain          // domain business logic and types
│   ├── httpserver      // echo http server
│   ├── mariadb         // mariadb repo implementation
│   │   └── migrations  // sql migrations
│   └── session         // session store
└── scripts             // dev handy scripts
```


### Solution
For a larger project the standard approach would have been defining 1 controller and 1 service
per model/business concept, in fact in the first commits of the project you can find that I was
following that approach, but I decided to make things simpler as I was finding myself writing a
lot of boilerplate code.

I've also set up a migration system along the exercise showing how the development of the different
exercises has been evolving the database schema.

The service is available as a docker image with no wasted space
```
❯ docker build -t muzz:local .
❯ dive --ci --highestWastedBytes 0 --highestUserWastedPercent 0 muzz:local
  Using default CI config
Image Source: docker://muzz:local
Fetching image... (this can take a while for large images)
Analyzing image...
  efficiency: 100.0000 %
  wastedBytes: 0 bytes (0 B)
  userWastedPercent: 0.0000 %
Inefficient Files:
Count  Wasted Space  File Path
None
Results:
  PASS: highestUserWastedPercent
  PASS: highestWastedBytes
  PASS: lowestEfficiency
Result:PASS [Total:3] [Passed:3] [Failed:0] [Warn:0] [Skipped:0]
```

## Part 1 - The Basics
I've implemented everything as it is defined, the only thing to point out is that at the `/profiles`
endpoint the following is stated:
> It should return other profiles that are potential matches for this user

I've decided that "potential matches" means that they are of the opposite gender, so I've applied that filter.

## Part 2 - Authentication
Here I've decided to use Redis instead of a simpler approach like: a map in memory,
a simple file, or a new table on the relational database.

I thought about using it because it is the most common approach for session storage on a production system.
Also, I've introduced at this point [utilities to inject and extract the session](./internal/domain/context.go)
to/from the `context.Context` from the `domain` package.

## Part 3 - Filtering
The profiles can be filtered with query params
```go
type req struct {
	AgeMin       uint8         `query:"ageMin"`
	AgeMax       uint8         `query:"ageMax"`
	Gender       domain.Gender `query:"gender"`
	OrderByLikes bool          `query:"orderByLikes"`
}
```
and an extra parameter `orderByLikes` has been added to order the results by `attractiveness`
which is evaluated as `ATTRACTIVENESS = LIKES - DISLIKES` on the file [internal/mariadb/list_matches.go](./internal/mariadb/list_matches.go#L21)
as `2*SUM(s.preference) - COUNT(s.user_id) as "score"`

If `orderByLikes` is `false` the profiles will be ordered by distance from the authenticated user.

This has been implemented with a MariaDB Geometric type Point column at the `users` table, and
thanks to the function [ST_DISTANCE_SPHERE](https://mariadb.com/kb/en/st_distance_sphere/)
I am able to determine the real distance on Earth.

[I've also implemented a custom scanner to parse the Point binary representation from the database.](./internal/domain/point.go#L28)

## Next improvements
### Testing
I've added some unitary testing to the `domain` package.
The following steps would be to increase the testing coverage of the codebase and add
integration tests that ensured the correct behavior of the service with its dependencies:
MariaDB and Redis.

## Error handling and logging
Right now the error handling is very simple, but I would be needed to have a better
handling of domain errors so `httpserver` can be able to determine the correct
HTTP status code for each error.
