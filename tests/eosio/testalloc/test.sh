mkdir -p build
eosio-go build -gen-code=false -o build/test.wasm . || exit 1
run-uuos -m pytest -x -s test.py -k test_hello
