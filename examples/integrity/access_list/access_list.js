import {fail} from "k6";
import {createSharedClients, deployContract, txContract, txInfoByHash} from "k6/x/gasper/loadtest";

const env = {
        CONFIG_PATH: "./examples/integrity/access_list/access_list.yml",
        UID: "access_list",
        ABI_PATH: "./bindings/DummyStorage.abi",
        BIN_PATH: "./bindings/DummyStorage.bin",
};

// setup is run once before the test starts
// we need to create shared clients in the setup
export function setup() {
    createSharedClients(env.CONFIG_PATH, env.UID);
}

export default function () {
    console.log("access list starting...");

    // first we need to deploy the contract that we're going to call with and without an access list
    // to ensure that the gas consumed by each transaction is different and that by using an access list
    // we use less gas
    const data = deployContract(env.UID, {
        "gas_limit": 5995000,
        "abi_path": env.ABI_PATH,
        "bin_path": env.BIN_PATH,
        "args": [
            {type: "uint256", value: "10"},
        ]
    });

    let txHash, contractAddress;
    for (const [_, res] of Object.entries(data)) {
        if (res.err) {
            fail(res.err);
        }
        txHash = res.data.transaction_hash;
        contractAddress = res.data.contract_address;
    }

    console.log("deployed contract with tx hash", txHash);
    console.log("deployed contract with address", contractAddress);

    let tx1Gas, tx2Gas, tx1Hash, tx2Hash;

    // now first we're going to call the contract without an access list
    const tx1 = txContract(env.UID, {
        "contract_address": contractAddress,
        "gas_price_multiplier": 2,
        "method": "addValue",
        "args": [
            {type: "uint256", value: "100"},
        ]
    });
    for (const [_, res] of Object.entries(tx1)) {
        if (res.err) {
            fail(res.err);
        }
        console.log("tx1 hash", res.data);
        tx1Hash = res.data;
    }

    // now we're going to call the contract with an access list
    const tx2 = txContract(env.UID,
        {
            "contract_address": contractAddress,
            "method": "addValue",
            "gas_price_multiplier": 2,
            "access_list": [
                {
                    "address": contractAddress,
                    "storage_keys": [
                        "0x0", "0x1", "0x2"
                    ]
                }
            ],
            "args": [
                {type: "uint256", value: "100"},
            ]
        });
    for (const [_, res] of Object.entries(tx2)) {
        if (res.err) {
            fail(res.err);
        }
        console.log("tx2 hash", res.data);
        tx2Hash = res.data;
    }

    // now we need to get the gas used by each transaction
    const tx1Info = txInfoByHash(env.UID, tx1Hash)
    for (const [_, res] of Object.entries(tx1Info)) {
        if (res.err) {
            fail(res.err);
        }
        console.log("tx1 gas used", res.data.gas_used);
        tx1Gas = res.data.gas_used;
    }

    const tx2Info = txInfoByHash(env.UID, tx2Hash)
    for (const [_, res] of Object.entries(tx2Info)) {
        if (res.err) {
            fail(res.err);
        }
        console.log("tx2 gas used", res.data.gas_used);
        tx2Gas = res.data.gas_used;
    }

    // the example here is very contrived and the access list tx actually consumes more gas than the regular tx does
    // 2400 + 3 × 1900 − (3 × 2000) = 2100 is the exact difference expected.  The cost of the access list is greater than
    // the saving from having it in the first place.

    if (tx2Gas - tx1Gas !== 2100) {
        const diff = tx2Gas - tx1Gas;
        const message = `FAIL! tx2 gas is not 2100 different than tx1 gas, actual diff: ${diff}`;
        fail(message);
    } else {
        console.log("PASS! tx2 gas is 2100 different than tx1 gas");
    }

    console.log("access list finished");
}
