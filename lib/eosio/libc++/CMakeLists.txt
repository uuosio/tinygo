SET(SRC_FILENAMES algorithm.cpp any.cpp bind.cpp condition_variable.cpp functional.cpp
                  future.cpp ios.cpp iostream.cpp locale.cpp memory.cpp mutex.cpp new.cpp optional.cpp
                  regex.cpp stdexcept.cpp string.cpp strstream.cpp system_error.cpp exception.cpp
                  typeinfo.cpp utility.cpp valarray.cpp variant.cpp vector.cpp eosio.cpp)

SET(SRC_FILES "")
FOREACH(FN ${SRC_FILENAMES})
   LIST(APPEND SRC_FILES "libcxx/src/${FN}")
ENDFOREACH(FN)

set(CMAKE_CXX_FLAGS "${CMAKE_CXX_FLAGS} -Wno-everything  -Wno-unknown-attributes")

add_library(c++
            ${SRC_FILES})

# Let the library headers know they are currently being used to build the
# library.

target_compile_definitions(c++ PRIVATE -D_LIBCPP_BUILDING_LIBRARY)

target_include_directories(c++
                           PUBLIC 
                           ${CMAKE_CURRENT_SOURCE_DIR}/libcxx/include
                           ${CMAKE_SOURCE_DIR}/libc/musl/include
                           ${CMAKE_SOURCE_DIR}/libc/musl/src/internal
                           ${CMAKE_SOURCE_DIR}/libc/musl/src/crypt
                           ${CMAKE_SOURCE_DIR}/libc/musl/arch/eos)

target_link_libraries(c++ c)

add_custom_command( TARGET c++ POST_BUILD COMMAND ${CMAKE_COMMAND} -E copy $<TARGET_FILE:c++> ${CMAKE_BINARY_DIR}/lib/ )

file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/libcxx/include/ DESTINATION ${CMAKE_BINARY_DIR}/include/libcxx)
