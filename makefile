test:
	go test -coverprofile cover.out ./...
test_verbose:
	go test  -v -coverprofile cover.out ./...
test_send:
	cat cover.out | ssh localhost -p 2222 msg="hello" expired="10s"
test_download:
	sleep 1
	wget --delete-after --timeout 1 --tries 1 --no-verbose http://localhost:3000/api/v1/transfer/sample
test_flow:
	make -j2 test_send test_download
run:
	go run main.go
