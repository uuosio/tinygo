tinygo build -x -gc=leaking -o test.wasm -target eosio -wasm-abi=generic -scheduler=none  -opt z -tags=math_big_pure_go .
if [ $? -ne 0 ]; then
    echo "build failed"
    exit $?
fi

if [ -z "$1" ]; then
run-uuos -m pytest -x -v test.py
else
run-uuos -m pytest -x -v test.py -k $1
fi
