// Package format
package format

import (
	"bytes"
	"encoding/json"

	"github.com/fatih/color"
	"github.com/henrywhitakercommify/netpol-fqdns/internal/collector"
	"github.com/rodaine/table"
)

func JSON(rules []collector.Rule) ([]byte, error) {
	by, err := json.Marshal(rules)
	if err != nil {
		return nil, err
	}
	return by, nil
}

func Table(rules []collector.Rule) ([]byte, error) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	out := &bytes.Buffer{}

	tbl := table.New("Name", "Namespace", "FQDN")
	tbl.WithHeaderFormatter(headerFmt).WithFirstColumnFormatter(columnFmt).WithWriter(out)

	for _, r := range rules {
		for _, f := range r.FQDNs {
			tbl.AddRow(r.Name, r.Namespace, f)
		}
	}

	return out.Bytes(), nil
}
