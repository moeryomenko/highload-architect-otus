import { check } from 'k6';
import http from 'k6/http';

import {
  randomIntBetween,
  randomString,
  randomItem,
  uuidv4,
  findBetween,
} from 'https://jslib.k6.io/k6-utils/1.2.0/index.js';

export const options = {
	discardResponseBodies: true,
	scenarios: {
		insert_users: {
			executor: 'constant-vus',
			exec: 'insertUsers',
			vus: 200,
			duration: '30s',
			gracefulStop: '0s',
		},
		fetch_profiles: {
			executor: 'constant-vus',
			exec: 'fetchProfiles',
			vus: 800,
			duration: '30s',
			gracefulStop: '30s',
		},
	},
};

export function insertUsers () {
	const resLogin = http.post(`http://${__ENV.HOSTNAME}/api/v1/signup`, JSON.stringify({
		"nickname": uuidv4(),
		"password": "test",
	}),{
		headers: {
			'Content-Type': 'application/json',
		},
	});
	check(resLogin, {'is status 200': (r) => r.status === 200})
}


export function fetchProfiles () {
	const url = `http://${__ENV.HOSTNAME}/api/v1/profiles?first_name=${__ENV.FIRST_NAME}&last_name=${__ENV.LAST_NAME}`;
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
