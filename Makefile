createdb:
	docker exec -it postgres16 createdb --username=root --owner=root save_plus

dropdb:
	docker exec -it postgres16 dropdb save_plus

migrateup:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/save_plus?sslmode=disable" -verbose up

migratedown:
	migrate -path db/migration -database "postgres://root:secret@localhost:5432/save_plus?sslmode=disable" -verbose down

sqlc:
	sqlc generate

test:
	go test -cover ./...

test_v:
	go test -v -cover ./...

server:
	go run main.go

.PHONY: createdb dropdb migrateup migratedown sqlc test server test_v