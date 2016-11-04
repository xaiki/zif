#!/usr/bin/env bash

ZIF_BIN="$GOPATH/bin/zifd"

go install
pushd zifd; go install; popd
pushd ui;  npm install; popd

if [ -d "/proc/sys/kernel/pax" ]; then
    paxctl -c "$ZIF_BIN" && setfattr -n user.pax.flags -v "emr" "$ZIF_BIN"
fi

