package exporter

import (
	"fmt"

	"github.com/mdryaan/resgate/internal/models"
)

type Report struct {
	Pools        []*models.Pool
	Tenants      []*models.Tenant
	Reservations []*models.Reservation
	History      []*models.Reservation
}

type Exporter interface {
	Export(r *Report) ([]byte, error)
}

func New(format string) (Exporter, error) {
	switch format {
	case "json":
		return &JSONExporter{}, nil
	case "csv":
		return &CSVExporter{}, nil
	default:
		return nil, fmt.Errorf("unsupported export format: %s", format)
	}
}
