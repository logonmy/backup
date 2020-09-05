#!/usr/bin/bash

externalIP=`curl -s cip.cc | grep IP | cut -d ":" -f 2 | sed "s/^[ ]*//g"`
#internalIP=`ip addr | grep "inet\b" | grep -v "127.0.0.1" | awk '{ print $2 }' | awk -F "/" '{print $1}'`

echo "外网IP地址：$externalIP"
#echo "内网IP地址：$internalIP"
