package worker

import (
	"fmt"

	"gitlab.com/jgillich/autominer/hardware"
	"gitlab.com/jgillich/autominer/stratum"
)

var workers = map[string]workerFactory{}

type workerFactory func(Config) Worker

func New(coin string, config Config) (Worker, error) {
	factory, ok := workers[coin]
	if !ok {
		return nil, fmt.Errorf("worker for coin '%v' does not exist")
	}

	return factory(config), nil
}

func List() map[string]Capabilities {
	list := map[string]Capabilities{}

	for name, factory := range workers {
		list[name] = factory(Config{}).Capabilities()
	}

	return list
}

type Worker interface {
	Work() error
	Capabilities() Capabilities
}

type Capabilities struct {
	CPU    bool
	OpenCL bool
	CUDA   bool
}

type Config struct {
	Stratum *stratum.Client
	Donate  int
	CPUSet  []CPUConfig
	GPUSet  []GPUConfig
}

type CPUConfig struct {
	Threads int
	CPU     hardware.CPU
}

type GPUConfig struct {
	Intensity int
	GPU       hardware.GPU
}

type Stats struct {
	CPUStats []CPUStats `json:"cpu_stats"`
	GPUStats []GPUStats `json:"gpu_stats"`
}

type CPUStats struct {
	Index    int     `json:"index"`
	Hashrate float32 `json:"hashrate"`
}

type GPUStats struct {
	Index    int     `json:"index"`
	Hashrate float32 `json:"hashrate"`
}