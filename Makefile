.PHONY: bin
bin:
	go build -o bin/dlux

build: static
	go build -o bin/dlux

static: staticsetup
	cd static && npm run build

staticsetup:
	cd static && npm ci

base:
	lxc image alias delete dlxbase || echo No base image to delete
	sudo MYUSER=bjk distrobuilder build-lxd contrib/distrobuilder/debian.yaml --import-into-lxd=dlxbase
	sudo rm rootfs.squashfs
	sudo rm lxd.tar.xz

clean:
	sudo rm rootfs.squashfs
	sudo rm lxd.tar.xz


launch:
	lxc rm --force tester
	lxc launch dlxbase tester
	sleep 3
	lxc login tester

aliases:
	lxc alias list
