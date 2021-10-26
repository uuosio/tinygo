mkdir -p build
pushd build
cmake .. -DLLVM_BUILD_DIR=$(pwd)/../../../llvm-build -DCMAKE_C_COMPILER_FORCED=TRUE -DCMAKE_CXX_COMPILER_FORCED=TRUE
make -j4
popd

mkdir -p sysroot || exit 1

mkdir -p sysroot/lib
cp -r build/lib sysroot || exit 1

mkdir -p sysroot/include || exit 1
cp -r build/include/* sysroot/include
