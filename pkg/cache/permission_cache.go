package cache

import "sync"

var permissionCache = make(map[uint]map[string]bool)
var permissionCacheMutex = &sync.RWMutex{}

func GetUserPermissionCache(userID uint) (map[string]bool, bool) {
	permissionCacheMutex.RLock()
	perms, exists := permissionCache[userID]
	permissionCacheMutex.RUnlock()
	return perms, exists
}

func SetUserPermissionCache(userID uint, permissions map[string]bool) {
	permissionCacheMutex.Lock()
	permissionCache[userID] = permissions
	permissionCacheMutex.Unlock()
}

func ClearUserPermissionCache(userID uint) {
	permissionCacheMutex.Lock()
	delete(permissionCache, userID)
	permissionCacheMutex.Unlock()
}
