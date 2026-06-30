package accounts

import (
	"encore.dev/storage/cache"
)

// The global cache feats the accounts service because its the list depended service.

// GlobalCache is a global cache instance that can be used throughout the application.
var GlobalCache = cache.NewCluster("global-cache", cache.ClusterConfig{
	EvictionPolicy: cache.AllKeysLRU,
})
