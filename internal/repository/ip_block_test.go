package repository

import (
	"testing"
	"time"
)

func TestGlobalIPBlock(t *testing.T) {
	t.Parallel()

	t.Run("RequestCount", func(t *testing.T) {
		t.Parallel()
		repo := NewGlobalIPBlock()
		ip := "192.168.1.1"
		tcs := []struct {
			name   string
			action func() int
			want   int
		}{
			{"initial count", func() int { return repo.requestCounts[ip] }, 0},
			{"increment", func() int { return repo.IncrementRequestCount(ip) }, 1},
			{"reset", func() int { repo.ResetRequestCount(ip); return repo.requestCounts[ip] }, 0},
		}
		for _, tc := range tcs {
			tc := tc
			t.Run(tc.name, func(t *testing.T) {
				t.Parallel()
				if got := tc.action(); got != tc.want {
					t.Errorf("%s: got %d, want %d", tc.name, got, tc.want)
				}
			})
		}
	})

	t.Run("BlockManagement", func(t *testing.T) {
		t.Parallel()
		repo := NewGlobalIPBlock()
		ip := "192.168.1.1"
		blockTime := time.Now()
		repo.AddBlock(ip, blockTime)
		gotTime, exists := repo.GetBlockTime(ip)
		if !exists {
			t.Error("Expected block to exist")
		}
		if !gotTime.Equal(blockTime) {
			t.Errorf("Block time mismatch: got %v, want %v", gotTime, blockTime)
		}
		repo.RemoveBlock(ip)
		_, exists = repo.GetBlockTime(ip)
		if exists {
			t.Error("Expected block to be removed")
		}
	})

	t.Run("CleanupExpiredBlocks", func(t *testing.T) {
		t.Parallel()
		repo := NewGlobalIPBlock()
		ip := "192.168.1.1"
		blockTime := time.Now().Add(-2 * time.Hour)
		repo.AddBlock(ip, blockTime)
		repo.CleanupExpiredBlocks(time.Now().Add(-1 * time.Hour))
		_, exists := repo.GetBlockTime(ip)
		if exists {
			t.Error("Expected expired block to be cleaned up")
		}
	})
}
