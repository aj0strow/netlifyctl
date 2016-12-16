package commands

import (
	"github.com/aj0strow/netlifyctl/commands/deploy"
	"github.com/aj0strow/netlifyctl/commands/middleware"
	"github.com/aj0strow/netlifyctl/commands/sites"
	"github.com/spf13/cobra"
)

func setupRunE(cmd *cobra.Command, f middleware.CommandFunc, m []middleware.Middleware) *cobra.Command {
	cmd.RunE = middleware.NewRunFunc(f, m)
	return cmd
}

func addCommands() {
	middlewares := []middleware.Middleware{
		middleware.AuthMiddleware,
		middleware.ClientMiddleware,
		middleware.LoggingMiddleware,
	}

	sCmd, sFunc := sites.Setup()
	rootCmd.AddCommand(setupRunE(sCmd, sFunc, middlewares))

	dCmd, dFunc := deploy.Setup()
	rootCmd.AddCommand(setupRunE(dCmd, dFunc, middlewares))

	rootCmd.AddCommand(versionCmd)
}
