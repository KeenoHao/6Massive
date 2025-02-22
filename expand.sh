#!/bin/bash

# 检查参数数量是否正确
if [ "$#" -ne 4 ]; then
    echo "Usage: $0 <input_file> <execution_times> <sample_size> <output_prefix>"
    exit 1
fi

# 获取参数
input_file=$1
execution_times=$2
sample_size=$3
output_prefix=$4

# 检查输入文件是否存在
if [ ! -f "$input_file" ]; then
    echo "Error: Input file '$input_file' not found!"
    exit 1
fi

# 检查执行次数和采样数量是否为正整数
if ! [[ "$execution_times" =~ ^[0-9]+$ && "$sample_size" =~ ^[0-9]+$ ]]; then
    echo "Error: Execution times and sample size must be positive integers!"
    exit 1
fi

# 对输入文件排序并去重
sort "$input_file" > "${input_file}.sorted"
uniq "${input_file}.sorted" > "${input_file}"
rm -rf "${input_file}.sorted"

# 执行随机采样
for ((i=1; i<=execution_times; i++)); do
    output_file="${output_prefix}${i}"
    shuf -n "$sample_size" "${input_file}" > "$output_file"
done

echo "Sampling completed. Output files are named as ${output_prefix}1, ${output_prefix}2, ..., ${output_prefix}${execution_times}."