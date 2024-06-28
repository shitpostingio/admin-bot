#!/bin/bash

ROOTDIR=$(pwd)
VERSION=$(git describe --tags --abbrev=0)
DEST=admin-bot-v$VERSION
mkdir $DEST

echo "[!!] version $VERSION"
echo "[+] building admin-bot..."
make build

mv admin-bot $DEST
cd database/cmd

for i in adminbot-*; do 
        cd $i
        echo "[+] building $i..."
        go build
        mv $i ../../../$DEST
        cd ../
done


cd $ROOTDIR
echo "[+] building tar.xz..."
tar -cf - $DEST | xz -9 -c - > $DEST.tar.xz
echo "[+] cleaning..."
rm -r $DEST
echo "[+] done!"