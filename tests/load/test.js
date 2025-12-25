import http from 'k6/http';
import { check, sleep } from 'k6';

export let options = {
    stages: [
        { duration: '30s', target: 20 }, // Ramp up to 20 users
        { duration: '1m', target: 20 },  // Stay at 20 users
        { duration: '30s', target: 0 },  // Ramp down
    ],
    thresholds: {
        http_req_duration: ['p(95)<500'], // 95% of requests should be faster than 500ms
    },
};

export default function () {
    const phrases = ['linux', 'binary', 'apple', 'xkcd', 'computer'];
    const phrase = phrases[Math.floor(Math.random() * phrases.length)];

    // Use 'api' hostname because k6 will run in docker network
    const res = http.get(`http://api:8080/api/search?phrase=${phrase}`);

    check(res, {
        'status is 200': (r) => r.status === 200,
        'content present': (r) => r.body.includes('comics'),
    });

    sleep(1);
}
