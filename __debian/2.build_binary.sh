#!/usr/bin/env bash

PKGDIR="vclt-2.00.00-0_amd64"

mkdir -p ${PKGDIR}/opt/bin ${PKGDIR}/DEBIAN
mkdir -p ${PKGDIR}/opt/bin ${PKGDIR}/DEBIAN
for i in control preinst prerm postinst postrm;do
  mv $i ${PKGDIR}/DEBIAN/
done

echo "Building binary from source"
cd ../src
go build -o ../__debian/${PKGDIR}/opt/bin/vclt2 .
strip ../__debian/${PKGDIR}/opt/bin/vclt2
sudo chown 0:0 ../__debian/${PKGDIR}/opt/bin/vclt2

echo "Binary built. Now packaging..."
cd ../__debian/
dpkg-deb -b ${PKGDIR}
