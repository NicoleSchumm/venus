package messagepool

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/ipfs/go-datastore"

	"github.com/filecoin-project/go-address"
	"github.com/filecoin-project/venus/pkg/repo"
	"github.com/filecoin-project/venus/venus-shared/types"
)

var (
	ReplaceByFeePercentageMinimum types.Percent = 110
	ReplaceByFeePercentageDefault types.Percent = 125
)

var (
	MemPoolSizeLimitHiDefault = 30000
	MemPoolSizeLimitLoDefault = 20000
	PruneCooldownDefault      = time.Minute
	GasLimitOverestimation    = 1.25

	ConfigKey = datastore.NewKey("/mpool/config")
)

type MpoolConfig struct {
	PriorityAddrs          []address.Address
	SizeLimitHigh          int
	SizeLimitLow           int
	ReplaceByFeeRatio      types.Percent
	PruneCooldown          time.Duration
	GasLimitOverestimation float64
}

func (mc *MpoolConfig) Clone() *MpoolConfig {
	r := new(MpoolConfig)
	*r = *mc
	return r
}

func loadConfig(ctx context.Context, ds repo.Datastore) (*MpoolConfig, error) {
	haveCfg, err := ds.Has(ctx, ConfigKey)
	if err != nil {
		return nil, err
	}

	if !haveCfg {
		return DefaultConfig(), nil
	}

	cfgBytes, err := ds.Get(ctx, ConfigKey)
	if err != nil {
		return nil, err
	}
	cfg := new(MpoolConfig)
	err = json.Unmarshal(cfgBytes, cfg)
	return cfg, err
}

func saveConfig(ctx context.Context, cfg *MpoolConfig, ds repo.Datastore) error {
	cfgBytes, err := json.Marshal(cfg)
	if err != nil {
		return err
	}
	return ds.Put(ctx, ConfigKey, cfgBytes)
}

func (mp *MessagePool) GetConfig() *MpoolConfig {
	mp.cfgLk.Lock()
	defer mp.cfgLk.Unlock()
	return mp.cfg.Clone()
}

func validateConfg(cfg *MpoolConfig) error {
	if cfg.ReplaceByFeeRatio < ReplaceByFeePercentageMinimum {
		return fmt.Errorf("'ReplaceByFeeRatio' is less than required %s < %s",
			cfg.ReplaceByFeeRatio, ReplaceByFeePercentageMinimum)
	}
	if cfg.GasLimitOverestimation < 1 {
		return fmt.Errorf("'GasLimitOverestimation' cannot be less than 1")
	}
	return nil
}

func (mp *MessagePool) SetConfig(ctx context.Context, cfg *MpoolConfig) error {
	if err := validateConfg(cfg); err != nil {
		return err
	}
	cfg = cfg.Clone()

	mp.cfgLk.Lock()
	mp.cfg = cfg
	err := saveConfig(ctx, cfg, mp.ds)
	if err != nil {
		log.Warnf("error persisting mpool config: %s", err)
	}
	mp.cfgLk.Unlock()

	return nil
}

func DefaultConfig() *MpoolConfig {
	return &MpoolConfig{
		SizeLimitHigh:          MemPoolSizeLimitHiDefault,
		SizeLimitLow:           MemPoolSizeLimitLoDefault,
		ReplaceByFeeRatio:      ReplaceByFeePercentageDefault,
		PruneCooldown:          PruneCooldownDefault,
		GasLimitOverestimation: GasLimitOverestimation,
	}
}
