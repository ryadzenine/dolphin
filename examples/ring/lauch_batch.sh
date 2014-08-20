#!/bin/bash
# $1 String : path to the learning data
# $2 Int : Tau 
# $3 int : Workers
# $4 String : Configuration file
# $5 Int : Number of instances to launch
for me in $(seq 1 $5)
 do
   echo "launching process" $me
   ./ring -network=$4 -me=$me -learning-data=$1  -workers=$3 -tau=$2 -smooth=0.01 >> $(basename $1)-tau-$2-wk-$3-inst-$5.log &
 done
