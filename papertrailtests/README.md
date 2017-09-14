# Goal: Download daily data from papertrail and convert it to the format to push to Keen

## STEPS:
1. `papertrailtests/main.go` to download files from S3 and concatenate them to daily data. The output will be in `papertrail/logs/<some number>/...`
2. Copied all the daily TSV data into the folder `all-data/JUN-JUL-ALL-APT`
3. Use `convertTSVforGo.sh` to concatenate all the TSV files into a single file (`.logs`) where each row has `<devicekey>|@|<json log message>` with `|@|` as separator.
4. Use `all-data/go_keenlogs/main.go` to convert the massive `.logs` file into the keen records and push them to Keen project
