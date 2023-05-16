DB_URL="postgresql://root:secret@localhost:5432/simple_bank?sslmode=disable"
postgres:
	docker run --name postgres-latest -p 5432:5432 -e POSTGRES_USER=root -e POSTGRES_PASSWORD=secret -d postgres 
createdb: 
	docker exec -it postgres-latest createdb --username=root --owner=root simple_bank
dropdb:
	docker exec -it postgres-latest dropdb --username=root simple_bank
migrateup:
	migrate -path db/migration -database $(DB_URL) -verbose up  
migratedown:
	migrate -path db/migration -database $(DB_URL) -verbose down  

sqlc:
	sqlc generate 