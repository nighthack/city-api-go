#! /bin/bash
mongoimport --host 192.168.0.101  --port 8001 --db nht_cities --collection city --type json --file ./cities.json 