package service

import (
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/fsnotify/fsnotify"
	"github.com/networkgcorefullcode/scp/context"
	"github.com/networkgcorefullcode/scp/factory"
	"github.com/networkgcorefullcode/scp/logger"
	"github.com/networkgcorefullcode/scp/proxy"
	"github.com/spf13/viper"
	"github.com/urfave/cli/v3"
	"go.uber.org/zap"
	"go.uber.org/zap/zapcore"
)

type SCP struct{}

type (
	// Config information.
	Config struct {
		cfg string
	}
)

var config Config

var scpCLi = []cli.Flag{
	&cli.StringFlag{
		Name:     "cfg",
		Usage:    "scp config file",
		Required: true,
	},
}

func (*SCP) GetCliCmd() (flags []cli.Flag) {
	return scpCLi
}

func (scp *SCP) Initialize(c *cli.Command) error {
	config = Config{
		cfg: c.String("cfg"),
	}

	absPath, err := filepath.Abs(config.cfg)
	if err != nil {
		logger.CfgLog.Errorln(err)
		return err
	}

	if err := factory.InitConfigFactory(absPath); err != nil {
		return err
	}

	scp.setLogLevel()

	if err := factory.CheckConfigVersion(); err != nil {
		return err
	}

	if _, err := os.Stat(absPath); err == nil {
		viper.SetConfigFile(absPath)
		viper.SetConfigType("yaml")
		err = viper.ReadInConfig() // Find and read the config file
		if err != nil {            // Handle errors reading the config file
			return err
		}
	} else if os.IsNotExist(err) {
		logger.AppLog.Errorln("file %s does not exists", absPath)
		return err
	}

	factory.ScpConfig.CfgLocation = absPath

	return nil
}

func (scp *SCP) setLogLevel() {
	if factory.ScpConfig.Logger == nil {
		logger.InitLog.Warnln("SCP config without log level setting")
		return
	}

	if factory.ScpConfig.Logger.SCP_Proxy != nil {
		if factory.ScpConfig.Logger.SCP_Proxy.DebugLevel != "" {
			if level, err := zapcore.ParseLevel(factory.ScpConfig.Logger.SCP_Proxy.DebugLevel); err != nil {
				logger.InitLog.Warnf("SCP_Proxy Log level [%s] is invalid, set to [info] level",
					factory.ScpConfig.Logger.SCP_Proxy.DebugLevel)
				logger.SetLogLevel(zap.InfoLevel)
			} else {
				logger.InitLog.Infof("SCP_Proxy Log level is set to [%s] level", level)
				logger.SetLogLevel(level)
			}
		} else {
			logger.InitLog.Warnln("SCP_Proxy Log level not set. Default set to [info] level")
			logger.SetLogLevel(zap.InfoLevel)
		}
	}
}

func (scp *SCP) WatchConfig() {
	viper.WatchConfig()
	viper.OnConfigChange(func(e fsnotify.Event) {
		logger.AppLog.Infoln("config file changed:", e.Name)
		if err := factory.UpdateConfig(factory.ScpConfig.CfgLocation); err != nil {
			logger.AppLog.Errorln("error in loading updated configuration")
		} else {
			self := context.SCP_Self()
			context.InitScpContext(self)
			logger.AppLog.Infoln("successfully updated configuration")
		}
	})
}

func (scp *SCP) Start() {
	logger.AppLog.Infoln("SCP started are starting")

	// new scp context
	self := context.SCP_Self()
	context.InitScpContext(self)

	// Start the proxy server http
	go proxy.Start_Proxy_Server(self.PortHttp)

	signalChannel := make(chan os.Signal, 1)
	signal.Notify(signalChannel, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-signalChannel
		scp.Terminate()
		os.Exit(0)
	}()
}
