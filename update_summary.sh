#!/bin/bash

update_string () {
    old="\[$1\].*md"
    #echo "$old"
    new="[$1]($2)"
    echo "$old --> $new"

    # in case "-i" doesn't work
    sed "s@$old@$new@g" "$3" > "$3.bak"
    mv -f "$3.bak" "$3"
}

if [ "$#" != "2" ]; then
    echo "help: update_summary.sh env (live or test) port (port number if env=test)" 
    exit -1
fi

echo "Updating SMMMARY.md with Parameters: env $1, port $2"

if [ "$1" == "live" ]; then
    echo "updating summary in live env"
    NEW="https://brofu.github.io"
else
    echo "updating summary in test env"
    NEW="http://localhost:$2/"
fi

# list to update
# to avoid to use hashmap in shell
# update here
topics=( "Golang,golang" "Elasticsearch,es" "Kafka,kafka" "Rpc,rpc" "Database Theory,database" \
    "Antlr,antlr" "Vim,vim")

for i in "${topics[@]}"
do
    #tms=$(echo "$i" | tr "," "\n")

    # reference: https://stackoverflow.com/questions/10586153/how-to-split-a-string-into-an-array-in-bash
    IFS=',' read -r -a tms <<< "$i"
    topic=${tms[0]}
    if [ "$1" == "live" ]; then
    new="$NEW/notes-${tms[1]}"
    else
    new=$NEW
    fi

    update_string "$topic" $new SUMMARY.md
done


