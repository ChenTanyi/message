#!/bin/sh

set -ex

curdir=`pwd`

cd $(dirname $0)

if [ -d "template" ]; then
    go get -v github.com/go-bindata/go-bindata/go-bindata
    go-bindata template/
fi

cd $curdir