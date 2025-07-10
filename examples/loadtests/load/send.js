import { sleep } from "k6";
import {
  createSharedClients,
  reportBlockMetrics,
  txPoolStatus,
  sendTransaction,
  sendERC20Transaction,
  sendERC721Transaction,
} from "k6/x/gasper/loadtest";
import { validateResult } from "../../utils.js";

// Environment configurations for different test scenarios
// Each environment defines:
// - CONFIG_PATH: Path to the YAML configuration file for this test
// - UID: Unique identifier for this test environment
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
    CONFIG_PATH: "./examples/loadtests/load/config_eip155.yml",
    UID: "eip155",
  },
  erc20: {
    CONFIG_PATH: "./examples/loadtests/load/config_erc20.yml",
    UID: "erc20",
  },
  erc721: {
    CONFIG_PATH: "./examples/loadtests/load/config_erc721.yml",
    UID: "erc721",
  },
};

// Setup function runs once before the test starts
// Creates shared clients for all test environments
// This allows multiple scenarios to use the same client connections
export function setup() {
  for (const [_, env] of Object.entries(envs)) {
    createSharedClients(env.CONFIG_PATH, env.UID);
  }
}

// Test options configuration
// Defines the test duration, response body handling, and scenarios
export const options = {
  setupTimeout: "10m",
  discardResponseBodies: true,
  scenarios: {
    // reporting tx pool status
    txpool_status: {
      executor: "constant-vus",
      exec: "fireTxPoolStatus",
      duration: "12m",
      env: envs["txpool_status"],
      vus: 1,
    },
    // reporting block metrics
    report: {
      executor: "constant-vus",
      exec: "reportBlocks",
      duration: "12m",
      env: envs["report"],
      vus: 1,
    },

    eip155: {
      executor: "constant-vus",
      exec: "fireEIP155",
      duration: "6m",
      vus: 100,
      env: envs["eip155"],
      gracefulStop: "0s",
    },
    erc20: {
      executor: "constant-vus",
      exec: "fireERC20",
      duration: "6m",
      vus: 100,
      env: envs["erc20"],
      gracefulStop: "0s",
    },
    erc721: {
      executor: "constant-vus",
      exec: "fireERC721",
      duration: "6m",
      vus: 100,
      env: envs["erc721"],
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

  const hash = sendTransaction(uid, {
    tx_count: 1000,
  });
  for (const [_, res] of Object.entries(hash)) {
    validateResult(res, "successful_transaction");
  }

  sleep(0.01);
}

export function fireERC20(data) {
  const uid = __ENV.UID;

  const hash = sendERC20Transaction(uid, {
    tx_count: 1000,
  });
  for (const [_, res] of Object.entries(hash)) {
    validateResult(res, "successful_erc20_transaction");
  }

  sleep(0.01);
}

export function fireERC721(data) {
  const uid = __ENV.UID;

  const hash = sendERC721Transaction(uid, {
    tx_count: 1000,
  });
  for (const [_, res] of Object.entries(hash)) {
    validateResult(res, "successful_erc721_transaction");
  }

  sleep(0.01);
}
