package mock

import (
	"bufio"
	"context"
	"fmt"
	"strings"
	"sync"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/testcontainers/testcontainers-go"
	"github.com/testcontainers/testcontainers-go/wait"
)

const (
	TestAccountsCount = 3
	TestChainID       = 1337
)

func StartTestAnvilContainer(t *testing.T, ctx context.Context) (string, []string) {
	anvil, err := testcontainers.GenericContainer(ctx, testcontainers.GenericContainerRequest{
		ContainerRequest: testcontainers.ContainerRequest{
			Image:        "ghcr.io/foundry-rs/foundry",
			ExposedPorts: []string{"8545/tcp"},
			Cmd: []string{
				fmt.Sprintf("anvil --host=0.0.0.0 --port=8545 --chain-id=%d --accounts=%d --balance=100000 --block-time=1", TestChainID, TestAccountsCount),
			},
			WaitingFor: wait.ForLog("Listening on"),
		},
		Started: true,
	})
	require.NoError(t, err)
	t.Cleanup(func() {
		anvil.Terminate(ctx) // nolint:errcheck
	})

	logReader, err := anvil.Logs(ctx)
	require.NoError(t, err)
	t.Cleanup(func() {
		logReader.Close() // nolint:errcheck
	})

	var privateKeys []string
	wg := sync.WaitGroup{}
	wg.Add(TestAccountsCount)

	go func() {
		scanner := bufio.NewScanner(logReader)
		inPrivateKeysSection := false

		for scanner.Scan() {
			line := scanner.Text()
			// fmt.Println("[anvil]:", line)

			if strings.Contains(line, "Private Keys") {
				inPrivateKeysSection = true
				continue
			}

			if inPrivateKeysSection {
				line = strings.TrimSpace(line)
				if strings.HasPrefix(line, "(") {
					// Example line: (0) 0xac0974be...
					parts := strings.Fields(line)
					if len(parts) >= 2 {
						key := strings.TrimPrefix(parts[1], "0x")
						privateKeys = append(privateKeys, key)
						wg.Done()
					}
				}
			}
		}
	}()

	wg.Wait()

	host, err := anvil.Host(ctx)
	require.NoError(t, err)
	port, err := anvil.MappedPort(ctx, "8545")
	require.NoError(t, err)
	return fmt.Sprintf("http://%s:%s", host, port.Port()), privateKeys
}
