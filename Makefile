fmt:
	golangci-lint run && treefmt
build:
	mkdir -p bin
	$(eval XK6_CMD := $(shell command -v xk6ea 2>/dev/null || echo xk6))
	$(XK6_CMD) build --with github.com/mysteryforge/gasper=. --with github.com/grafana/xk6-output-influxdb --output bin/gasper
gosec:
	gosec -conf=./.gosec.json .
test:
	go test --count=1 ./...
init-submodule:
	git submodule update --init --recursive

# Load tests
send-contract:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/loadtests/contract/send.js
send-nonce-offset:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/loadtests/nonce_offset/send.js
send-load:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/loadtests/load/send.js
send-marketing:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/loadtests/marketing/send.js
send-single:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/loadtests/single/send.js

# Report
report:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/loadtests/report.js

# Integrity
hello:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/integrity/hello/hello.js
access_list:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/integrity/access_list/send.js
rpc:
	./bin/gasper run --out xk6-influxdb=http://localhost:8086/gasper examples/integrity/rpc/send.js