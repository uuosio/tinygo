cmake_minimum_required(VERSION 3.4.3)
# hack for OSX
set(CMAKE_OSX_SYSROOT="${CMAKE_BINARY_DIR}")
set(CMAKE_OSX_DEPLOYMENT_TARGET="")
set(CMAKE_SYSROOT "${CMAKE_BINARY_DIR}")

project(eosio)

file(MAKE_DIRECTORY ${CMAKE_BINARY_DIR}/lib)

set(CMAKE_SYSTEM_NAME Generic)
set(CMAKE_SYSTEM_PROCESSOR wasm)
set(CMAKE_CROSSCOMPILING 1)
set(CMAKE_C_COMPILER_WORKS 1)
set(CMAKE_CXX_COMPILER_WORKS 1)

#include(${CMAKE_SOURCE_DIR}/modules/WasmToolchain.cmake)

set(CMAKE_C_COMPILER "${LLVM_BUILD_DIR}/bin/clang")
set(CMAKE_CXX_COMPILER "${LLVM_BUILD_DIR}/bin/clang")
set(CMAKE_ASM_COMPILER "${LLVM_BUILD_DIR}/bin/clang")

set(CMAKE_C_FLAGS "--target=wasm32 -ffreestanding -nostdlib -fno-builtin -fno-threadsafe-statics -fno-exceptions -fno-rtti -fmodules-ts")
set(CMAKE_CXX_FLAGS "${CMAKE_C_FLAGS} -DBOOST_DISABLE_ASSERTS -DBOOST_EXCEPTION_DISABLE -mllvm -use-cfl-aa-in-codegen=both -O3 --std=c++17")
set(CMAKE_ASM_FLAGS " -fnative -fasm ")
set(CMAKE_AR "${LLVM_BUILD_DIR}/bin/llvm-ar")
set(CMAKE_RANLIB "${LLVM_BUILD_DIR}/bin/llvm-ranlib")

set(WASM_LINKER "${LLVM_BUILD_DIR}/bin/wasm-ld")

set(CMAKE_C_LINK_EXECUTABLE "${WASM_LINKER} <LINK_FLAGS> <OBJECTS> -o <TARGET> <LINK_LIBRARIES>")
set(CMAKE_CXX_LINK_EXECUTABLE "${WASM_LINKER} <LINK_FLAGS> <OBJECTS> -o <TARGET> <LINK_LIBRARIES>")

set(CMAKE_FIND_ROOT_PATH "${CMAKE_BINARY_DIR}")
set(CMAKE_FIND_ROOT_PATH_MODE_PROGRAM NEVER)
set(CMAKE_FIND_ROOT_PATH_MODE_LIBRARY ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_INCLUDE ONLY)
set(CMAKE_FIND_ROOT_PATH_MODE_PACKAGE ONLY)


set(BASE_BINARY_DIR "${CMAKE_BINARY_DIR}")

add_subdirectory(libc)
add_subdirectory(libc++)
add_subdirectory(rt)
add_subdirectory(eosiolib)
add_subdirectory(boost)
