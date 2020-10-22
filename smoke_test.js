import { check } from 'k6';
import http from 'k6/http';

export let options = {
    vus: 10,
    duration: '10s',
};

export default function () {
    let params = {
        headers: {
            'x-traefik-fault-delay-request': Math.floor(Math.random() * 1000),
            'x-traefik-fault-delay-request-percentage': Math.floor(Math.random() * 100),
            'x-traefik-fault-abort-request': 404,
            'x-traefik-fault-abort-request-percentage': 50
        },
    };

    var res = http.get(`http://${__ENV.MY_HOSTNAME}/`, params);

    check(res, { 'duration was <= 500ms': (r) => r.timings.duration <= 500 });
    check(res, { 'status was 200': (r) => r.status == 200 });
}