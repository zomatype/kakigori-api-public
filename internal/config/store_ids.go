package config

import (
	"os"
	"strconv"
	"strings"
)

type StoreConfig struct {
	ID        string
	MaxOrders int
}

func GetMaxOrders() int {
	maxOrders := 100
	if maxStr := os.Getenv("KAKIGORI_MAX_ORDERS"); maxStr != "" {
		if max, err := strconv.Atoi(maxStr); err == nil && max > 0 {
			maxOrders = max
		}
	}
	return maxOrders
}

var GetStoreIDs = func() []string {
	storeIDsStr := os.Getenv("KAKIGORI_STORE_IDS")
	if storeIDsStr == "" {
		return []string{}
	}

	storeIDs := strings.Split(storeIDsStr, ",")
	var validIDs []string
	for _, id := range storeIDs {
		if trimmed := strings.TrimSpace(id); trimmed != "" {
			validIDs = append(validIDs, trimmed)
		}
	}

	return validIDs
}

func IsValidStoreID(id string) bool {
	for _, storeID := range GetStoreIDs() {
		if storeID == id {
			return true
		}
	}
	return false
}
