file(GLOB HEADERS "*.hpp"
                  "*.h")

set(INCLUDES "${CMAKE_CURRENT_SOURCE_DIR}/../libc/musl/include"
    ${CMAKE_CURRENT_SOURCE_DIR}/../libc/musl/src/internal
    ${CMAKE_CURRENT_SOURCE_DIR}/../libc/musl/src/crypt
    ${CMAKE_CURRENT_SOURCE_DIR}/../libc/musl/arch/eos
    "${CMAKE_CURRENT_SOURCE_DIR}/../libc++/libcxx/include"
    "${CMAKE_CURRENT_SOURCE_DIR}/../boost/include"
)

add_library(eosio
            eosiolib.cpp
            crypto.cpp
            ${HEADERS})

add_library(eosio_malloc
            malloc.cpp
            ${HEADERS})
target_include_directories(eosio_malloc PUBLIC ${INCLUDES})

add_library(eosio_dsm
            simple_malloc.cpp
            ${HEADERS})
target_include_directories(eosio_dsm PUBLIC ${INCLUDES})

add_library(eosio_cmem
            memory.cpp
            ${HEADERS})
message(STATUS "Adding eosio_cmem ${CMAKE_SOURCE_DIR}/libc++/libcxx/include")
target_include_directories(eosio_cmem PUBLIC
    ${INCLUDES}
    "${CMAKE_SOURCE_DIR}/libc++/libcxx/include"
)

set_target_properties(eosio_malloc PROPERTIES LINKER_LANGUAGE C)

target_include_directories(eosio PUBLIC ${INCLUDES}})

target_link_libraries( eosio c c++ )

add_custom_command( TARGET eosio POST_BUILD COMMAND ${CMAKE_COMMAND} -E copy $<TARGET_FILE:eosio> ${BASE_BINARY_DIR}/lib )
add_custom_command( TARGET eosio_malloc POST_BUILD COMMAND ${CMAKE_COMMAND} -E copy $<TARGET_FILE:eosio_malloc> ${BASE_BINARY_DIR}/lib )
add_custom_command( TARGET eosio_dsm POST_BUILD COMMAND ${CMAKE_COMMAND} -E copy $<TARGET_FILE:eosio_dsm> ${BASE_BINARY_DIR}/lib )
add_custom_command( TARGET eosio_cmem POST_BUILD COMMAND ${CMAKE_COMMAND} -E copy $<TARGET_FILE:eosio_cmem> ${BASE_BINARY_DIR}/lib )

file(COPY ${CMAKE_CURRENT_SOURCE_DIR}/../eosiolib DESTINATION ${BASE_BINARY_DIR}/include FILES_MATCHING PATTERN "*.h" PATTERN "*.hpp")
