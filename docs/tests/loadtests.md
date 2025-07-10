# Load Testing

## Configuration Options

### Load testing

- `http`: (required) Ethereum node RPC endpoint URL (must be http or https)
- `private_keys`: List of private keys of the accounts used to fund new wallets
- `num_wallets`: Number of new wallets to create and fund
- `fund_amount`: Amount of ETH to fund new wallets with (in wei)
- `wallets`: List of pre-funded wallet private keys to use
- `target_addresses`: Target addresses for transactions (required when `num_wallets > 0` or `wallets` are provided)
- `num_target_addresses`: Number of new target addresses to use (required when `num_wallets > 0` or `wallets` are provided)
- `erc20`: Enable ERC20 token testing (boolean)
- `erc20_address`: Address of the ERC20 token contract (required when `erc20 = true` and `private_keys` not set)
- `erc20_mint_amount`: Amount to transfer in ERC20 token contract per wallet
- `erc721`: Enable ERC721 token testing (boolean)
- `erc721_address`: Address of the ERC721 token contract (required when `erc721 = true` and `private_keys` not set)
- `erc721_mint`: Whether to mint ERC721 token contract on startup
- `db_path`: Path to the database where we store transaction hashes to track latency
- `rate_limit`: Rate limit for transaction sending (in transactions per second)
- `adaptive_rate_limit`: Whether to use adaptive rate limiting (boolean)
- `min_gas_price`: Minimum gas price to use for transactions (in wei)

## Available Functions

### Load testing

The following functions are exposed to k6 test scripts.

#### Setup
- `createSharedClients(configPath, uid)`: Initialize shared clients with configuration file

#### Chain Information
- `chainID(uid)`: Get the chain ID
- `txPoolStatus(uid)`: Get transaction pool status
- `reportBlockMetrics(uid)`: Report current block metrics

#### Wallet Management
- `requestSharedWallet(uid)`: Request a shared wallet
- `releaseSharedWallet(uid, address)`: Release a shared wallet

#### Transaction Params:
- `tx_count`: Number of transactions to send
- `confirmation_delay`: Delay in seconds to wait for confirmation
- `no_send`: Whether to simulate sending the transaction (don't actually send)
- `nonce_offset`: Nonce offset for the transaction
- `gas_price_multiplier`: Multiplier for the gas price
- `wallets`: List of wallet addresses to use for the transaction

#### Transaction Operations
- `sendTransaction(uid, params)`: Send a basic transaction

#### Token Operations
- `sendERC20Transaction(uid, params)`: Send an ERC20 token transfer
- `sendERC721Transaction(uid, params)`: Send an ERC721 token transfer

#### Deploy Contract Params:
- `abi_path`: Path to the contract ABI file
- `bin_path`: Path to the contract binary file
- `gas_limit`: Gas limit for the transaction
- `args`: List of arguments to pass to the contract constructor

#### Deployment Operations
- `deployContract(uid, params)`: Deploy a new contract

#### Contract Params:
- `contract_address`: Address of the contract
- `method`: Method to call
- `args`: List of arguments to pass to the method

#### Contract Operations
- `txContract(uid, params)`: Send a contract transaction
- `callContract(uid, params)`: Call a contract method (read-only)
