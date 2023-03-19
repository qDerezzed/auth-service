set -o allexport; source .env; set +o allexport

goose -dir ./migrations postgres "postgres://postgres:${POSTGRES_PASSWORD}@172.18.0.1:5432/postgres?sslmode=disable" status

goose -dir ./migrations postgres "postgres://postgres:${POSTGRES_PASSWORD}@172.18.0.1:5432/postgres?sslmode=disable" down
goose -dir ./migrations postgres "postgres://postgres:${POSTGRES_PASSWORD}@172.18.0.1:5432/postgres?sslmode=disable" up