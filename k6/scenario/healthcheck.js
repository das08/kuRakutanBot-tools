import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
    const res = http.get(`http://rakutan-app:8081/hello`);
    check(res, {
        'status was 200': (r) => r.status === 200,
        'body was OK': (r) => r.body === "Hello World!!",
    });
    sleep(1);
}

//docker compose run --rm -T k6 run - < scenario/healthcheck.js --vus 2000 --duration 5s
