package deploy

import (
	"fmt"

	"github.com/Sirupsen/logrus"

	"path/filepath"

	"github.com/aj0strow/netlifyctl/commands/middleware"
	"github.com/aj0strow/netlifyctl/configuration"
	"github.com/aj0strow/netlifyctl/context"
	"github.com/spf13/cobra"
)

type deployCmd struct {
	path   string
	siteId string
}

func Setup() (*cobra.Command, middleware.CommandFunc) {
	cmd := &deployCmd{}
	c := &cobra.Command{
		Use:   "deploy",
		Short: "Deploy your site",
		Long:  "Deploy your site",
	}
	c.Flags().StringVarP(&cmd.path, "path", "p", "", "path to build directory")
	c.Flags().StringVarP(&cmd.siteId, "site-id", "s", "", "target site id to deploy")
	return c, cmd.deploySite
}

func (*deployCmd) deploySite(ctx context.Context, cmd *cobra.Command, args []string) error {
	conf, err := configuration.Load()
	if err != nil {
		return err
	}
	client := context.GetClient(ctx)

	id, err := cmd.Flags().GetString("site-id")
	if err != nil {
		return err
	}

	relPath, err := cmd.Flags().GetString("path")
	if err != nil {
		return err
	}

	path := filepath.Join(conf.Root(), relPath)
	logrus.WithFields(logrus.Fields{"site": id, "path": path}).Debug("deploying site")

	d, err := client.DeploySite(ctx, id, path)
	if err != nil {
		return err
	}

	if len(d.Required) > 0 {
		ready, err := client.WaitUntilDeployReady(ctx, d)
		if err != nil {
			return err
		}
		d = ready
	}
	fmt.Printf("=> Done, your website is live in %s\n", d.URL)

	return nil
}
