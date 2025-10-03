package exporter

import (
	"bytes"
	"encoding/csv"
	"fmt"
	"strconv"

	"github.com/mdryaan/resgate/internal/utils"
)

type CSVExporter struct{}

func (e *CSVExporter) Export(r *Report) ([]byte, error) {
	var buf bytes.Buffer
	w := csv.NewWriter(&buf)

	w.Write([]string{"section", "id", "tenant", "pool", "cpu", "memory", "gpu", "status", "created_at", "expires_at"})

	for _, res := range r.History {
		expiresAt := ""
		if res.ExpiresAt != nil {
			expiresAt = utils.FormatTime(*res.ExpiresAt)
		}
		w.Write([]string{
			"history",
			res.ID,
			res.TenantName,
			res.PoolName,
			strconv.Itoa(res.Resources.CPU),
			strconv.Itoa(res.Resources.Memory),
			strconv.Itoa(res.Resources.GPU),
			string(res.Status),
			utils.FormatTime(res.CreatedAt),
			expiresAt,
		})
	}

	for _, p := range r.Pools {
		w.Write([]string{
			"pool",
			p.Name,
			"",
			p.Name,
			strconv.Itoa(p.Total.CPU),
			strconv.Itoa(p.Total.Memory),
			strconv.Itoa(p.Total.GPU),
			fmt.Sprintf("cpu=%.0f%% mem=%.0f%%", p.CPUPercent(), p.MemoryPercent()),
			utils.FormatTime(p.CreatedAt),
			"",
		})
	}

	w.Flush()
	return buf.Bytes(), w.Error()
}
