module github.com/Infamous003/go-blog

go 1.24.5

require github.com/go-chi/chi/v5 v5.2.3 // direct

require github.com/lib/pq v1.10.9 // direct

require (
	github.com/tomasen/realip v0.0.0-20180522021738-f0c99a92ddce // direct
	golang.org/x/time v0.14.0 // direct
)

require (
	github.com/wneessen/go-mail v0.7.2 // direct
	golang.org/x/crypto v0.45.0 // direct
	golang.org/x/text v0.31.0 // indirect
)

require (
	github.com/BurntSushi/toml v1.4.1-0.20240526193622-a339e1f7089c // indirect
	golang.org/x/exp/typeparams v0.0.0-20231108232855-2478ac86f678 // indirect
	golang.org/x/mod v0.29.0 // indirect
	golang.org/x/sync v0.18.0 // indirect
	golang.org/x/tools v0.38.0 // indirect
	honnef.co/go/tools v0.6.1 // indirect
)

tool honnef.co/go/tools/cmd/staticcheck
