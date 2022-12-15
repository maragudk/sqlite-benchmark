.PHONY: benchmark-postgres
benchmark-postgres:
	go test -tags postgres -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench BenchmarkDB_Postgres

.PHONY: benchmark-postgres-csv
benchmark-postgres-csv:
	go test -tags postgres -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench BenchmarkDB_Postgres | bench2csv -freq >>benchmark.csv

.PHONY: benchmark-sqlite
benchmark-sqlite:
	go test -timeout 30m -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench .

.PHONY: benchmark-sqlite-csv
benchmark-sqlite-csv:
	go test -timeout 30m -cpu 1,2,4,8,16,32,64,128,256,512 -bench . | bench2csv -freq >>benchmark.csv

.PHONY: down
down:
	docker compose down

.PHONY: up
up:
	docker compose up -d
