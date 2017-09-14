#!/bin/sh

DIR_LOGS="/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/myplayground/papertrailtests/mcMaster/*.tsv"

STRING_SEARCH="button pressed"

for flog in $DIR_LOGS
do
    logdate=$(echo $flog  | rev | cut -f1 -d'/' | cut -f2 -d'.' | rev)
    grep -i "$STRING_SEARCH" $flog | cut -f5 -d$'\t' | sort | uniq -c | while read LINE
    do
        echo $logdate $LINE
    done
done
