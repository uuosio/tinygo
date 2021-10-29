tinygo build -x -gc=leaking -target eosio -wasm-abi=generic -scheduler=none -opt z -tags=math_big_pure_go -gen-code=true -strip=true -o test.wasm .
