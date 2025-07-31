module github.com/entry

go 1.20

require github.com/common v0.0.1

replace (
	github.com/common => ../common
)

require (
	github.com/go-chi/chi/v5 v5.0.10
	github.com/google/uuid v1.3.1
	github.com/redis/go-redis/v9 v9.2.1
)

require (
	github.com/cespare/xxhash/v2 v2.2.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/gorilla/websocket v1.5.0 // indirect
)
