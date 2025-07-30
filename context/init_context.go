package context

import (
	"github.com/google/uuid"
	"github.com/networkgcorefullcode/scp/factory"
	"github.com/networkgcorefullcode/scp/logger"
)

var (
	scpContext = ScpContext{}
)

// Create new SCP context
func SCP_Self() *ScpContext {
	return &scpContext
}

// Initialize SCP context
func InitScpContext(context *ScpContext) {
	config := factory.ScpConfig
	logger.UtilLog.Infof("scpconfig Info: Version[%s] Description[%s]", config.Info.Version, config.Info.Description)
	configuration := config.Configuration

	if context.ScpId == "" {
		logger.UtilLog.Infoln("context.ScpId empty, generating a new")
		context.ScpId = uuid.New().String()
	} else {
		logger.UtilLog.Infoln("context.ScpId is present: ", context.ScpId)
	}

	if configuration.ScpName != "" {
		logger.UtilLog.Infoln("Has been added ScpName from file config: ", configuration.ScpName)
		context.ScpName = configuration.ScpName
	} else {
		logger.UtilLog.Warnln("ScpName not defined in file config")
	}

	if configuration.ScpDBName != "" {
		logger.UtilLog.Infoln("Has been added ScpDBName from file config: ", configuration.ScpDBName)
		context.ScpDBName = configuration.ScpDBName
	} else {
		logger.UtilLog.Warnln("ScpDBName not defined in file config")
	}

	if configuration.PortHttp != 0 {
		logger.UtilLog.Infoln("Has been added PortHttp from file config: ", configuration.PortHttp)
		context.PortHttp = configuration.PortHttp
	} else {
		logger.UtilLog.Warnln("PortHttp not defined in file config")
		context.PortHttp = factory.SCP_HTTP_PORT // Default port
	}
}
