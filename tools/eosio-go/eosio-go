#!/bin/bash
if [ $1 == "build" ]; then
    tinygo build -x -gc=leaking -target eosio -wasm-abi=generic -scheduler=none -opt z ${@:2}
elif [ $1 == "init" ]; then
    tinygo init ${@:2}
else
    tinygo $@
fi
