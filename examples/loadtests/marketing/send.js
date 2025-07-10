import { sleep } from "k6";
import {
  createSharedClients,
  reportBlockMetrics,
  txPoolStatus,
  sendTransaction,
  requestSharedWallet,
  sendERC20Transaction,
  sendERC721Transaction,
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
    CONFIG_PATH: "./examples/loadtests/marketing/config_eip155.yml",
    UID: "eip155",
  },
  erc20: {
    CONFIG_PATH: "./examples/loadtests/marketing/config_erc20.yml",
    UID: "erc20",
  },
  erc721: {
    CONFIG_PATH: "./examples/loadtests/marketing/config_erc721.yml",
    UID: "erc721",
  },
};

export function setup() {
  for (const [_, env] of Object.entries(envs)) {
    createSharedClients(env.CONFIG_PATH, env.UID);
  }

  // request a shared wallet that will be used for all vus in the eip155 scenario
  const wallets = {};
  for (const [_, env] of Object.entries(envs)) {
    if (env.UID == "txpool_status" || env.UID == "report") {
      continue;
    }
    const res = requestSharedWallet(env.UID);
    for (const [_, r] of Object.entries(res)) {
      if (r.err == undefined) {
        if (wallets[env.UID] == undefined) {
          wallets[env.UID] = [];
        }
        wallets[env.UID].push(r.data);
      }
    }
  }
  return wallets;
}

export const options = {
  setupTimeout: "10m",
  discardResponseBodies: true,
  scenarios: {
    txpool_status: {
      executor: "constant-vus",
      exec: "reportTxPoolStatus",
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

    eip155: {
      executor: "constant-vus",
      exec: "fireEIP155",
      duration: "6m",
      vus: 80,
      env: envs["eip155"],
      gracefulStop: "0s",
    },
    erc20: {
      executor: "constant-vus",
      exec: "fireERC20",
      duration: "6m",
      vus: 10,
      env: envs["erc20"],
      gracefulStop: "0s",
    },
    erc721: {
      executor: "constant-vus",
      exec: "fireERC721",
      duration: "6m",
      vus: 10,
      env: envs["erc721"],
      gracefulStop: "0s",
    },
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

export function reportTxPoolStatus(data) {
  const uid = __ENV.UID;

  const status = txPoolStatus(uid);
  for (const [_, res] of Object.entries(status)) {
    validateResult(res, "successful_tx_pool_status");
  }

  sleep(0.5);
}

export function fireEIP155(data) {
  const uid = __ENV.UID;
  const wallets = data[uid];

  const hash = sendTransaction(uid, {
    tx_count: 100000,
    wallets: wallets,
  });
  for (const [_, res] of Object.entries(hash)) {
    validateResult(res, "successful_transaction");
  }
}

export function fireERC20(data) {
  const uid = __ENV.UID;
  const wallets = data[uid];

  const hash = sendERC20Transaction(uid, {
    tx_count: 1000,
    wallets: wallets,
  });
  for (const [_, res] of Object.entries(hash)) {
    validateResult(res, "successful_erc20_transaction");
  }
}

export function fireERC721(data) {
  const uid = __ENV.UID;
  const wallets = data[uid];

  const hash = sendERC721Transaction(uid, {
    tx_count: 1000,
    wallets: wallets,
  });
  for (const [_, res] of Object.entries(hash)) {
    validateResult(res, "successful_erc721_transaction");
  }
}
