# eosio-cpp -o test-cpp.wasm cpp/test.cpp
# if [ $? -ne 0 ]; then
#     echo "build cpp failed"
#     exit $?
# fi

# eosio-wasm2wast test.wasm -o test.wast

eosio-go build -o test.wasm .
if [ $? -ne 0 ]; then
    echo "build go failed"
    exit $?
fi

if [ -z "$1" ]; then
run-ipyeos -m pytest -x -s test.py
else
run-ipyeos -m pytest -x -s test.py -k $1
fi
