package service

import "github.com/networkgcorefullcode/scp/logger"

func (scp *SCP) Terminate() {
	logger.AppLog.Infoln("SCP is terminating...")

	// Perform any necessary cleanup here
	// For example, close database connections, stop background tasks, etc.

	logger.AppLog.Infoln("SCP terminated successfully")
}
