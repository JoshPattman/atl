build:
	mkdir -p ./bin
	go build -o bin/atlgen .

install: build
	sudo cp bin/atlgen /usr/local/bin

local-release:
	goreleaser release --snapshot --skip publish --clean

local-release-install: local-release
	sudo dpkg -i dist/*_linux_amd64.deb
