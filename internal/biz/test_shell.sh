#!/bin/sh
# 获取索引列表，并排序。
list=$(curl --location --request GET 'http://iot.admin.jiaoranyouxuan.com:9201/_cat/indices/jaeger-span*?h=i&s=i:desc' --header 'Authorization: Basic ZWxhc3RpYzpsamgyMDIyLg==')
index=10 # 保留最近的指定数量
for item in $list; do
  if ((index>0)); then
    ((index--))
    echo "跳过索引: $index $item \r\n"
    continue
  fi
  echo "删除索引: $item " # 执行删除语句
  res=$(curl --location --request DELETE "http://iot.admin.jiaoranyouxuan.com:9201/$item" --header='Authorization: Basic ZWxhc3RpYzpsamgyMDIyLg==')
  echo "$res \r\n"
done
echo "任务执行完毕..."

