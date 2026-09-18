// Package format
package format

import (
	"encoding/csv"
	"encoding/json"
	"strings"

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

func CSV(rules []collector.Rule) ([]byte, error) {
	out := &strings.Builder{}

	w := csv.NewWriter(out)
	if err := w.Write([]string{"name", "namespace", "fqdn"}); err != nil {
		return nil, err
	}

	for _, r := range rules {
		for _, f := range r.FQDNs {
			if err := w.Write([]string{r.Name, r.Namespace, f}); err != nil {
				return nil, err
			}
		}
	}

	w.Flush()
	if err := w.Error(); err != nil {
		return nil, err
	}

	return []byte(out.String()), nil
}

func Table(rules []collector.Rule) ([]byte, error) {
	headerFmt := color.New(color.FgGreen, color.Underline).SprintfFunc()
	columnFmt := color.New(color.FgYellow).SprintfFunc()

	out := &strings.Builder{}

	tbl := table.New("Name", "Namespace", "FQDN").
		WithHeaderFormatter(headerFmt).
		WithFirstColumnFormatter(columnFmt).
		WithWriter(out)

	for _, r := range rules {
		for _, f := range r.FQDNs {
			tbl.AddRow(r.Name, r.Namespace, f)
		}
	}

	tbl.Print()

	return []byte(out.String()), nil
}
