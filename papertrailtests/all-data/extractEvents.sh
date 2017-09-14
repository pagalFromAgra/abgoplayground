#!/bin/sh

DIR_LOGS="/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/myplayground/papertrailtests/ccbcc/*.tsv"

COMPANY="syscoboston"
CSV_PATH="/Users/adityabansal/kineticdevs/go/src/github.com/wearkinetic/digger/exportedCSVs/export-$COMPANY.csv"


pullEvent(){
    STRING_SEARCH=$1

    for flog in $DIR_LOGS
    do
        logdate=$(echo $flog  | rev | cut -f1 -d'/' | cut -f2 -d'.' | rev)
        grep -i "$STRING_SEARCH" $flog | cut -f5 -d$'\t' | sort | uniq -c | sed 's/-1//g' | while read LINE
        do
            key=$(echo $LINE | cut -f2 -d' ')

            if [ $(grep -c $key $CSV_PATH) -ne 0 ]; then
                echo $COMPANY $logdate $LINE
            fi
        done
    done
}

pullEvent "BUTTON_PRESS_RISKY_LIFTS_GOAL"
pullEvent "BUTTON_PRESS_RANK"
pullEvent "BUTTON_PRESS_BASELINE"
pullEvent "BUTTON_PRESS_GOAL"
pullEvent "BUTTON_PRESS_TIME"
pullEvent "BUTTON_PRESS_STAY_SAFE"
