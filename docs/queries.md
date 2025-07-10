## Filter by Txs

```
from(bucket: "gasper")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r._measurement == "gasper_txs")
  |> filter(fn: (r) => r._field == "value")
  |> filter(fn: (r) => r._value >= 204 and r._value <= 206)
  |> map(fn: (r) => ({
       _time: r._time,
       _value: r._value,
       block: r.block
     }))
```

## Filter by Mgas

```
from(bucket: "gasper")
  |> range(start: v.timeRangeStart, stop: v.timeRangeStop)
  |> filter(fn: (r) => r._measurement == "gasper_block_mgas")
  |> filter(fn: (r) => r._field == "value")
  |> filter(fn: (r) => r._value >= 3812 and r._value <= 3813)
  |> map(fn: (r) => ({
       _time: r._time,
       _value: r._value,
       mgas: r.mgas
     }))
```