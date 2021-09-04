eosio-go build -o test.wasm .
if [ $? -ne 0 ]; then
    echo "build failed"
    exit $?
fi

if [ -z "$1" ]; then
run-uuos -m pytest -x -s test.py
else
run-uuos -m pytest -x -s test.py -k $1
fi
