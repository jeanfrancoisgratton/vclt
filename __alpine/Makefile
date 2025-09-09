all: 
	cd ../src
	CGO_ENABLED=0 go build . -o ./vclt2.exe

prepare:
	sudo /opt/bin/install_golang.sh 1.24.1 amd64
	sudo apk update --no-cache

install:
	test -d $(DESTDIR)$(PREFIX)/bin || mkdir -p $(DESTDIR)$(PREFIX)/bin
	install -pm 755 $(BUILDDIR)/vclt2.exe $(DESTDIR)$(PREFIX)/bin/vclt2
	strip $(DESTDIR)$(PREFIX)/bin/vclt2
