#!/bin/sh
# 获取索引列表，并排序。
# list=$(curl http://175.178.108.84:6123/_cat/indices/dev-jaeger-span-2022-11*?h=i\&s=i:desc)
list="dev-jaeger-span-2022-11-10
dev-jaeger-span-2022-11-09
dev-jaeger-span-2022-11-08
dev-jaeger-span-2022-11-07
dev-jaeger-span-2022-11-06
dev-jaeger-span-2022-11-05
dev-jaeger-span-2022-11-04
dev-jaeger-span-2022-11-03
dev-jaeger-span-2022-11-02
dev-jaeger-span-2022-11-01"
index=5 # 保留最近的指定数量
echo "-----------------------"
for item in $list; do
  if ((index>0)); then
    ((index--))
    echo "index: $index"
    continue
  fi
  echo "echo: $item" # 执行删除语句
done
echo "任务执行完毕..."