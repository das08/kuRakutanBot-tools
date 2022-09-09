import http from 'k6/http';
import { check, sleep } from 'k6';

export default function () {
    const res = http.get(`http://localhost:8081/hello`);
    check(res, {
        'status was 200': (r) => r.status === 200,
        'body was OK': (r) => r.body === "Hello World!!",
    });
    sleep(1);
}