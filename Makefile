bin_name = coordinate

build:
	$(MAKE) clean
	mkdir -p bin
	go build -o bin/${bin_name} cmd/main.go

clean:
	rm -rf bin