containers:
	docker compose up -d
createdb:
	docker exec -it dotshop_db createdb dotshop --username=root --owner=root dotshop
dropdb:
	docker exec -it dotshop_db dropdb dotshop
migrate-create:
	migrate create -ext sql -dir migrations -seq init_schema
migrateup:
	migrate -path migrations -database "postgresql://root:rand0mPassword@localhost:5432/dotshop?sslmode=disable" -verbose up
migratedown:
	migrate -path migrations -database "postgresql://root:rand0mPassword@localhost:5432/dotshop?sslmode=disable" -verbose down
run:
	go run cmd/main.go
gen-swagger:
	swag fmt && bash scripts/gen-swagger.sh

.PHONY: containers createdb dropdb migrate-create migrateup migratedown run gen-swagger
