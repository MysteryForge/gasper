# Contracts

- ERC20
- ERC721
- StateFiller

# Generate GO Bindings

```
make gen
```

# Instructions

When generating `*.abi` and `*.bin` files of contract make sure to follow the instructions below:

```
cat ./out/<>.sol/<>.json | jq -r '.abi'              | tr -d '\n' > ../bindings/<>.abi
cat ./out/<>.sol/<>.json | jq -r '.bytecode.object'  | tr -d '\n' > ../bindings/<>.bin
```