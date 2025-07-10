package loadtest

import (
	"strconv"
	"time"

	"github.com/mysteryforge/gasper/k6/eth"
	"go.k6.io/k6/js/modules"
	"go.k6.io/k6/metrics"
)

type EthMetrics struct {
	RequestDuration *metrics.Metric
	TimeToMine      *metrics.Metric
	GasUsed         *metrics.Metric
	BlockGasUsed    *metrics.Metric
	Txs             *metrics.Metric
	BlockTxs        *metrics.Metric
	TPS             *metrics.Metric // number of confirmed transactions that were successfully mined into blocks on the chain per second.

	Mgas              *metrics.Metric
	BlockMgas         *metrics.Metric
	EOA               *metrics.Metric // number of initiated ETH transactions, submitted via eth_sendRawTransaction.
	BlockTime         *metrics.Metric
	BlockPerSec       *metrics.Metric
	PoolStatusPending *metrics.Metric
	PoolStatusQueued  *metrics.Metric
}

func RegisterMetrics(vu modules.VU) *EthMetrics {
	r := vu.InitEnv().Registry
	return &EthMetrics{
		RequestDuration:   r.MustNewMetric("gasper_req_duration", metrics.Trend, metrics.Time),
		TimeToMine:        r.MustNewMetric("gasper_time_to_mine", metrics.Trend, metrics.Time),
		GasUsed:           r.MustNewMetric("gasper_gas_used", metrics.Trend, metrics.Default),
		BlockGasUsed:      r.MustNewMetric("gasper_block_gas_used", metrics.Trend, metrics.Default),
		Txs:               r.MustNewMetric("gasper_txs", metrics.Trend, metrics.Default),
		BlockTxs:          r.MustNewMetric("gasper_block_txs", metrics.Trend, metrics.Default),
		TPS:               r.MustNewMetric("gasper_tps", metrics.Trend, metrics.Default),
		Mgas:              r.MustNewMetric("gasper_mgas", metrics.Trend, metrics.Default),
		BlockMgas:         r.MustNewMetric("gasper_block_mgas", metrics.Trend, metrics.Default),
		EOA:               r.MustNewMetric("gasper_eoa", metrics.Gauge, metrics.Default),
		BlockTime:         r.MustNewMetric("gasper_block_time", metrics.Trend, metrics.Time),
		BlockPerSec:       r.MustNewMetric("gasper_block_per_sec", metrics.Trend, metrics.Default),
		PoolStatusPending: r.MustNewMetric("gasper_pool_status_pending", metrics.Trend, metrics.Default),
		PoolStatusQueued:  r.MustNewMetric("gasper_pool_status_queued", metrics.Trend, metrics.Default),
	}
}

func ReportReqDurationFromStats(vu modules.VU, m *EthMetrics, clientUID string, call string, dur time.Duration) {
	if vu.State() == nil {
		return
	}

	metrics.PushIfNotDone(vu.Context(), vu.State().Samples, metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: m.RequestDuration,
			Tags:   metrics.NewRegistry().RootTagSet().With("call", call).With("client_uid", clientUID).With("test_uid", TestUID),
		},
		Value: float64(dur / time.Millisecond),
		Time:  time.Now(),
	})
}

func ReportEoaFromStats(vu modules.VU, m *EthMetrics, clientUID string, eoa uint64, txType eth.TransactionType) {
	if vu.State() == nil {
		return
	}

	metrics.PushIfNotDone(vu.Context(), vu.State().Samples, metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: m.EOA,
			Tags:   metrics.NewRegistry().RootTagSet().With("client_uid", clientUID).With("test_uid", TestUID).With("tx_type", string(txType)),
		},
		Value: float64(eoa),
		Time:  time.Now(),
	})
}

func ReportTimeToMineFromStats(vu modules.VU, m *EthMetrics, clientUID string, dur time.Duration) {
	if vu.State() == nil {
		return
	}

	metrics.PushIfNotDone(vu.Context(), vu.State().Samples, metrics.Sample{
		TimeSeries: metrics.TimeSeries{
			Metric: m.TimeToMine,
			Tags:   metrics.NewRegistry().RootTagSet().With("client_uid", clientUID).With("test_uid", TestUID),
		},
		Value: float64(dur / time.Millisecond),
		Time:  time.Now(),
	})
}

func ReportTxPoolStatusFromStats(vu modules.VU, m *EthMetrics, clientUID string, status *eth.PoolStatus) {
	if vu.State() == nil {
		return
	}

	metrics.PushIfNotDone(vu.Context(), vu.State().Samples, metrics.ConnectedSamples{
		Samples: []metrics.Sample{
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.PoolStatusPending,
					Tags:   metrics.NewRegistry().RootTagSet().With("client_uid", clientUID).With("test_uid", TestUID),
				},
				Value: float64(status.Pending),
				Time:  time.Now(),
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.PoolStatusQueued,
					Tags:   metrics.NewRegistry().RootTagSet().With("client_uid", clientUID).With("test_uid", TestUID),
				},
				Value: float64(status.Queued),
				Time:  time.Now(),
			},
		},
	})
}

func ReportBlockMetrics(
	vu modules.VU,
	m *EthMetrics,
	clientUID string,
	block *eth.SlimBlock,
	tps float64,
	mgas float64,
	blockTimestampDiffMili uint64,
	t time.Time,
) {
	if vu.State() == nil {
		return
	}

	rootTS := metrics.NewRegistry().RootTagSet()
	blockTime := time.Unix(int64(block.Timestamp), 0)
	blockNumStr := strconv.FormatUint(block.Number.Uint64(), 10)
	txsLn := strconv.Itoa(len(block.Transactions))
	metrics.PushIfNotDone(vu.Context(), vu.State().Samples, metrics.ConnectedSamples{
		Samples: []metrics.Sample{
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.Txs,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"block":      blockNumStr,
					}),
				},
				Value: float64(len(block.Transactions)),
				Time:  t,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.BlockTxs,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"txs":        txsLn,
					}),
				},
				Value: float64(block.Number.Uint64()),
				Time:  t,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.GasUsed,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"block":      blockNumStr,
					}),
				},
				Value: float64(block.GasUsed),
				Time:  t,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.BlockGasUsed,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"gas_used":   strconv.FormatUint(block.GasUsed, 10),
					}),
				},
				Value: float64(block.Number.Uint64()),
				Time:  t,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.TPS,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"block":      blockNumStr,
					}),
				},
				Value: tps,
				Time:  blockTime,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.Mgas,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"block":      blockNumStr,
					}),
				},
				Value: mgas,
				Time:  t,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.BlockMgas,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"mgas":       strconv.FormatFloat(mgas, 'f', 2, 64),
					}),
				},
				Value: float64(block.Number.Uint64()),
				Time:  t,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.BlockTime,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"block":      blockNumStr,
					}),
				},
				Value: float64(blockTimestampDiffMili),
				Time:  t,
			},
			{
				TimeSeries: metrics.TimeSeries{
					Metric: m.BlockPerSec,
					Tags: rootTS.WithTagsFromMap(map[string]string{
						"client_uid": clientUID,
						"test_uid":   TestUID,
						"block":      blockNumStr,
					}),
				},
				Value: 1.0,
				Time:  blockTime,
			},
		},
	})
}
