all:
	/Users/danjenkins/go/bin/c-for-go ndi.yml
	mv NDI/* .
	rm -rf NDI
	go mod init NDI
	# go mod tidy

clean:
	rm -f cgo_helpers.go cgo_helpers.h cgo_helpers.c
	rm -f doc.go types.go const.go
	rm -f NDI.go
	rm -rf go.mod

test:
	cd NDI && go build