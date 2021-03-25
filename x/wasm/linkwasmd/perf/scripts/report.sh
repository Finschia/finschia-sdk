#!/usr/bin/env bash
log_path=$1
output_file=$2

echo "performance report"  >> $output_file
echo ""  >> $output_file

elapsed_total=0
txs_total=0

while read line
do
  case ${line:0:1} in
    t)
      txs=${line:2}
      ;;
    s)
      start=$(date -d "${line:2}" +%s%N)
      ;;
    e)
      end=$(date -d "${line:2}" +%s%N)
      elapsed=$((end-start))

      if [ $((txs-10)) -ge 0 ];then
        echo "$txs txs in $((elapsed/1000000)) msec." >> $output_file
        txs_total=$((txs_total+txs))
        elapsed_total=$((elapsed_total+elapsed))
      else
        echo "$txs txs in $((elapsed/1000000)) msec. drop!" >> $output_file
      fi
      ;;
  esac
done < $log_path

echo "" >> $output_file
tps=$((1000000000 * txs_total/ elapsed_total))
echo "tps:$tps" >> $output_file
echo "TPS $tps"
