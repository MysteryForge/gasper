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
  offset_nonce_tx: {
    CONFIG_PATH: "./examples/loadtests/tx_nonce_offset/config.yml",
    UID: "offset_nonce_tx",
  },
};

// setup is run once before the test starts
// we need to create shared clients in the setup
export function setup() {
  for (const [_, env] of Object.entries(envs)) {
    createSharedClients(env.CONFIG_PATH, env.UID);
  }
}

export const options = {
  timeout: "10m",
  discardResponseBodies: true,
  scenarios: {
    txpool_status: {
      executor: "constant-vus",
      exec: "fireTxPoolStatus",
      duration: "15s",
      env: envs["txpool_status"],
      vus: 1,
    },
    report: {
      executor: "constant-vus",
      duration: "15s",
      exec: "reportBlocks",
      vus: 1,
      env: envs["report"],
    },

    offset_nonce: {
      executor: "constant-vus",
      exec: "fireOffsetNonce",
      duration: "5s",
      vus: 1,
      env: envs["offset_nonce_tx"],
    },
    fire_and_verify: {
      executor: "constant-vus",
      exec: "fireAndVerify",
      duration: "5s",
      startTime: "8s",
      vus: 1,
      env: envs["offset_nonce_tx"],
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

export function fireAndVerify(data) {
  const uid = __ENV.UID;
  const hashs = sendTransaction(uid, {
    confirmation_delay: 12,
  });
  for (const [_, res] of Object.entries(hashs)) {
    validateResult(res, "successful_confirmed_transaction");
  }

  sleep(0.1);
}

export function fireOffsetNonce(data) {
  const uid = __ENV.UID;
  const hashes = sendTransaction(uid, {
    nonce_offset: 10,
  });
  for (const [_, res] of Object.entries(hashes)) {
    validateResult(res, "successful_offset_transaction");
  }

  sleep(0.1);
}
