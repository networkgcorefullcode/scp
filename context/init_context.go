package context

var (
	scpContext = ScpContext{}
)

// Create new SCP context
func SCP_Self() *ScpContext {
	return &scpContext
}

// Initialize SCP context
func InitScpContext(self *ScpContext) {
	scpContext = ScpContext{
		Rcvd: false,
	}
}
