#tinygo build -x -gc=leaking -target eosio -wasm-abi=generic -scheduler=none -opt 0 -tags=math_big_pure_go -gen-code=true -strip=false -o test.wasm .
mkdir -p build
SYSROOT=/Users/newworld/dev/github/uuosio.gscdk/tinygo/lib/eosio/sysroot
tinygo clang --target=wasm32--wasi --sysroot=$SYSROOT -Oz -I$SYSROOT/include -I$SYSROOT/include/libc -I$SYSROOT/include/libcxx -I$SYSROOT/include/eosiolib/capi -I$SYSROOT/include/eosiolib/core -I$SYSROOT/include/eosiolib/contracts -g -I$(pwd) -MD -MV -MTdeps -Xclang -internal-isystem -Xclang $SYSROOT/include/libc  -c -std=c++17 -Wno-unknown-attributes -o build/test.o test.cpp || exit 1
tinygo build -x -gc=leaking -target eosio -wasm-abi=generic -scheduler=none -opt z -tags=math_big_pure_go -gen-code=true -strip=true -o build/test.wasm . || exit 1
