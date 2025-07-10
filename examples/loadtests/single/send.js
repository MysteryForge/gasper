import { sleep } from "k6";
import {
  createSharedClients,
  reportBlockMetrics,
  txPoolStatus,
  sendTransaction,
} from "k6/x/gasper/loadtest";
import { validateResult } from "../../utils.js";

const envs = {
  txpool_status: {
    CONFIG_PATH: "./examples/loadtests/config_txpool_status.yml",
    UID: "txpool_status",
  },
  report: {
    CONFIG_PATH: "./examples/loadtests/config_report.yml",
    UID: "report",
  },
  eip155: {
    CONFIG_PATH: "./examples/loadtests/single/config_eip155.yml",
    UID: "eip155",
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
    // reporting tx pool status
    txpool_status: {
      executor: "constant-vus",
      exec: "fireTxPoolStatus",
      duration: "20s",
      env: envs["txpool_status"],
      vus: 1,
    },
    // reporting block metrics
    report: {
      executor: "constant-vus",
      exec: "reportBlocks",
      duration: "20s",
      env: envs["report"],
      vus: 1,
    },

    eip155: {
      executor: "constant-vus",
      exec: "fireEIP155",
      duration: "10s",
      vus: 2,
      env: envs["eip155"],
      gracefulStop: "0s",
    },
  },

  thresholds: {
    http_req_duration: ["p(99)<1500"], // 99% of requests must complete below 1.5s
  },
};

export function reportBlocks() {
  // this case we are using the same HTTP endpoint in both scenarios so we need to report only once
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

export function fireEIP155(data) {
  const uid = __ENV.UID;

  const hash = sendTransaction(uid);
  for (const [_, res] of Object.entries(hash)) {
    const succ = validateResult(res, "successful_transaction");
    if (succ) {
      console.log(res.data);
    }
  }

  sleep(1);
}
