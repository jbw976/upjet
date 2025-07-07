// SPDX-FileCopyrightText: 2023 The Crossplane Authors <https://crossplane.io>
//
// SPDX-License-Identifier: Apache-2.0

package main

import (
	"os"

	"github.com/alecthomas/kong"
	"github.com/crossplane/crossplane-runtime/pkg/logging"
	"sigs.k8s.io/controller-runtime/pkg/log/zap"

	"github.com/crossplane/upjet/cmd/upjet/batch"
)

var _ = kong.Must(&cli{})

type (
	verboseFlag bool
)

// BeforeApply ensures that when the verbose flag is set that we set up and use
// a debug level logger
func (v verboseFlag) BeforeApply(ctx *kong.Context) error { //nolint:unparam // BeforeApply requires this signature with error return type (even though it's not used).
	logger := logging.NewLogrLogger(zap.New(zap.UseDevMode(true)))
	ctx.BindTo(logger, (*logging.Logger)(nil))
	return nil
}

// upjet CLI top level commands and flags
type cli struct {
	// commands - keep these alphabetized
	Batch batch.BatchCmd `cmd:"" help:"Batch build and push a family of service-scoped provider packages."`
	Help  helpCmd        `cmd:"" help:"Show help."`

	// flags - keep these alphabetized too
	Verbose verboseFlag `help:"Print verbose logging statements." name:"verbose"`
}

// helpCmd ensures that the general help content is shown.
type helpCmd struct{}

func (h *helpCmd) Run(ctx *kong.Context) error {
	_, err := ctx.Parse([]string{"--help"})
	return err
}

const helpDescription = `The Upjet CLI.

Please report issues and feature requests at https://github.com/crossplane/upjet.`

func main() {
	logger := logging.NewNopLogger()

	// set up the kong runtime to parse and run our CLI logic
	parser := kong.Must(&cli{},
		kong.Name("upjet"),
		kong.Description(helpDescription),
		// Binding a variable to kong context makes it available to all commands
		// at runtime.
		kong.BindTo(logger, (*logging.Logger)(nil)),
		kong.ConfigureHelp(kong.HelpOptions{
			FlagsLast:      true,
			Compact:        true,
			WrapUpperBound: 80,
		}),
		kong.UsageOnError())

	if len(os.Args) == 1 {
		// no args provided, show help
		_, err := parser.Parse([]string{"--help"})
		parser.FatalIfErrorf(err)
		return
	}

	// parse and run the given commands and arguments
	ctx, err := parser.Parse(os.Args[1:])
	parser.FatalIfErrorf(err)

	err = ctx.Run()
	ctx.FatalIfErrorf(err)
}
