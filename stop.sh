#!/bin/bash

Proc=`ps -ef |grep -w "./metricsvr" |grep -v grep|wc -l`
if [ $Proc -le 0 ];then
    echo "Node havn't running .. "
else
    ps -ef | grep "./metricsvr" | grep -v grep | awk '{print $2}' | xargs kill -9
fi
