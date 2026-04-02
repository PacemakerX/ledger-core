import http from "k6/http";
import { check } from "k6";
import { uuidv4 } from "https://jslib.k6.io/k6-utils/1.4.0/index.js";

const accounts = [
  "587b276a-67ae-4bf2-8b74-9a12a2b4df1b",
  "e7ece2c9-7a26-4e39-9a12-3f94e58021a8",
  "1f41a179-f485-4008-95b4-128a82015d9f",
  "e92e7023-2d5b-4336-a59e-466dc9feed10",
  "1e5e5b18-146f-4075-8c2e-7c47f85b2472",
  "c557b21f-9d29-4cb2-9a4b-f5a9b0f3d4ce",
  "0b7613ad-638f-4ce5-af23-df9a5da0d65f",
  "4ef92e48-8749-4aff-8546-78f34048d727",
  "78614102-fa02-4a2d-b062-945c60258124",
  "4ac8ddcd-61be-4b13-b5c0-c39e2e487e2d",
  "bb5c3703-e46e-4664-84dc-a7274f9c97da",
  "97b86294-df50-476a-890a-3fd0cd9ceeb6",
  "15746868-130a-4962-afa9-87c8ccecf63d",
  "4b498f2b-388d-418d-a1b7-4668575ca9fe",
  "e28fbdb4-764b-4cdb-b36c-dfbbc2827b37",
  "132805c3-b113-40c7-addb-1963df8f688d",
  "890b8e4c-c9da-4b11-9c2f-c9877ae53660",
  "ea48d3b5-c548-40c1-89e1-9d8a057bcfe6",
  "ca6bdbf1-ce97-4815-a3f1-d22afd30a6d0",
  "6eb37917-b4ee-4abc-b9f4-8d8132a32e1b",
];

export const options = {
  stages: [
    { duration: "30s", target: 50 },
    { duration: "60s", target: 100 },
    { duration: "60s", target: 200 },
    { duration: "30s", target: 0 },
  ],
  thresholds: {
    http_req_duration: ["p(99)<600", "p(95)<500"],
    http_req_failed: ["rate<0.01"],
  },
};

export default function () {
  // Pick two different random accounts
  const fromIndex = Math.floor(Math.random() * accounts.length);
  let toIndex = Math.floor(Math.random() * accounts.length);
  while (toIndex === fromIndex) {
    toIndex = Math.floor(Math.random() * accounts.length);
  }

  const payload = JSON.stringify({
    from_account_id: accounts[fromIndex],
    to_account_id: accounts[toIndex],
    amount: 1,
    currency: "INR",
    idempotency_key: uuidv4(),
  });

  const params = {
    headers: { "Content-Type": "application/json" },
  };

  const res = http.post(
    "http://localhost:8080/api/v1/transfers",
    payload,
    params,
  );

  check(res, {
    "status is 200": (r) => r.status === 200,
  });
}
