package ignitecmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/ignite/cli/ignite/pkg/cliui"
	"github.com/ignite/cli/ignite/pkg/cliui/icons"
	"github.com/ignite/cli/ignite/pkg/gomodulepath"
	"github.com/ignite/cli/ignite/pkg/placeholder"
	"github.com/ignite/cli/ignite/pkg/xgenny"
	"github.com/ignite/cli/ignite/templates/app"
)

func NewDoctor() *cobra.Command {
	return &cobra.Command{
		Use:    "doctor",
		Short:  "Try to fix things",
		Hidden: true,
		RunE: func(cmd *cobra.Command, args []string) error {
			session := cliui.New(cliui.StartSpinner())
			defer session.End()

			if err := doctorDepTools(session); err != nil {
				return fmt.Errorf("dep tools: %w", err)
			}

			return nil
		},
	}
}

// doctorDepTools ensures that tools.go is present and properly set.
func doctorDepTools(session *cliui.Session) error {
	const toolsGoFile = "tools/tools.go"
	session.StartSpinner(fmt.Sprintf("Checking %s", toolsGoFile))
	_, err := os.Stat(toolsGoFile)
	switch {
	case err == nil:
		// tools.go exists
		// TODO ensure it has the required dependencies
		session.Printf("%s Checked %s\n", icons.OK, toolsGoFile)

	case os.IsNotExist(err):
		// create tools.go
		pathInfo, err := gomodulepath.ParseAt(".")
		if err != nil {
			return err
		}
		g, err := app.NewToolsGoGenerator(&app.Options{
			ModulePath:       pathInfo.RawPath,
			AppName:          pathInfo.Package,
			BinaryNamePrefix: pathInfo.Root,
		})
		if err != nil {
			return fmt.Errorf("generator: %w", err)
		}
		_, err = xgenny.RunWithValidation(placeholder.New(), g)
		if err != nil {
			return fmt.Errorf("xgenny.run: %w", err)
		}
		session.Printf("%s %s\n", createPrefix, toolsGoFile)
		// TODO run go mod tidy or go install dependencies
		// TODO remove go get from install dependencies fonctions?
		// TODO ask to run `ignite chain doctor` is install dependency fail

	default:
		return err
	}
	return nil
}
