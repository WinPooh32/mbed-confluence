cmake_minimum_required(VERSION 3.0)
project(mbedconfluence)

set(TARGET_LIB mbedconfluence.lib)

# Go configurations.
set(GO_LIBNAME libmbedconfluence.so)

# List all .go files.
file(GLOB_RECURSE GO_SOURCES ./**.go)

add_custom_target(src SOURCES 
  "mbed_confluence.h" 
  ${GO_SOURCES}
)

set(GO_BUILD_COMMAND go build -buildmode=c-shared -o "${CMAKE_CURRENT_BINARY_DIR}/${GO_LIBNAME}" ${CMAKE_GO_FLAGS})

# Custom command for 'go build -buildmode=c-shared ...'
# to create a library from Go codes.
add_custom_command(
  OUTPUT ${GO_LIBNAME}
  DEPENDS ${GO_SOURCES}
  WORKING_DIRECTORY ${CMAKE_CURRENT_SOURCE_DIR}
  COMMAND ${GO_BUILD_COMMAND}
)

# Add a custom target for the library.
add_custom_target(${TARGET_LIB} DEPENDS ${GO_LIBNAME})

# Export target name of the library.
set(TARGET_LIB_ADD ${TARGET_LIB} PARENT_SCOPE)
