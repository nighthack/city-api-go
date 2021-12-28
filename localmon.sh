#! /bin/bash
mongoimport --host localhost --port 27017 --db nht_cities --collection city --type json --file ./cities.json 