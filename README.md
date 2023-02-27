# Golang and Redis App
This app spins and golang app and redis in local

## Setup in local
Run in local ```sh start-local.sh```
check out server running at [http://localhost:8080/healthcheck](http://localhost:8080/healthcheck)

##
CRUD call curls
* POST ```curl --location --request DELETE 'http://localhost:8080/movies/4080f4b3-0e2b-4671-89cd-678c4d8d2a68'```
* Get ```curl --location --request GET 'http://localhost:8080/movies' --header 'Content-Type:application/json' --data-raw '{"title":"Mai hoon na","description":"srk movie"}'```
* Update curl ```curl --location --request DELETE 'http://localhost:8080/movies/4080f4b3-0e2b-4671-89cd-678c4d8d2a68'```
* Delete curl ```curl --location --request DELETE 'http://localhost:8080/movies/4080f4b3-0e2b-4671-89cd-678c4d8d2a68'```
