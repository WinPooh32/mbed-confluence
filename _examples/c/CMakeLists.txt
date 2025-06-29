cmake_minimum_required(VERSION 3.0)
project(app)

set(SOURCES main.c)

set(LIB_MBEDCONFLUENCE_DIR  ../../)
set(LIB_MBEDCONFLUENCE_BINARY_DIR "mbdebconfluence")

# Include mbdebconfluence project.
add_subdirectory(${LIB_MBEDCONFLUENCE_DIR} ${LIB_MBEDCONFLUENCE_BINARY_DIR})

# A library and a header are generated in the binary directory.
include_directories(
    ${LIB_MBEDCONFLUENCE_DIR}
    ${CMAKE_CURRENT_BINARY_DIR}/${LIB_MBEDCONFLUENCE_BINARY_DIR}
)

link_directories(${CMAKE_CURRENT_BINARY_DIR}/${LIB_MBEDCONFLUENCE_BINARY_DIR})

add_executable(app ${SOURCES})
add_dependencies(app ${TARGET_LIB_ADD})
target_link_libraries(app mbedconfluence pthread)