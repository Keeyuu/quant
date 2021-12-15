# CMakeList.txt: cpp 的 CMake 项目，在此处包括源代码并定义
# 项目特定的逻辑。
#
cmake_minimum_required (VERSION 3.8)

project ("cpp")

set(SOURCES
    src/lib.cpp
)

# 将源代码添加到此项目的可执行文件。
add_executable (cpp ${SOURCES})

target_include_directories(cpp
    PRIVATE
        ${PROJECT_SOURCE_DIR}/include
)
# TODO: 如有需要，请添加测试并安装目标。
