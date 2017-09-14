#!/bin/sh

DIR_LOGS="/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/myplayground/papertrailtests/all-data/*.tsv"

OUTFILE="Jun-Jul-2017.logs"

for flog in $DIR_LOGS
do
    echo "Working on $flog"
    cat $flog | awk -F$'\t' '{print $5 "|@|" $10}' >> $OUTFILE
done
