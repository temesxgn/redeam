FROM mongo

COPY data.json /data.json

CMD mongoimport --host mongo --db redeam --collection book --type json --file /data.json --jsonArray