# DistributedDatabaseFromScratch

I dont know what I am doing 

## Prerequisites

- Postgres database with three empty databases named Raft1, Raft2, Raft3

## QuickStartup 

Build

`go build`

Run three instances

`./main.exe --raft-id 1 --http-port 8081 --raft-port 8881 --db-name Raft1`
`./main.exe --raft-id 2 --http-port 8082 --raft-port 8882 --db-name Raft2`
`./main.exe --raft-id 3 --http-port 8083 --raft-port 8883 --db-name Raft3`


Add first follower

`curl http://localhost:8081/join?followerAddr=localhost:8882&followerId=2`

Add second follower

`curl http://localhost:8081/join?followerAddr=localhost:8883&followerId=3`

Add values

`curl -X POST http://127.0.0.1:8081/value -d '{"key": "x", "value": "23"}'`