mkdir -p build
pushd build
CXX=$(pwd)/../../../llvm-build/bin/clang C=$(pwd)/../../../llvm-build/bin/clang cmake .. -DLLVM_BUILD_DIR=$(pwd)/../../../llvm-build
make -j4
popd
