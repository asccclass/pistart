

build:
	GOOS=linux GOARCH=amd64 go build -o app .

s:
	git push -u origin main
