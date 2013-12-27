#!/bin/sh -e

rm -rf /tmp/seed-db
curl -d @users.json http://localhost:9000/user
curl http://localhost:9000/user


