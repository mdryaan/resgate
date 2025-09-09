package pool

import (
	"fmt"

	"github.com/mdryaan/resgate/internal/models"
)

func ValidateNew(p *models.Pool) error {
	if p.Name == "" {
		return fmt.Errorf("pool name cannot be empty")
	}
	if p.Total.CPU <= 0 && p.Total.Memory <= 0 && p.Total.GPU <= 0 {
		return fmt.Errorf("pool must have at least one resource (cpu, memory, or gpu)")
	}
	if p.Total.CPU < 0 {
		return fmt.Errorf("cpu must be non-negative")
	}
	if p.Total.Memory < 0 {
		return fmt.Errorf("memory must be non-negative")
	}
	if p.Total.GPU < 0 {
		return fmt.Errorf("gpu must be non-negative")
	}
	return nil
}

func ValidateReservation(pool *models.Pool, req models.Resources) error {
	avail := pool.Available()
	if req.CPU > avail.CPU {
		return fmt.Errorf("insufficient CPU: requested %d, available %d", req.CPU, avail.CPU)
	}
	if req.Memory > avail.Memory {
		return fmt.Errorf("insufficient memory: requested %dMB, available %dMB", req.Memory, avail.Memory)
	}
	if req.GPU > avail.GPU {
		return fmt.Errorf("insufficient GPU: requested %d, available %d", req.GPU, avail.GPU)
	}
	return nil
}
