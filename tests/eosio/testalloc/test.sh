mkdir -p build
eosio-go build -o build/test.wasm . || exit 1
run-uuos -m pytest -x -s test.py -k test_hello
