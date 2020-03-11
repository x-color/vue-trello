#!/bin/bash

set -eu

echo "#############################################"
echo "############       Start        #############"
echo "#############################################"
echo ""
echo "#############################################"
echo "############        Auth        #############"
echo "#############################################"

echo "----SignUp----"
curl -s -i -X POST localhost:8080/auth/signup -H 'Content-Type:application/json; charset=UTF-8' -d '{"name":"testuser", "password":"pass"}'
echo ""

echo "----Signin----"
curl -s -i -X POST localhost:8080/auth/signin -H 'Content-Type:application/json; charset=UTF-8' -d '{"name":"testuser", "password":"pass"}' -c /tmp/cookie.file
echo ""

echo "#############################################"
echo "############        Get         #############"
echo "#############################################"

echo "----Get Resources----"
curl -s -i localhost:8080/api/resources -H 'X-XSRF-TOKEN:csrf' -H 'Content-Type:application/json; charset=UTF-8' -b /tmp/cookie.file
echo ""

echo "----Get Boards----"
curl -s -i localhost:8080/api/boards -H 'X-XSRF-TOKEN:csrf' -H 'Content-Type:application/json; charset=UTF-8' -b /tmp/cookie.file
echo ""

echo "#############################################"
echo "############       Create       #############"
echo "#############################################"

echo "----Create Board----"
curl -s -i -X POST localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"title": "test", "color":"red"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file
echo ""

BID=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

echo "----Create List----"
curl -s -i -X POST localhost:8080/api/lists \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"board_id":"'$BID'", "title": "list"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file
echo ""

LID=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

echo "----Create Item----"
curl -s -i -X POST localhost:8080/api/items \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"list_id":"'$LID'", "title": "item", "text": "hahaha", "tags":[]}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file
echo ""

echo "----Create List----"
curl -s -i -X POST localhost:8080/api/lists \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"board_id":"'$BID'", "title": "list"}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file
echo ""

LID=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

echo "----Create Item----"
curl -s -i -X POST localhost:8080/api/items \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"list_id":"'$LID'", "title": "item", "text": "hahaha", "tags":[]}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file
echo ""

echo "----Create Item----"
curl -s -i -X POST localhost:8080/api/items \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"list_id":"'$LID'", "title": "item", "text": "hahaha", "tags":[]}' \
-b /tmp/cookie.file \
| tee /tmp/tmp.file
echo ""

IID=$(cat /tmp/tmp.file | tail -1 | jq .id -r)

echo "#############################################"
echo "############       Check        #############"
echo "#############################################"

echo "----Get Board----"
curl -s -i localhost:8080/api/boards/$BID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file
echo ""

echo "#############################################"
echo "############       Update       #############"
echo "#############################################"

echo "----Update Board----"
curl -s -i -X PATCH localhost:8080/api/boards/$BID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"title": "edited title", "color":"blue", "text":"Additional text"}' \
-b /tmp/cookie.file
echo ""

echo "----Update List----"
curl -s -i -X PATCH localhost:8080/api/lists/$LID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"board_id":"'$BID'", "title": "edited list"}' \
-b /tmp/cookie.file
echo ""

echo "----Update Item----"
curl -s -i -X PATCH localhost:8080/api/items/$IID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-d '{"list_id":"'$LID'", "title": "edited item", "text": "fofofo", "tags":["1","2"]}' \
-b /tmp/cookie.file
echo ""

echo "#############################################"
echo "############       Check        #############"
echo "#############################################"

echo "----Get Boards----"
curl -s -i localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file
echo ""

ID=$(cat /tmp/tmp.file | tail -1 | jq .boards[0].id -r)

echo "----Get Board----"
curl -s -i localhost:8080/api/boards/$BID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file
echo ""

echo "#############################################"
echo "############       Delete       #############"
echo "#############################################"

echo "----Delete Item----"
curl -s -i -X DELETE localhost:8080/api/items/$IID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file
echo ""

echo "----Delete List----"
curl -s -i -X DELETE localhost:8080/api/lists/$LID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file
echo ""

echo "----Delete Board----"
curl -s -i -X DELETE localhost:8080/api/boards/$BID \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file
echo ""

echo "#############################################"
echo "############       Check        #############"
echo "#############################################"

echo "----Get Boards----"
curl -s -i localhost:8080/api/boards \
-H 'X-XSRF-TOKEN:csrf' \
-H 'Content-Type:application/json; charset=UTF-8' \
-b /tmp/cookie.file
echo ""

ID=$(cat /tmp/tmp.file | tail -1 | jq .boards[0].id -r)

echo "#############################################"
echo "############       Finish       #############"
echo "#############################################"
