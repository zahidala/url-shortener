include .env

create_migration:
	migrate create -ext=sql -dir=internal/database/migrations -seq init

migrate_up:
	migrate -path=internal/database/migrations -database "${DATABASE_URL}" -verbose up

migrate_down:
	migrate -path=internal/database/migrations -database "${DATABASE_URL}" -verbose down

.PHONY: create_migration migrate_up migrate_down
