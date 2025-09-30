import http from 'k6/http';
import { check, sleep } from 'k6';
import { Counter } from 'k6/metrics';
import { scenario } from 'k6/execution';

export let success200 = new Counter('success200');
export let status429 = new Counter('status429');

export const options = {
    scenarios: {
        phase1: {
          executor: 'constant-arrival-rate',
          rate: 3,                 // inteiro, ~3 req/s
          timeUnit: '1s',
          duration: '34s',         // ~102 reqs em 34s ≈ 100
          preAllocatedVUs: 20,
          maxVUs: 50,
        },
        phase2: {
          executor: 'constant-arrival-rate',
          rate: 5,                 // 5 req/s
          timeUnit: '1s',
          duration: '10s',         // 50 reqs
          preAllocatedVUs: 20,
          maxVUs: 50,
          startTime: '34s',     
        },
      },
  thresholds: {
    'success200': ['count>=100'],
    'status429': ['count>=50'],
  },
};

const BASE_URL = __ENV.TARGET_URL || 'http://localhost:9999/example';
const CLIENT_ID = __ENV.CLIENT_ID || '6e421fce-2887-4b38-bff4-5887749cfb62';

export default function () {
  const currentScenario = scenario.name;

  const params = {
    headers: { 'X-Api-Id': CLIENT_ID },
    tags: { endpoint: 'example' },
  };

  let res = http.get(BASE_URL, params);

  if (currentScenario === 'phase1') {
    const ok = check(res, { 'phase1: status is 200': (r) => r.status === 200 });
    console.log(`Phase1 -> status: ${res.status}`);
    if (ok) success200.add(1);
  } else if (currentScenario === 'phase2') {
    const ok = check(res, { 'phase2: status is 429': (r) => r.status === 429 });
    console.log(`Phase2 -> status: ${res.status}`);
    if (ok) status429.add(1);
  }

  sleep(0.1);
}
