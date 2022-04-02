import { check } from 'k6';
import http from 'k6/http';

export const options = {
	discardResponseBodies: true,
	scenarios: {
		first_page: {
			executor: 'constant-vus',
			exec: 'firstPage',
			vus: 1000,
			duration: '30s',
			gracefulStop: '10s',
		},
		second_page: {
			executor: 'constant-vus',
			exec: 'secondPage',
			vus: 1000,
			duration: '30s',
			gracefulStop: '10s',
			startTime: '40s',
		},
	},
};

export function firstPage () {
	const url = `http://${__ENV.HOSTNAME}/api/v1/profiles?first_name=${__ENV.FIRST_NAME}`;
	getPage(url);
}

export function secondPage () {
	const url = `http://${__ENV.HOSTNAME}/api/v1/profiles?first_name=${__ENV.FIRST_NAME}?page_token=${__ENV.NEXT_PAGE}`;
	getPage(url);
}

function getPage(url) {
	const params = {
		headers: {
			'Authorization': `Bearer ${__ENV.TOKEN}`,
		},
	};

	const res = http.get(url, params)
	check(res, {'is status 200': (r) => r.status === 200})
}
