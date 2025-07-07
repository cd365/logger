module main

go 1.23.0

require (
	github.com/cd365/logger/v9 v9.0.0
	github.com/redis/go-redis/v9 v9.11.0
	github.com/rs/zerolog v1.34.0
)

require (
	github.com/cespare/xxhash/v2 v2.3.0 // indirect
	github.com/dgryski/go-rendezvous v0.0.0-20200823014737-9f7001d12a5f // indirect
	github.com/mattn/go-colorable v0.1.14 // indirect
	github.com/mattn/go-isatty v0.0.20 // indirect
	golang.org/x/sys v0.33.0 // indirect
)

replace github.com/cd365/logger/v9 v9.0.0 => ../../logger
