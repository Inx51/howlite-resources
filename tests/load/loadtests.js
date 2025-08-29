import http from 'k6/http';
import { sleep } from 'k6';
import { randomString } from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
  duration: '2m',
  rps: 1000,
  vus: 1000,
};

// The default exported function is gonna be picked up by k6 as the entry point for the test script. It will be executed repeatedly in "iterations" for the whole duration of the test.
export default function () {
    // Make a GET request to the target URL
    const body = JSON.stringify({
        key1: 'value1',
        key2: 'value2',
    });

    const headers = {
        headers: {
            'Content-Type': 'application/json',
        },
    };
    http.post(`http://howlite-resources:8080/${randomString(30)}`, body, headers);
}