package repository

// var todoKeyList []string

// const (
// 	TodoByUserCacheKeyPrefix = "todo_user" // todo_user:userid
// 	TodoByIDCacheKeyPrefix   = "todo_id"   // todo_id:todoid
// )

// func GetTodoByUserCacheKey(id uint) string {
// 	return TodoByUserCacheKeyPrefix + ":" + strconv.Itoa(int(id))
// }

// func GetTodoByIDCacheKey(id uint) string {
// 	return TodoByIDCacheKeyPrefix + ":" + strconv.Itoa(int(id))
// }

// func ClearAllTodoCache(cache cache.ICache[string, any]) {
// 	for _, key := range todoKeyList {
// 		cache.Delete(key)
// 	}
// 	todoKeyList = []string{}
// }