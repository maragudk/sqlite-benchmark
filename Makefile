.PHONY: benchmark
benchmark: up
	go test -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench .
	make down

.PHONY: benchmark-csv
benchmark-csv: up
	go test -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench . | bench2csv -freq >benchmark.csv
	make down

.PHONY: down
down:
	docker compose down

.PHONY: up
up:
	docker compose up -d
