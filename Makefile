.PHONY: benchmark
benchmark:
	go test -cpu 1,2,4,8,16,32,64,128,256,512,1024 -bench .
