#!/bin/sh

mkdir -p package/
cd ../xj9/ || exit 1
echo "... building"
go build || exit 1
cd ../installer/
echo "... compressing files"
tar cfJ package/package.tar.xz ../xj9/xj9 ../xj9/images/*.png
