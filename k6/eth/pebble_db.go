package eth

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"time"

	"github.com/cockroachdb/pebble"
	"github.com/cockroachdb/pebble/bloom"
)

type PebbleDb struct {
	pth string
	db  *pebble.DB
}

func NewPebbleDb(pth string) (*PebbleDb, error) {
	dir := filepath.Dir(pth)
	if err := os.MkdirAll(dir, 0755); err != nil {
		return nil, err
	}

	if _, err := os.Stat(pth); err == nil {
		date := time.Now().Format("2006-01-02_15-04-05")
		newPath := fmt.Sprintf("%s_%s", strings.TrimSuffix(pth, ".db"), date)
		if err := os.Rename(pth, newPath); err != nil {
			return nil, err
		}
	}

	db, err := pebble.Open(pth, pebbleDbOpt())
	if err != nil {
		return nil, err
	}

	return &PebbleDb{pth: pth, db: db}, nil
}

func (pdb *PebbleDb) Close() error {
	if pdb == nil || pdb.db == nil {
		return nil
	}
	if err := pdb.db.Flush(); err != nil {
		return err
	}
	return pdb.db.Close()
}

func (pdb *PebbleDb) Db() *pebble.DB {
	return pdb.db
}

func (pdb *PebbleDb) GenKey(keys ...string) []byte {
	return []byte(strings.Join(keys, "_"))
}

func pebbleDbOpt() *pebble.Options {
	opt := &pebble.Options{
		MaxOpenFiles:                16,
		MemTableSize:                1<<30 - 1, // Max 1 GB
		MemTableStopWritesThreshold: 2,
		// MaxConcurrentCompactions: func() int { return runtime.NumCPU() },
		Levels: []pebble.LevelOptions{
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
			{TargetFileSize: 2 * 1024 * 1024, FilterPolicy: bloom.FilterPolicy(10)},
		},
	}
	opt.Experimental.ReadSamplingMultiplier = -1

	return opt
}
