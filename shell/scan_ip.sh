#!/usr/bin/bash

declare -a arr
arr_index=0

prefix=192.168.2. # 前缀
start=0         # 开始
end=255           # 结束
timeout=0.1       # 单位 秒 100ms
ping_count=1      # ping 的次数

for ((i = $start; i <= $end; i++)); do

	ping -c$ping_count -W$timeout $prefix$i &>/dev/null
	flag=$?
	if [ $flag -eq 0 ]; then
		echo -e "\033[32m $prefix$i ping success!\033[0m"
		arr[$index]=$prefix$i
		index=$((index + 1))
	else
		echo -e "\033[31m $prefix$i ping failed!\033[0m"
	fi
done

echo -e "\t\t\t\t\tsuccess list:"

for i in ${arr[*]}; do
	echo -e "\033[32m ip: $i existed\033[0m"
done

echo "total: $((end-start+1)), found: ${#arr[*]}"
