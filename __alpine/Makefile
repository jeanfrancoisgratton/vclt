all: 
	cd ../src
	CGO_ENABLED=0 go build . -o ./vclt.exe

prepare:
	sudo /opt/bin/install_golang.sh 1.23.0 amd64
	sudo apk update --no-cache

install:
	test -d $(DESTDIR)$(PREFIX)/bin || mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -pm 755 $(BUILDDIR)/vclt.exe $(DESTDIR)$(PREFIX)/bin/vclt
	strip $(DESTDIR)$(PREFIX)/bin/vclt
