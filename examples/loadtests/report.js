import { sleep, fail } from "k6";
import {
  createSharedClients,
  reportBlockMetrics,
  txPoolStatus,
} from "k6/x/gasper/loadtest";
import { validateResult } from "../utils.js";

const envs = {
  txpool_status: {
    CONFIG_PATH: "./examples/loadtests/config_txpool_status.yml",
    UID: "txpool_status",
  },
  report: {
    CONFIG_PATH: "./examples/loadtests/config_report.yml",
    UID: "report",
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
    txpool_status: {
      executor: "constant-vus",
      exec: "fireTxPoolStatus",
      duration: "12m",
      env: envs["txpool_status"],
      vus: 1,
    },
    report: {
      executor: "constant-vus",
      exec: "reportBlocks",
      duration: "12m",
      env: envs["report"],
      vus: 1,
    },
  },

  thresholds: {
    http_req_duration: ["p(99)<1500"], // 99% of requests must complete below 1.5s
  },
};

export function reportBlocks() {
  const uid = __ENV.UID;

  const result = reportBlockMetrics(uid);
  for (const [_, res] of Object.entries(result)) {
    validateResult(res, "successful_block_report");
  }

  sleep(0.075); // lower is better, 1/3 of the block time
}

export function fireTxPoolStatus(data) {
  const uid = __ENV.UID;

  const status = txPoolStatus(uid);
  for (const [_, res] of Object.entries(status)) {
    validateResult(res, "successful_tx_pool_status");
  }

  sleep(1);
}
