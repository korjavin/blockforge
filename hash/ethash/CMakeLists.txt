set(FILES
    util.h
    io.c
    internal.c
    ethash.h
    endian.h
    compiler.h
    fnv.h
    data_sizes.h
    sha3.c sha3.h
)

if (WIN32)
    list(APPEND FILES io_win32.c mmap_win32.c)
else()
    list(APPEND FILES io_posix.c)
endif()

add_library(ethash OBJECT ${FILES})
set_target_properties(ethash PROPERTIES C_STANDARD 99)
