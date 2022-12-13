.PHONY: benchmark
benchmark:
	go test -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench .

.PHONY: benchmark-csv
benchmark-csv:
	go test -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench . | bench2csv -freq >benchmark.csv
