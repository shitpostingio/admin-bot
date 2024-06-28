package cache

import (
	"fmt"
	"strconv"
	"sync"
	"time"

	"github.com/patrickmn/go-cache"
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														STRUCTS														   *
 *																													   *
 ***********************************************************************************************************************
 */

// Action represents an action that must be performed just once.
type Action struct {
	Performed bool
	Mutex     sync.Mutex
}

/*
 ***********************************************************************************************************************
 *																													   *
 *												CONSTS AND VARS														   *
 *																													   *
 ***********************************************************************************************************************
 */

const (
	//deletionCacheExpiration = 5 * time.Minute
	//deletionCacheCleanup    = 10 * time.Minute
	banCacheExpiration = 12 * time.Hour
	banCacheCleanup    = 24 * time.Hour
)

var (
	//deletionCache *cache.Cache
	banCache *cache.Cache
)

/*
 ***********************************************************************************************************************
 *																													   *
 *														START														   *
 *																													   *
 ***********************************************************************************************************************
 */

//CreateActionsCache creates the caches for deletions and bans.
func CreateActionsCache() {
	//go DeletionManager()
	//deletionCache = cache.New(deletionCacheExpiration, deletionCacheCleanup)
	banCache = cache.New(banCacheExpiration, banCacheCleanup)
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													ADDITIONS														   *
 *																													   *
 ***********************************************************************************************************************
 */

//// AddDeletionToCache adds a deletion to its cache.
//func AddDeletionToCache(messageID int) (*Action, error) {
//	var outputAction Action
//	err := deletionCache.Add(strconv.Itoa(messageID), &outputAction, cache.DefaultExpiration)
//	return &outputAction, err
//}

// AddBanToCache adds a ban to its cache.
func AddBanToCache(userID int64) (*Action, error) {
	var outputAction Action
	err := banCache.Add(strconv.FormatInt(userID, 10), &outputAction, cache.DefaultExpiration)
	return &outputAction, err
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													REMOVALS														   *
 *																													   *
 ***********************************************************************************************************************
 */

// RemoveBanFromCache removes a ban from its cache.
func RemoveBanFromCache(userID int64) {
	banCache.Delete(strconv.FormatInt(userID, 10))
}

/*
 ***********************************************************************************************************************
 *																													   *
 *													  GETTERS														   *
 *																													   *
 ***********************************************************************************************************************
 */

//// GetDeletionFromCache gets a deletion from its cache.
//func GetDeletionFromCache(messageID int) (*Action, error) {
//	value, found := deletionCache.Get(strconv.Itoa(messageID))
//	if !found {
//		return nil, fmt.Errorf("deletion request not found")
//	}
//
//	return value.(*Action), nil
//}

// GetBanFromCache gets a ban from its cache.
func GetBanFromCache(userID int64) (*Action, error) {
	value, found := banCache.Get(strconv.FormatInt(userID, 10))
	if !found {
		return nil, fmt.Errorf("ban request not found")
	}

	return value.(*Action), nil
}
