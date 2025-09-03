package models

import "time"

type Pool struct {
	Name      string    `json:"name"`
	Total     Resources `json:"total"`
	Reserved  Resources `json:"reserved"`
	CreatedAt time.Time `json:"created_at"`
}

func (p *Pool) Available() Resources {
	return p.Total.Sub(p.Reserved)
}

func (p *Pool) CanFit(r Resources) bool {
	return r.Fits(p.Available())
}

func (p *Pool) Allocate(r Resources) {
	p.Reserved = p.Reserved.Add(r)
}

func (p *Pool) Release(r Resources) {
	p.Reserved = p.Reserved.Sub(r)
	if p.Reserved.CPU < 0 {
		p.Reserved.CPU = 0
	}
	if p.Reserved.Memory < 0 {
		p.Reserved.Memory = 0
	}
	if p.Reserved.GPU < 0 {
		p.Reserved.GPU = 0
	}
}

func (p *Pool) CPUPercent() float64 {
	if p.Total.CPU == 0 {
		return 0
	}
	return float64(p.Reserved.CPU) / float64(p.Total.CPU) * 100
}

func (p *Pool) MemoryPercent() float64 {
	if p.Total.Memory == 0 {
		return 0
	}
	return float64(p.Reserved.Memory) / float64(p.Total.Memory) * 100
}

func (p *Pool) GPUPercent() float64 {
	if p.Total.GPU == 0 {
		return 0
	}
	return float64(p.Reserved.GPU) / float64(p.Total.GPU) * 100
}
