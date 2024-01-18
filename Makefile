run-container: 
	docker run --name segment3d-db -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=postgres -d postgres:12.17-alpine3.19

db-up: 
	docker start segment3d-db

db-down: 
	docker stop segment3d-db

create-db: 
	docker exec -it segment3d-db createdb --username=root --owner=root segment3d_db

drop-db: 
	docker exec -it segment3d-db dropdb -U root segment3d_db

migrate-up: 
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/segment3d_db?sslmode=disable" --verbose up

migrate-down:
	migrate -path db/migration -database "postgresql://root:postgres@localhost:5432/segment3d_db?sslmode=disable" --verbose down

sqlc:
	sqlc generate

server-dev:
	air

server-prod:
	go run main.go

test:
	go test ./... -cover -v