// Stores things like message codes, etc.

package zif

var (
	// Protocol header, so we know this is a zif client.
	// Version should follow.
	ProtoZif     = 0x7a66
	ProtoVersion = 0x0000 //version 0 atm :D (change when spec is stable)

	ProtoHeader = 0x0000

	// inform a peer on the status of the latest request
	ProtoOk        = 0x0001
	ProtoNo        = 0x0002
	ProtoTerminate = 0x0003
	ProtoCookie    = 0x0004
	ProtoSig       = 0x0005

	ProtoPing      = 0x0200
	ProtoPong      = 0x0201
	ProtoBootstrap = 0x0202
	ProtoSearch    = 0x0203
	ProtoRecent    = 0x0204
	ProtoHashList  = 0x0205

	ProtoDhtQuery    = 0x0300
	ProtoDhtAnnounce = 0x0301
)
