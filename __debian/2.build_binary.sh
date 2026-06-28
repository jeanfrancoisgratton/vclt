#!/usr/bin/env bash

PKGDIR="vclt-2.3.0-1_amd64"

mkdir -p ${PKGDIR}/opt/bin ${PKGDIR}/DEBIAN
mkdir -p ${PKGDIR}/opt/bin ${PKGDIR}/DEBIAN
for i in control preinst prerm postinst postrm;do
  mv $i ${PKGDIR}/DEBIAN/
done

echo "Building binary from source"
cd ../src
go build -trimpath -mod=readonly -modcacherw -ldflags="-s -w -buildid=" -o vclt .
strip ../__debian/${PKGDIR}/opt/bin/vclt
sudo chown 0:0 ../__debian/${PKGDIR}/opt/bin/vclt

echo "Binary built. Now packaging..."
cd ../__debian/
dpkg-deb -b ${PKGDIR}
