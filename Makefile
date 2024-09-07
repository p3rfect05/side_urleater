run:
	go run cmd/main.go

migrate:
	migrate -database "postgres://postgres:postgres@localhost:5432/postgres?sslmode=disable" -path build/migrations up