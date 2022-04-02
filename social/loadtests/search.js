import { check } from 'k6';
import http from 'k6/http';

export const options = {
	discardResponseBodies: true,
	scenarios: {
		constants: {
			executor: 'constant-vus',
			vus: 1000,
			duration: '30s',
		},
	},
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
