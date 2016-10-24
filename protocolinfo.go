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
	ProtoPing      = 0x0006
	ProtoPong      = 0x0007

	ProtoBootstrap = 0x0102 // Request a bootstrap
	ProtoSearch    = 0x0103 // Request a search
	ProtoRecent    = 0x0104 // Request recent posts

	ProtoHashList = 0x0205
	ProtoEntry    = 0x0206 // An individual DHT entry in Content
	ProtoPosts    = 0x0207 // A list of posts in Content

	ProtoDhtQuery    = 0x0300
	ProtoDhtAnnounce = 0x0301
)
