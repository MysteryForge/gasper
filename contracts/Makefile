gen:
	forge build --force

	cat ./out/ERC20.sol/ERC20.json | jq -r '.abi'             | tr -d '\n' > ../bindings/ERC20.abi
	cat ./out/ERC20.sol/ERC20.json | jq -r '.bytecode.object' | tr -d '\n' > ../bindings/ERC20.bin
	abigen --abi=../bindings/ERC20.abi --bin=../bindings/ERC20.bin --pkg=bindings --type=ERC20 --out=../bindings/ERC20.go

	cat ./out/ERC721.sol/ERC721.json | jq -r '.abi'             | tr -d '\n' > ../bindings/ERC721.abi
	cat ./out/ERC721.sol/ERC721.json | jq -r '.bytecode.object' | tr -d '\n' > ../bindings/ERC721.bin
	abigen --abi=../bindings/ERC721.abi --bin=../bindings/ERC721.bin --pkg=bindings --type=ERC721 --out=../bindings/ERC721.go

	cat ./out/StateFiller.sol/StateFiller.json | jq -r '.abi'              | tr -d '\n' > ../bindings/StateFiller.abi
	cat ./out/StateFiller.sol/StateFiller.json | jq -r '.bytecode.object'  | tr -d '\n' > ../bindings/StateFiller.bin
	abigen --abi=../bindings/StateFiller.abi --bin=../bindings/StateFiller.bin --pkg=bindings --type=StateFiller --out=../bindings/StateFiller.go

	cat ./out/BatchFunder.sol/BatchFunder.json | jq -r '.abi'              | tr -d '\n' > ../bindings/BatchFunder.abi
	cat ./out/BatchFunder.sol/BatchFunder.json | jq -r '.bytecode.object'  | tr -d '\n' > ../bindings/BatchFunder.bin
	abigen --abi=../bindings/BatchFunder.abi --bin=../bindings/BatchFunder.bin --pkg=bindings --type=BatchFunder --out=../bindings/BatchFunder.go

	cat ./out/DummyStorage.sol/DummyStorage.json | jq -r '.abi'              | tr -d '\n' > ../bindings/DummyStorage.abi
	cat ./out/DummyStorage.sol/DummyStorage.json | jq -r '.bytecode.object'  | tr -d '\n' > ../bindings/DummyStorage.bin
	abigen --abi=../bindings/DummyStorage.abi --bin=../bindings/DummyStorage.bin --pkg=bindings --type=DummyStorage --out=../bindings/DummyStorage.go

