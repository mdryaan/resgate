package pool

import "github.com/mdryaan/resgate/internal/models"

type UtilStatus string

const (
	StatusHealthy  UtilStatus = "healthy"
	StatusWarning  UtilStatus = "warning"
	StatusExhausted UtilStatus = "exhausted"
)

func UtilizationStatus(pool *models.Pool) UtilStatus {
	cpu := pool.CPUPercent()
	mem := pool.MemoryPercent()
	gpu := pool.GPUPercent()

	max := cpu
	if mem > max {
		max = mem
	}
	if gpu > max {
		max = gpu
	}

	switch {
	case max >= 90:
		return StatusExhausted
	case max >= 60:
		return StatusWarning
	default:
		return StatusHealthy
	}
}
