package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	"github.com/henrywhitakercommify/netpol-fqdns/internal/collector"
	"github.com/henrywhitakercommify/netpol-fqdns/internal/format"
	"github.com/spf13/pflag"
)

var (
	kubeconfig   string
	outputFormat string
)

func main() {
	ctx, cancel := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer cancel()

	flags := pflag.NewFlagSet("", pflag.ContinueOnError)
	flags.StringVar(&kubeconfig, "kubeconfig", "~/.kube/config", "The path to the kubeconfig file")
	flags.StringVar(&outputFormat, "format", "table", "The format to output the rules as")

	if err := flags.Parse(os.Args[1:]); err != nil {
		panic(err)
	}

	rules, err := collector.Collect(ctx, kubeconfig)
	if err != nil {
		panic(err)
	}

	var output []byte
	switch outputFormat {
	case "json":
		output, err = format.JSON(rules)
	case "table":
		output, err = format.Table(rules)
	default:
		panic(fmt.Sprintf("unknown format %s", outputFormat))
	}
	if err != nil {
		panic(err)
	}
	fmt.Println(string(output))
}
