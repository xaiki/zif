#!/usr/bin/env bash

ZIF_BIN="$GOPATH/bin/zifd"
VERBF=""

if [ "$1" == "-v" ]; then
    VERBF="-v -x"
fi

go install $VERB

pushd zifd
go install $VERB
popd

which npm >/dev/null 2>&1
if [ $? -eq 0 ]; then
    pushd ui
    npm install
    popd
fi

if [ -d "/proc/sys/kernel/pax" ]; then
    paxctl -c "$ZIF_BIN" && setfattr -n user.pax.flags -v "emr" "$ZIF_BIN"
fi

