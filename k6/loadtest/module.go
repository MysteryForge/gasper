package loadtest

import (
	"context"
	"fmt"
	"log/slog"
	"sync"

	"os"

	"github.com/ethereum/go-ethereum/common"
	"github.com/go-logr/logr"
	"github.com/google/uuid"
	"github.com/mysteryforge/gasper/k6/eth"
	"go.k6.io/k6/js/modules"
)

var (
	sharedClients    map[string]*Clients
	initOnceClients  map[string]*sync.Once
	sharedDBs        map[string]*eth.PebbleDb
	initOnceDBs      map[string]*sync.Once
	sharedLog        logr.Logger
	sharedWallets    map[string]struct{}
	sharedWalletsMux *sync.Mutex
	TestUID          string
)

func init() {
	id, _ := uuid.NewUUID()
	TestUID = id.String()[:8]
	sharedLog = logr.FromSlogHandler(slog.NewJSONHandler(os.Stdout, &slog.HandlerOptions{
		AddSource: false,
	}))
	sharedLog = sharedLog.WithValues("test_uid", TestUID)
	sharedClients = make(map[string]*Clients)
	initOnceClients = make(map[string]*sync.Once)
	sharedDBs = make(map[string]*eth.PebbleDb)
	initOnceDBs = make(map[string]*sync.Once)
	sharedWallets = make(map[string]struct{})
	sharedWalletsMux = &sync.Mutex{}
}

type RootModule struct{}

type ModuleInstance struct {
	vu      modules.VU
	metrics *EthMetrics
}

var (
	_ modules.Instance = &ModuleInstance{}
	_ modules.Module   = &RootModule{}
)

func New() *RootModule {
	return &RootModule{}
}

func (*RootModule) NewModuleInstance(vu modules.VU) modules.Instance {
	return &ModuleInstance{
		vu:      vu,
		metrics: RegisterMetrics(vu),
	}
}

func (mi *ModuleInstance) Exports() modules.Exports {
	return modules.Exports{
		Named: map[string]interface{}{
			"createSharedClients": func(pth, uid string) interface{} {
				if pth == "" {
					panic("CONFIG_PATH is not set")
				}
				if uid == "" {
					panic("uid is not set")
				}
				createSharedClients(pth, uid) // create shared clients
				return nil
			},
			"requestSharedWallet": func(uid string) interface{} {
				panicIfNotInitialized(uid)
				sharedWalletsMux.Lock()
				defer sharedWalletsMux.Unlock()

				res := sharedClients[uid].RequestSharedWallet()
				for _, r := range res {
					if r.Err != nil {
						panic(r.Err)
					}

					walletUid := fmt.Sprintf("%s_%s", r.Data.(*SharedWallet).Address, r.Data.(*SharedWallet).UID)
					if _, exists := sharedWallets[walletUid]; exists {
						panic("wallet already requested and shared, cannot have the same wallet shared across multiple scenarios")
					}
					sharedWallets[walletUid] = struct{}{}
				}

				return res
			},
			"releaseSharedWallet": func(uid string, address string) interface{} {
				panicIfNotInitialized(uid)
				if address == "" {
					panic("address is not set")
				}
				sharedClients[uid].ReleaseSharedWallet(common.HexToAddress(address))
				return nil
			},

			"chainID": func(uid string) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].ChainID(mi.vu, mi.metrics)
			},
			"txPoolStatus": func(uid string) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].TxPoolStatus(mi.vu, mi.metrics)
			},
			"reportBlockMetrics": func(uid string) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].ReportBlockMetrics(mi.vu, mi.metrics)
			},

			"sendTransaction": func(uid string, params map[string]interface{}) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].SendTransaction(mi.vu, mi.metrics, params)
			},
			"sendERC20Transaction": func(uid string, params map[string]interface{}) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].SendERC20Transaction(mi.vu, mi.metrics, params)
			},
			"sendERC721Transaction": func(uid string, params map[string]interface{}) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].SendERC721Transaction(mi.vu, mi.metrics, params)
			},

			"deployContract": func(uid string, params map[string]interface{}) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].DeployContract(mi.vu, params)
			},
			"txContract": func(uid string, params map[string]interface{}) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].TxContract(mi.vu, mi.metrics, params)
			},
			"callContract": func(uid string, params map[string]interface{}) interface{} {
				panicIfNotInitialized(uid)

				return sharedClients[uid].CallContract(mi.vu, mi.metrics, params)
			},
			"txInfoByHash": func(uid string, hash string) interface{} {
				panicIfNotInitialized(uid)
				return sharedClients[uid].TxInfoByHash(mi.vu, mi.metrics, hash)
			},
		},
	}
}

func createSharedClients(pth, uid string) {
	var err error
	if initOnceClients[uid] == nil {
		initOnceClients[uid] = &sync.Once{}
	}

	initOnceClients[uid].Do(func() {
		var cfgs []*clientConfig
		cfgs, err = ReadConfigYML(pth)
		if err != nil {
			panic(fmt.Errorf("failed to read config file: %w", err))
		}

		if len(cfgs) == 0 {
			panic(fmt.Errorf("no configurations found in the file"))
		}

		for _, cfg := range cfgs {
			if cfg.DBPath == "" {
				continue
			}
			if initOnceDBs[cfg.DBPath] == nil {
				initOnceDBs[cfg.DBPath] = &sync.Once{}
			}
			initOnceDBs[cfg.DBPath].Do(func() {
				sharedDBs[cfg.DBPath], err = eth.NewPebbleDb(cfg.DBPath)
				if err != nil {
					panic(err)
				}
			})
		}

		sharedClients[uid], err = NewClients(context.Background(), cfgs, uid, sharedDBs, sharedLog)
		if err != nil {
			panic(err)
		}
	})
}

func panicIfNotInitialized(uid string) {
	if uid == "" {
		panic("uid is not set")
	}
	if sharedClients[uid] == nil {
		panic("sharedClients is not initialized")
	}
}
