package irc

const (
	Protocol string = "irc"

	PermPRIVMSG int = 4
	PermACTION

	PermJOIN int = 8
	PermKICK
	PermNICK
	PermINVITE
	PermMODE

	PermCONNECT int = 16
	PermDISCONNECT

	PermIRCBridge int = 1 << 30
)
