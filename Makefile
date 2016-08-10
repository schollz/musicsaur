SOURCEDIR=.
SOURCES := $(shell find $(SOURCEDIR) -name '*.go')

BINARY=musicsaur

VERSION=1.5
BUILD_TIME=`date +%FT%T%z`
BUILD=`git rev-parse HEAD`

LDFLAGS=-ldflags "-X main.VersionNum=${VERSION} -X main.Build=${BUILD} -X main.BuildTime=${BUILD_TIME}"

.DEFAULT_GOAL: $(BINARY)

$(BINARY): $(SOURCES)
	go get github.com/mholt/caddy/caddyhttp
	go get github.com/BurntSushi/toml
	go get github.com/mholt/caddy
	go get github.com/mholt/caddy/caddytls
	go get github.com/toqueteos/webbrowser
	go get github.com/bobertlo/go-id3/id3
	go get github.com/tcolgate/mp3
	go get gopkg.in/tylerb/graceful.v1
	go build ${LDFLAGS} -o ${BINARY} ${SOURCES}

.PHONY: install
install:
	go install ${LDFLAGS} ./...

.PHONY: clean
clean:
	if [ -f ${BINARY} ] ; then rm ${BINARY} ; fi
	rm -rf builds
	rm -rf musicsaur*

.PHONY: binaries
binaries:
	rm -rf builds
	mkdir builds
	# Build Windows
	env GOOS=windows GOARCH=amd64 go build ${LDFLAGS} -o musicsaur.exe -v *.go
	zip -r musicsaur_${VERSION}_windows_amd64.zip musicsaur.exe LICENSE ./templates/* ./static/*
	mv musicsaur_${VERSION}_windows_amd64.zip builds/
	rm musicsaur.exe
	# Build Linux
	env GOOS=linux GOARCH=amd64 go build ${LDFLAGS} -o musicsaur -v *.go
	zip -r musicsaur_${VERSION}_linux_amd64.zip musicsaur LICENSE ./templates/* ./static/*
	mv musicsaur_${VERSION}_linux_amd64.zip builds/
	rm musicsaur
	# Build OS X
	env GOOS=darwin GOARCH=amd64 go build ${LDFLAGS} -o musicsaur -v *.go
	zip -r musicsaur_${VERSION}_osx.zip musicsaur LICENSE ./templates/* ./static/*
	mv musicsaur_${VERSION}_osx.zip builds/
	rm musicsaur
	# Build Raspberry Pi / Chromebook
	env GOOS=linux GOARCH=arm go build ${LDFLAGS} -o musicsaur -v *.go
	zip -r musicsaur_${VERSION}_linux_arm.zip musicsaur LICENSE ./templates/* ./static/*
	mv musicsaur_${VERSION}_linux_arm.zip builds/
	rm musicsaur
