package types

var (
	MsgTypeProxySendCoinsFrom   = "proxy_send_coins_from"
	MsgTypeProxyDisapproveCoins = "proxy_disapprove_coins"
	MsgTypeProxyApproveCoins    = "proxy_approve_coins"

	EventProxyApproveCoins    = "proxy_approve_coins"
	EventProxyDisapproveCoins = "proxy_disapprove_coins"
	EventProxySendCoinsFrom   = "proxy_send_coins_from"

	AttributeKeyProxyAddress           = "proxy_address"
	AttributeKeyProxyOnBehalfOfAddress = "proxy_on_behalf_of_address"
	AttributeKeyProxyToAddress         = "proxy_to_address"
	AttributeKeyProxyDenom             = "proxy_denom"
	AttributeKeyProxyAmount            = "proxy_amount"

	AttributeValueCategory = ModuleName
)
