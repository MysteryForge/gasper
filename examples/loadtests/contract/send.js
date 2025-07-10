import { sleep, fail } from "k6";
import {
  createSharedClients,
  reportBlockMetrics,
  txPoolStatus,
  deployContract,
  txContract,
  callContract,
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
  state_filler: {
    CONFIG_PATH: "./examples/loadtests/contract/config.yml",
    UID: "state_filler",
    ABI_PATH: "./bindings/StateFiller.abi",
    BIN_PATH: "./bindings/StateFiller.bin",
  },
};

// setup is run once before the test starts
// we need to create shared clients in the setup
export function setup() {
  const contracts = {};
  for (const [_, env] of Object.entries(envs)) {
    createSharedClients(env.CONFIG_PATH, env.UID);

    if (env.ABI_PATH && env.BIN_PATH) {
      const data = deployContract(env.UID, {
        abi_path: env.ABI_PATH,
        bin_path: env.BIN_PATH,
        gas_limit: 5995000,
        args: [{ type: "uint256", value: "250" }],
      });
      for (const [_, res] of Object.entries(data)) {
        if (res.err != undefined) {
          fail(res.err);
        }
        contracts[env.UID] = res.data.contract_address;
      }
    }
  }

  return contracts;
}

export const options = {
  setupTimeout: "10m",
  discardResponseBodies: true,
  scenarios: {
    // reporting tx pool status
    txpool_status: {
      executor: "constant-vus",
      exec: "fireTxPoolStatus",
      duration: "5s",
      env: envs["txpool_status"],
      vus: 1,
    },
    // report has to be a scenario because we cannot use vu state
    report: {
      executor: "constant-vus",
      duration: "5s",
      exec: "reportBlocks",
      vus: 1,
      env: envs["report"],
    },
    state_filler: {
      executor: "constant-vus",
      duration: "5s",
      exec: "fillState",
      vus: 1,
      env: envs["state_filler"],
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

export function fillState(contracts) {
  const uid = __ENV.UID;

  const hashes = txContract(uid, {
    contract_address: contracts[uid],
    method: "deleteRandomState",
    args: [],
  });
  for (const [_, res] of Object.entries(hashes)) {
    const status = validateResult(res, "successful_tx_delete_random_state");
    if (status) {
      console.log(res.data);
    }
  }

  const items = callContract(uid, {
    contract_address: contracts[uid],
    method: "size",
    args: [],
  });
  for (const [_, res] of Object.entries(items)) {
    const status = validateResult(res, "successful_call_contract_size");
    if (status) {
      console.log(res.data);
    }
  }

  const owner = callContract(uid, {
    contract_address: contracts[uid],
    method: "getItems",
    args: [{ type: "uint256", value: "0" }],
  });
  for (const [_, res] of Object.entries(owner)) {
    const status = validateResult(res, "successful_call_contract_get_items");
    if (status) {
      console.log(res.data);
    }
  }

  sleep(1);
}
