migrate-create:
	sh -c "migrate create -ext sql -dir configs/migrations -seq $(name)"

migrate-up:
	sh -c "migrate -path configs/migrations -database postgres://api:pwd@localhost:5433/todo?sslmode=disable up"

test:
	go test -p 1000 ./...
