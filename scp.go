package main

// Here, we are importing the necessary packages to handle HTTP requests and responses,
// read and write data, and log any errors that occur.
import (
	"context"
	"fmt"
	"os"

	"github.com/networkgcorefullcode/scp/logger"
	"github.com/networkgcorefullcode/scp/service"
	"github.com/urfave/cli/v3"
)

var SCP = &service.SCP{}

func main() {
	app := &cli.Command{}
	app.Name = "scp"
	logger.AppLog.Infoln(app.Name)
	app.Usage = "Service Communication Proxy (SCP) for Aether"
	app.UsageText = "scp -cfg <scp_config_file.conf>"
	app.Action = action
	app.Flags = SCP.GetCliCmd()
	if err := app.Run(context.Background(), os.Args); err != nil {
		logger.AppLog.Fatalf("SCP run error: %v", err)
	}
}

func action(ctx context.Context, c *cli.Command) error {
	if err := SCP.Initialize(c); err != nil {
		logger.CfgLog.Errorf("%+v", err)
		return fmt.Errorf("failed to initialize")
	}

	SCP.WatchConfig()

	SCP.Start()

	return nil
}
