package worker

import (
	crand "crypto/rand"
	"encoding/hex"
	"fmt"
	"math"
	"math/big"
	"math/rand"
	"strings"
	"time"

	metrics "github.com/armon/go-metrics"
	"gitlab.com/blockforge/blockforge/algo/ethash"
	"gitlab.com/blockforge/blockforge/log"
)

var maxUint256 = new(big.Int).Exp(big.NewInt(2), big.NewInt(256), big.NewInt(0))

type Ethash struct {
	Work   <-chan *ethash.Work
	Shares chan<- ethash.Share

	config Config
	// random source for nonces
	rand *rand.Rand

	metrics *metrics.Metrics
}

func (worker *Ethash) Configure(config Config) error {
	seed, err := crand.Int(crand.Reader, big.NewInt(math.MaxInt64))
	if err != nil {
		return err
	}

	worker.config = config
	worker.rand = rand.New(rand.NewSource(seed.Int64()))

	worker.metrics = config.Metrics

	return nil
}

func (worker *Ethash) Start() error {
	totalThreads := len(worker.config.CLDevices)
	for _, c := range worker.config.Processors {
		totalThreads += c.Threads
	}

	workChannels := make([]chan *ethash.Work, totalThreads)
	for i := 0; i < totalThreads; i++ {
		workChannels[i] = make(chan *ethash.Work, 1)
		defer func(i int) {
			close(workChannels[i])
		}(i)
	}

	var light *ethash.Light
	var full *ethash.Full
	defer func() {
		if light != nil {
			light.Release()
		}
		if full != nil {
			full.Release()
		}
	}()

	var seedhash string

	for work := range worker.Work {
		if seedhash != work.Seedhash {
			seedhash = work.Seedhash

			seedhashBytes, err := hex.DecodeString(strings.TrimPrefix(work.Seedhash, "0x"))
			if err != nil {
				return err
			}

			// when DAG changes, we shutdown and recreate all threads
			for i := 0; i < totalThreads; i++ {
				close(workChannels[i])
				workChannels[i] = make(chan *ethash.Work, 1)
				if light != nil {
					light.Release()
					light = nil
				}
				if full != nil {
					full.Release()
					full = nil
				}
			}

			light, err = ethash.NewLight(seedhashBytes)
			if err != nil {
				return err
			}

			if len(worker.config.Processors) > 0 {
				log.Info("DAG is being initialized, this may take a while")
				full, err = ethash.NewFull(light)
				if err != nil {
					return err
				}
				log.Info("DAG initialized")
			}

			for cpuIndex, conf := range worker.config.Processors {
				for i := 0; i < conf.Threads; i++ {
					key := []string{"cpu", fmt.Sprintf("%v", cpuIndex), fmt.Sprintf("%v", i)}
					go worker.thread(key, full, workChannels[len(worker.config.CLDevices)+i])
				}
			}

			if len(worker.config.CLDevices) > 0 {
				for i, d := range worker.config.CLDevices {
					cl, err := newEthashCL(d, light)
					if err != nil {
						return err
					}
					key := []string{"opencl", fmt.Sprintf("%v", d.Device.Platform.Index), fmt.Sprintf("%v", d.Device.Index)}

					go worker.clThread(key, cl, workChannels[i])
				}
			}
		}

		for _, ch := range workChannels {
			ch <- work
		}
	}

	return nil
}

func (worker *Ethash) thread(key []string, full *ethash.Full, workChan chan *ethash.Work) {
	work := <-workChan
	var ok bool

	nonce := uint64(worker.rand.Uint32())
	stepping := uint64(64 * 1024)

	for {
		select {
		case work, ok = <-workChan:
			if !ok {
				return
			}
			nonce = uint64(worker.rand.Uint32())

		default:
			start := time.Now()
			if err := work.VerifyRange(full, nonce, stepping, worker.Shares); err != nil {
				workerError(err)
			}
			nonce += stepping
			worker.metrics.IncrCounter(key, float32(float64(stepping)/time.Since(start).Seconds()))
		}
	}
}

func (worker *Ethash) clThread(key []string, cl *ethashCL, workChan chan *ethash.Work) {
	defer cl.Release()

	work := <-workChan
	if err := cl.Update(work.Header, work.Target); err != nil {
		workerError(err)
	}

	var ok bool
	var results [2]uint32

	nonce := uint64(worker.rand.Uint32())

	for {
		select {
		case work, ok = <-workChan:
			if !ok {
				return
			}
			if err := cl.Update(work.Header, work.Target); err != nil {
				workerError(err)
			}
			nonce = uint64(worker.rand.Uint32())

		default:
			start := time.Now()

			if err := cl.Run(work.ExtraNonce+nonce, results); err != nil {
				workerError(err)
			}
			if results[0] > 0 {
				worker.Shares <- ethash.Share{
					JobId: work.JobId,
					Nonce: nonce + uint64(results[1]),
				}
			}
			nonce += uint64(cl.globalWorkSize)
			worker.metrics.IncrCounter(key, float32(float64(cl.globalWorkSize)/time.Since(start).Seconds()))
		}
	}

}

func (w *Ethash) Capabilities() Capabilities {
	return Capabilities{
		CPU:    true,
		OpenCL: true,
		CUDA:   false,
	}
}
