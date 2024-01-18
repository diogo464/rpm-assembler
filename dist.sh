#!/bin/sh

if [ "$VERSION" = "" ]; then
	echo "VERSION not specified"
	exit 1
fi

mkdir -p dist
go build . || exit 1
rpm-assembler \
	--name rpm-assembler \
	--summary "Assemble rpm packages from artifacts" \
	--version $VERSION \
	--arch x86_64 \
	--url https://git.d464.sh/code/rpm-assembler \
	rpm-assembler:/usr/bin/rpm-assembler:0755

mv *.rpm dist
mv rpm-assembler dist/rpm-assembler_linux_x86-64
