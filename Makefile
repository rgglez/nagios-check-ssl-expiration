build:
	cd ./src && go build -o ../dist/check_ssl_expiration *.go

install:
	cp -v ./dist/check_ssl_expiration /usr/local/nagios/libexec/