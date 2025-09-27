package repository

import (
	"strconv"

	"github.com/lin-snow/ech0/internal/cache"
)

var echoKeyList []string

const (
	EchoPageCacheKeyPrefix = "echo_page" // echo_page:page:pageSize:search:showPrivate
)

func GetEchoPageCacheKey(page, pageSize int, search string, showPrivate bool) string {
	var showPrivateStr string
	if showPrivate {
		showPrivateStr = "true"
	} else {
		showPrivateStr = "false"
	}
	return EchoPageCacheKeyPrefix + ":" + strconv.Itoa(page) + ":" + strconv.Itoa(pageSize) + ":" + search + ":" + showPrivateStr
}

func ClearEchoPageCache(cache cache.ICache[string, any]) {
	for _, key := range echoKeyList {
		cache.Delete(key)
	}
	echoKeyList = []string{}
}

func GetEchoByIDCacheKey(id uint) string {
	return "echo_id:" + strconv.Itoa(int(id))
}

func GetTodayEchosCacheKey(showPrivate bool) string {
	return "echo_today:" + strconv.FormatBool(showPrivate)
}