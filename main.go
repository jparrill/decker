package main

import (
	"context"
	"fmt"
	"os"
	"os/signal"
	"syscall"

	diagnosticcmd "github.com/jparrill/decker/cmd/diagnostic"
	generatecmd "github.com/jparrill/decker/cmd/generate"
	validatecmd "github.com/jparrill/decker/cmd/validate"
	verifycmd "github.com/jparrill/decker/cmd/verify"
	"github.com/spf13/cobra"
)

func main() {
	cmd := &cobra.Command{
		Use:              "decker",
		SilenceUsage:     true,
		TraverseChildren: true,

		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
			os.Exit(1)
		},
	}

	ctx, cancel := context.WithCancel(context.Background())

	defer cancel()

	cmd.AddCommand(diagnosticcmd.NewCommand())
	cmd.AddCommand(generatecmd.NewCommand())
	cmd.AddCommand(validatecmd.NewCommand())
	cmd.AddCommand(verifycmd.NewCommand())

	sigs := make(chan os.Signal, 1)
	signal.Notify(sigs, syscall.SIGINT)
	go func() {
		<-sigs
		fmt.Fprintln(os.Stderr, "\nAborting...")
		cancel()
	}()

	if err := cmd.ExecuteContext(ctx); err != nil {
		//		fmt.Fprintf(os.Stderr, "%v\n", err)
		os.Exit(1)
	}
}
