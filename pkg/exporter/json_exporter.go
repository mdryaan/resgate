package exporter

import "encoding/json"

type JSONExporter struct{}

func (e *JSONExporter) Export(r *Report) ([]byte, error) {
	return json.MarshalIndent(r, "", "  ")
}
