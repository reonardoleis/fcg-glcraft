run:
	go run main.go

build-win:
	go build -o fcg-glcraft.exe main.go 

build-linux:
	go build -o fcg-glcraft main.go