set(CMAKE_C_COMPILER "${LLVM_BUILD_DIR}/bin/clang")
set(CMAKE_CXX_COMPILER "${LLVM_BUILD_DIR}/bin/clang++")
set(CMAKE_ASM_COMPILER "${LLVM_BUILD_DIR}/bin/clang")

set(CMAKE_C_FLAGS "--target=wasm32 -ffreestanding -nostdlib -fno-builtin -fno-threadsafe-statics -fno-exceptions -fno-rtti -fmodules-ts")
set(CMAKE_CXX_FLAGS "${CMAKE_C_FLAGS} -DBOOST_DISABLE_ASSERTS -DBOOST_EXCEPTION_DISABLE -mllvm -use-cfl-aa-in-codegen=both -O3 --std=c++17")
set(CMAKE_ASM_FLAGS " -fnative -fasm ")
set(CMAKE_AR "${LLVM_BUILD_DIR}/bin/llvm-ar")
set(CMAKE_RANLIB "${LLVM_BUILD_DIR}/bin/llvm-ranlib")

set(WASM_LINKER "${LLVM_BUILD_DIR}/bin/wasm-ld")

set(CMAKE_C_LINK_EXECUTABLE "${WASM_LINKER} <LINK_FLAGS> <OBJECTS> -o <TARGET> <LINK_LIBRARIES>")
set(CMAKE_CXX_LINK_EXECUTABLE "${WASM_LINKER} <LINK_FLAGS> <OBJECTS> -o <TARGET> <LINK_LIBRARIES>")
