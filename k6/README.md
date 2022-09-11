# k6 load test
```shell
docker compose run --rm -T k6 run - < scenario/omikuji.js  --vus 100 --duration 300s
```