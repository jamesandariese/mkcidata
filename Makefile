.PHONY: all always-build

RELPATH := $(shell realpath .)

all: mkcidata mkcidata.win32 version.txt
always-build:

version.txt: version.txt.tmpl always-build
	gomplate < version.txt.tmpl > version.txt
	git add version.txt

prepare-release: version.txt

test:
	echo "Add tests" 1>&2
	false

mkcidata: $(wildcard *.go **/*.go)
	go build -o mkcidata

mkcidata.win32: $(wildcard *.go **/*.go)
	GOOS=windows go build -o mkcidata.win32

clean:
	rm -f mkcidata mkcidata.win32
