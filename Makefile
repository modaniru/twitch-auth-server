.PHONY: fmt
fmt:
	go fmt ./...

.PHONY: go
go: fmt
	go run src/main.go


.PHONY: migrateUp
migrateUp:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5555/postgres?sslmode=disable" up

.PHONY: migrateDown
migrateDown:
	migrate -path db/migrations -database "postgresql://postgres:postgres@localhost:5555/postgres?sslmode=disable" down

.PHONY: sqlc
sqlc:
	sqlc generate
