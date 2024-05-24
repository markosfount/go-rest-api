# go-rest-api

Golang REST API example. Contains endpoints for getting, creating, updating and deleting (wip) movies.

The service will get information from `The Movie Database` API when a POST request is received to
create a movie, by searching for the provided title.

Uses gorilla/mux for the api server and postgresql for the database.
