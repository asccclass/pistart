

build:
	GOOS=linux GOARCH=amd64 go build .

s:
	git push -u origin main
