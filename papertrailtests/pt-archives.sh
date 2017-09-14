#!/bin/sh

#.. On Linux replace
#.. --date='{} days ago' to -v-{}d


START=$1
NDAYS=$2
HTTPAPIKEY=wAOHfYwtG7viLeec9BQ

seq 1 $NDAYS | xargs -I {} date -u -v-{}d +%Y-%m-%d | \
    xargs -I {} curl --progress-bar --no-include -o {}.tsv.gz \
    -L -H "X-Papertrail-Token: $HTTPAPIKEY" https://papertrailapp.com/api/v1/archives/{}/download
