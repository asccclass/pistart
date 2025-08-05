

build:
	GOOS=linux GOARCH=arm go build -o app .

s:
	git push -u origin main
