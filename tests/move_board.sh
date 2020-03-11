#!/bin/bash

set -eu


curl -s -X POST localhost:8080/auth/signup -H 'Content-Type:application/json; charset=UTF-8' -d '{"name":"testuser", "password":"pass"}'

curl -s -X POST localhost:8080/auth/signin -H 'Content-Type:application/json; charset=UTF-8' -d '{"name":"testuser", "password":"pass"}' -c /tmp/cookie.file

curl -s localhost:8080/api/resources -H 'X-XSRF-TOKEN:csrf' -H 'Content-Type:application/json; charset=UTF-8' -b /tmp/cookie.file

curl -s localhost:8080/api/boards -H 'X-XSRF-TOKEN:csrf' -H 'Content-Type:application/json; charset=UTF-8' -b /tmp/cookie.file

# Create

curl -s -X POST localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"title": "first", "color":"red"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file

BID1=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

curl -s -X POST localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"title": "first", "color":"red"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file

BID2=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

curl -s -X POST localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"title": "first", "color":"red"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file

BID3=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

# Create lists
curl -s -X POST localhost:8080/api/lists \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"board_id":"'$BID1'", "title": "first_list"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file

LID1=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

curl -s -X POST localhost:8080/api/lists \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"board_id":"'$BID2'", "title": "second_list"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file

LID2=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

curl -s -X POST localhost:8080/api/lists \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"board_id":"'$BID3'", "title": "first_list"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file

LID3=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

# Create items
curl -s -X POST localhost:8080/api/items \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"list_id":"'$LID1'", "title": "first_item", "text": "hahaha", "tags":[]}' \
-b /tmp/cookie.file \

curl -s -X POST localhost:8080/api/items \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"list_id":"'$LID2'", "title": "fourth_item", "text": "hahaha", "tags":[]}' \
-b /tmp/cookie.file \

curl -s -X POST localhost:8080/api/items \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"list_id":"'$LID3'", "title": "first_item", "text": "hahaha", "tags":[]}' \
-b /tmp/cookie.file \


# Show

echo 
echo "## Show ##"
echo 

curl -s localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file \
| jq .

# Move

curl -s -X PATCH localhost:8080/api/boards/$BID2/move \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"before": "'$BID3'"}' \
-b /tmp/cookie.file \

echo 
echo "## Moved ##"
echo 

curl -s localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file \
| jq .

# Delete

curl -s -X DELETE localhost:8080/api/boards/$BID2 \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"before": "'$BID3'"}' \
-b /tmp/cookie.file \

echo 
echo "## Deleted ##"
echo 

curl -s localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file \
| jq .
