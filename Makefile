build:
	cd servo; GOARM=6 GOOS=linux GOARCH=arm go build -o ../robot .
	cd controller; go build -o ../ct .

build-debug:
	cd servo; go build -o ../robot .
	cd controller; go build -o ../ct .

