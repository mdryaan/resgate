package output

import (
	"os"

	"github.com/olekukonko/tablewriter"
)

func NewTable(headers []string) *tablewriter.Table {
	colors := make([]tablewriter.Colors, len(headers))
	for i := range colors {
		colors[i] = tablewriter.Colors{tablewriter.Bold, tablewriter.FgCyanColor}
	}
	t := tablewriter.NewWriter(os.Stdout)
	t.SetHeader(headers)
	t.SetBorder(true)
	t.SetRowLine(false)
	t.SetHeaderAlignment(tablewriter.ALIGN_LEFT)
	t.SetAlignment(tablewriter.ALIGN_LEFT)
	t.SetHeaderColor(colors...)
	return t
}
