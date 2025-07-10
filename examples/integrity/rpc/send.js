import { sleep } from "k6";
import { createSharedClients, call } from "k6/x/gasper/loadtest";
import { validateResult } from "../../utils.js";

const envs = {
  rpc: {
    CONFIG_PATH: "./examples/integrity/rpc/config.yml",
    UID: "rpc",
  },
};

export function setup() {
  for (const [_, env] of Object.entries(envs)) {
    createSharedClients(env.CONFIG_PATH, env.UID);
  }
}

export const options = {
  setupTimeout: "10m",
  discardResponseBodies: true,
  scenarios: {
    rpc: {
      executor: "constant-vus",
      duration: "10s",
      env: envs["rpc"],
      vus: 1,
    },
  },

  thresholds: {
    http_req_duration: ["p(99)<1500"], // 99% of requests must complete below 1.5s
  },
};

export default function () {
  const uid = __ENV.UID;

  const id = call(uid, "eth_chainId");
  logResult(id, "eth_chainId");
  sleep(1);

  const blockNumber = call(uid, "eth_blockNumber");
  logResult(blockNumber, "eth_blockNumber");
  sleep(1);

  const latestBlockNumber = call(uid, "eth_getBlockByNumber", [
    "latest",
    false,
  ]);
  logResult(latestBlockNumber, "eth_getBlockByNumber");
  sleep(1);
}

function logResult(obj, call) {
  for (const [_, res] of Object.entries(obj)) {
    const status = validateResult(res, "successful_" + call);
    if (status) {
      console.log(res.data);
    }
  }
}
