package ratelimit

import (
	"fmt"
	"kakigori-api/internal/repository"
	"time"
)

type IPBlocker struct {
	repo          *repository.GlobalIPBlock
	maxRequests   int
	window        time.Duration
	blockDuration time.Duration
}

func NewIPBlocker(maxRequests int, window time.Duration, blockDuration time.Duration) *IPBlocker {
	blocker := &IPBlocker{
		repo:          repository.NewGlobalIPBlock(),
		maxRequests:   maxRequests,
		window:        window,
		blockDuration: blockDuration,
	}
	go blocker.cleanupExpiredBlocks()
	return blocker
}

func (b *IPBlocker) cleanupExpiredBlocks() {
	ticker := time.NewTicker(1 * time.Minute)
	for range ticker.C {
		threshold := time.Now().Add(-b.blockDuration)
		b.repo.CleanupExpiredBlocks(threshold)
	}
}

func (b *IPBlocker) IsBlocked(ip string) (bool, string) {
	blockTime, exists := b.repo.GetBlockTime(ip)
	if exists {
		if time.Since(blockTime) < b.blockDuration {
			remainingTime := b.blockDuration - time.Since(blockTime)
			return true, fmt.Sprintf("Your IP address is blocked. Block will be lifted in %d minutes", int(remainingTime.Minutes()))
		}
		return false, ""
	}
	return false, ""
}

func (b *IPBlocker) IncrementAndCheck(ip string) (bool, error) {
	blockTime, exists := b.repo.GetBlockTime(ip)
	if exists {
		if time.Since(blockTime) < b.blockDuration {
			remainingTime := b.blockDuration - time.Since(blockTime)
			return false, fmt.Errorf("IP address is blocked for %d more minutes", int(remainingTime.Minutes()))
		}
		b.repo.RemoveBlock(ip)
	}

	count := b.repo.IncrementRequestCount(ip)
	if count > b.maxRequests {
		b.repo.AddBlock(ip, time.Now())
		return false, fmt.Errorf("IP address has been blocked for %d hours due to excessive requests", int(b.blockDuration.Hours()))
	}

	go func(ip string) {
		time.Sleep(b.window)
		b.repo.ResetRequestCount(ip)
	}(ip)

	return true, nil
}

func (b *IPBlocker) GetBlockedIPs() map[string]time.Duration {
	blocks := b.repo.GetAllBlocks()
	blockedIPs := make(map[string]time.Duration)
	now := time.Now()
	for ip, blockTime := range blocks {
		if remaining := b.blockDuration - now.Sub(blockTime); remaining > 0 {
			blockedIPs[ip] = remaining
		}
	}
	return blockedIPs
}
