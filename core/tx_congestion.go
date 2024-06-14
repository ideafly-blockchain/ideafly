package core

import (
	"github.com/ethereum/go-ethereum/core/txpool"
	"math/big"
	"sort"
	"sync"
	"time"

	"github.com/ethereum/go-ethereum/core/types"
	"github.com/ethereum/go-ethereum/log"
	"github.com/ethereum/go-ethereum/metrics"
)

var (
	congestionMeter = metrics.NewRegisteredMeter("txpool/congestion", nil)
)

var oneGwei = big.NewInt(1e9)

var DefaultCongestionConfig = TxCongestionConfig{
	PeriodsSecs:         3,
	CongestionSecs:      15,
	UnderPricedFactor:   3,
	PendingFactor:       1,
	MaxValidPendingSecs: 300,
}

type TxCongestionConfig struct {
	PeriodsSecs       int // how many seconds to do a reflesh to the congestion index
	CongestionSecs    int // how many seconds for a tx pending that will meams a tx-congestion.
	UnderPricedFactor int
	PendingFactor     int

	MaxValidPendingSecs int //
}

func (c *TxCongestionConfig) sanity() TxCongestionConfig {
	cfg := *c
	if cfg.PeriodsSecs < 1 {
		log.Info("CongestionConfig sanity PeriodsSecs", "old", cfg.PeriodsSecs, "new", DefaultCongestionConfig.PeriodsSecs)
		cfg.PeriodsSecs = DefaultCongestionConfig.PeriodsSecs
	}
	if cfg.CongestionSecs < 3 {
		log.Info("CongestionConfig sanity CongestionSecs", "old", cfg.CongestionSecs, "new", DefaultCongestionConfig.CongestionSecs)
		cfg.CongestionSecs = DefaultCongestionConfig.CongestionSecs
	}
	if cfg.UnderPricedFactor < 1 {
		log.Info("CongestionConfig sanity UnderPricedFactor", "old", cfg.UnderPricedFactor, "new", DefaultCongestionConfig.UnderPricedFactor)
		cfg.UnderPricedFactor = DefaultCongestionConfig.UnderPricedFactor
	}
	if cfg.PendingFactor < 1 {
		log.Info("CongestionConfig sanity PendingFactor", "old", cfg.PendingFactor, "new", DefaultCongestionConfig.PendingFactor)
		cfg.PendingFactor = DefaultCongestionConfig.PendingFactor
	}
	if cfg.MaxValidPendingSecs <= cfg.CongestionSecs {
		log.Info("CongestionConfig sanity MaxValidPendingSecs", "old", cfg.MaxValidPendingSecs, "new", DefaultCongestionConfig.MaxValidPendingSecs)
		cfg.MaxValidPendingSecs = DefaultCongestionConfig.MaxValidPendingSecs
	}
	return cfg
}

// TxCongestionRecorder try to give a quantitative index to reflects the tx congestion.
type TxCongestionRecorder struct {
	cfg  TxCongestionConfig
	pool *txpool.TxPool
	head *types.Header

	underPricedCounter *underPricedCounter
	currentCongestion  int

	congestionLock sync.RWMutex

	quit        chan struct{}
	chainHeadCh chan *types.Header
}

func NewTxCongestionRecorder(cfg TxCongestionConfig, pool *txpool.TxPool) *TxCongestionRecorder {
	cfg = (&cfg).sanity()

	recorder := &TxCongestionRecorder{
		cfg:                cfg,
		pool:               pool,
		underPricedCounter: newUnderPricedCounter(cfg.PeriodsSecs),
		quit:               make(chan struct{}),
		chainHeadCh:        make(chan *types.Header, 1),
	}

	go recorder.updateLoop()

	return recorder
}

// Stop stops the loop goroutines of this congestion index
func (recorder *TxCongestionRecorder) Stop() {
	recorder.underPricedCounter.Stop()
	close(recorder.quit)
}

// CongestionRecord returns the current congestion record
func (recorder *TxCongestionRecorder) CongestionRecord() int {
	recorder.congestionLock.RLock()
	defer recorder.congestionLock.RUnlock()
	return recorder.currentCongestion
}

func (recorder *TxCongestionRecorder) updateLoop() {
	tick := time.NewTicker(time.Second * time.Duration(recorder.cfg.PeriodsSecs))
	defer tick.Stop()

	for {
		select {
		case h := <-recorder.chainHeadCh:
			recorder.head = h
		case <-tick.C:
			d := recorder.underPricedCounter.Sum()
			pendings := recorder.pool.Pending(false)
			if d == 0 && len(pendings) == 0 {
				break
			}
			// flatten
			var p int
			max := recorder.cfg.MaxValidPendingSecs
			congestionSecs := recorder.cfg.CongestionSecs
			maxGas := uint64(10000000)
			if recorder.head != nil {
				maxGas = (recorder.head.GasLimit / 10) * 6
			}
			durs := make([]time.Duration, 0, 1024)
			for _, txs := range pendings {
				for _, tx := range txs {
					// filtering
					if tx.GasPrice().Cmp(oneGwei) < 0 ||
						tx.Gas() > maxGas {
						continue
					}

					dur := time.Since(tx.LocalSeenTime())
					sec := int(dur / time.Second)
					if sec > max {
						continue
					}

					durs = append(durs, dur)
					if sec >= congestionSecs {
						p += sec / congestionSecs
					}
				}
			}
			nTotal := len(durs)

			if nTotal == 0 {
				p = 0
			} else {
				p = 100 * p / nTotal
			}

			idx := d*recorder.cfg.UnderPricedFactor + p*recorder.cfg.PendingFactor
			recorder.congestionLock.Lock()
			recorder.currentCongestion = idx
			recorder.congestionLock.Unlock()
			congestionMeter.Mark(int64(idx))

			var dists []time.Duration
			sort.Slice(durs, func(i, j int) bool {
				return durs[i] < durs[j]
			})
			if nTotal > 10 {
				dists = append(dists, durs[0])
				for i := 1; i < 10; i++ {
					dists = append(dists, durs[nTotal*i/10])
				}
				dists = append(dists, durs[nTotal-1])
			} else {
				dists = durs
			}

			log.Trace("TxCongestion", "congestion", idx, "d", d, "p", p, "n", nTotal, "dists", dists)
		case <-recorder.quit:
			return
		}
	}
}

func (recorder *TxCongestionRecorder) UpdateHeader(h *types.Header) {
	recorder.chainHeadCh <- h
}

func (recorder *TxCongestionRecorder) UnderPricedInc() {
	recorder.underPricedCounter.Inc()
}

type underPricedCounter struct {
	counts  []int // the length of this slice is 2 times of periodSecs
	periods int   //how many periods to cache, each period cache records of 0.5 seconds.
	idx     int   //current index
	sum     int   //current sum

	inCh       chan struct{}
	quit       chan struct{}
	queryCh    chan struct{}
	queryResCh chan int
}

func newUnderPricedCounter(periodSecs int) *underPricedCounter {
	c := &underPricedCounter{
		counts:     make([]int, 2*periodSecs),
		periods:    2 * periodSecs,
		inCh:       make(chan struct{}, 10),
		quit:       make(chan struct{}),
		queryCh:    make(chan struct{}),
		queryResCh: make(chan int),
	}
	go c.loop()
	return c
}

func (c *underPricedCounter) loop() {
	tick := time.NewTicker(500 * time.Millisecond)
	defer tick.Stop()

	for {
		select {
		case <-tick.C:
			c.idx = (c.idx + 1) % c.periods
			c.sum -= c.counts[c.idx]
			c.counts[c.idx] = 0
		case <-c.inCh:
			c.counts[c.idx]++
			c.sum++
		case <-c.queryCh:
			c.queryResCh <- c.sum
		case <-c.quit:
			return
		}
	}
}

func (c *underPricedCounter) Sum() int {
	var sum int
	c.queryCh <- struct{}{}
	sum = <-c.queryResCh
	return sum
}

func (c *underPricedCounter) Inc() {
	c.inCh <- struct{}{}
}

func (c *underPricedCounter) Stop() {
	close(c.quit)
}
