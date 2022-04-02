import { check } from 'k6';
import http from 'k6/http';

export const options = {
	stages: [
		{duration: '10s', vus: 20,   target: 20}, // warmup.
		{duration: '5m',  vus: 100,  target: 100},
		{duration: '5m',  vus: 1000, target: 1000},
	],
};

export default function () {
	const url = `http://${__ENV.HOSTNAME}/api/v1/profiles?first_name=${__ENV.FIRST_NAME}`;

	const params = {
		headers: {
			'Authorization': `Bearer ${__ENV.TOKEN}`,
		},
	};

	const res = http.get(url, params)
	check(res, {'is status 200': (r) => r.status === 200})
}
