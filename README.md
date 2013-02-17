# ServiceDirector #

An API director, like an out-of-line load balancer, where multiple versions of various services can be registered.

Currently very very very alpha, as in not at all functional.  Plus it's my first Go code.

Features:
* Persists all info to Redis
* Allows for a cluster of read-only slaves (eg put one on each webserver to only ever read locally) which replicate from the master
* HTTP interface for getting a load-balanced API URL (on port 8081 by default, under /services/)
* ReSTful HTTP interface for updating and seeing the full current stats (on port 8082 by default)
