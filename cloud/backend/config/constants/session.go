package constants

type key int

const (
	SessionIsAuthorized key = iota
	SessionSkipAuthorization
	SessionID
	SessionIPAddress
	SessionProxies
	SessionUser
	SessionUserCompanyName
	SessionUserRole
	SessionUserID
	SessionUserUUID
	SessionUserTimezone
	SessionUserName
	SessionUserFirstName
	SessionUserLastName
	SessionUserStoreID
	SessionUserStoreName
	SessionUserStoreLevel
	SessionUserStoreTimezone
)
