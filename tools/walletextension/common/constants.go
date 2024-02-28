package common

import (
	"time"
)

const (
	Localhost = "127.0.0.1"

	JSONKeyAddress      = "address"
	JSONKeyData         = "data"
	JSONKeyErr          = "error"
	JSONKeyFrom         = "from"
	JSONKeyID           = "id"
	JSONKeyMethod       = "method"
	JSONKeyParams       = "params"
	JSONKeyResult       = "result"
	JSONKeyRoot         = "root"
	JSONKeyRPCVersion   = "jsonrpc"
	JSONKeySignature    = "signature"
	JSONKeySubscription = "subscription"
	JSONKeyCode         = "code"
	JSONKeyMessage      = "message"
)

const (
	PathRoot                            = "/"
	PathReady                           = "/ready/"
	PathViewingKeys                     = "/viewingkeys/"
	PathGenerateViewingKey              = "/generateviewingkey/"
	PathSubmitViewingKey                = "/submitviewingkey/"
	PathJoin                            = "/join/"
	PathAuthenticate                    = "/authenticate/"
	PathAuthenticatePersonal            = "/authenticate_personal/"
	PathQuery                           = "/query/"
	PathRevoke                          = "/revoke/"
	PathObscuroGateway                  = "/"
	PathHealth                          = "/health/"
	PathNetworkHealth                   = "/network-health/"
	WSProtocol                          = "ws://"
	HTTPProtocol                        = "http://"
	DefaultUser                         = "defaultUser"
	UserQueryParameter                  = "u"
	EncryptedTokenQueryParameter        = "token"
	AddressQueryParameter               = "a"
	MessageUserIDLen                    = 40
	EthereumAddressLen                  = 42
	GetStorageAtUserIDRequestMethodName = "getUserID"
	SuccessMsg                          = "success"
	APIVersion1                         = "/v1"
	MethodEthSubscription               = "eth_subscription"
	PathVersion                         = "/version/"
	DeduplicationBufferSize             = 20
)

var ReaderHeadTimeout = 10 * time.Second
