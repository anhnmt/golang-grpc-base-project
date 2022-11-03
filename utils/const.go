package utils

const (
	// RedisKeyAuthSession is the redis key of the auth session.
	RedisKeyAuthSession = "auth:%s:session:%s"
	// RedisKeyListAuthPermissions is the redis key of the list of auth permissions.
	RedisKeyListAuthPermissions = "auth:permissions"
)
