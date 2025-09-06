package ratelimit

import (
	"testing"
	"time"
)

func TestIPBlocker_BlockAndCheck(t *testing.T) {
	blocker := NewIPBlocker(2, 1*time.Minute, 3*time.Hour)
	ip := "192.168.1.1"

	for i := 0; i < 2; i++ {
		ok, err := blocker.IncrementAndCheck(ip)
		if !ok {
			t.Errorf("Request %d should be allowed, got error: %v", i+1, err)
		}
	}

	ok, err := blocker.IncrementAndCheck(ip)
	if ok {
		t.Error("Request should be blocked after exceeding limit")
	}
	if err == nil {
		t.Error("Expected error message for blocked request")
	}

	blocked, message := blocker.IsBlocked(ip)
	if !blocked {
		t.Error("IP should be blocked")
	}
	if message == "" {
		t.Error("Expected block message with remaining time")
	}
}

func TestIPBlocker_BlockExpiration(t *testing.T) {
	blocker := NewIPBlocker(1, 1*time.Second, 1*time.Second)
	ip := "192.168.1.1"

	blocker.IncrementAndCheck(ip)
	blocker.IncrementAndCheck(ip)

	blocked, _ := blocker.IsBlocked(ip)
	if !blocked {
		t.Error("IP should be blocked immediately after exceeding limit")
	}

	time.Sleep(1 * time.Second)
	blocked, _ = blocker.IsBlocked(ip)
	if blocked {
		t.Error("IP should be unblocked after block duration")
	}
}

func TestIPBlocker_GetBlockedIPs(t *testing.T) {
	blocker := NewIPBlocker(1, 1*time.Minute, 3*time.Hour)
	ip1 := "192.168.1.1"
	ip2 := "192.168.1.2"

	blocker.IncrementAndCheck(ip1)
	blocker.IncrementAndCheck(ip1)
	blocker.IncrementAndCheck(ip2)
	blocker.IncrementAndCheck(ip2)

	blockedIPs := blocker.GetBlockedIPs()
	if len(blockedIPs) != 2 {
		t.Errorf("Expected 2 blocked IPs, got %d", len(blockedIPs))
	}

	for ip, remaining := range blockedIPs {
		if remaining > 3*time.Hour {
			t.Errorf("Remaining time for IP %s should not exceed block duration", ip)
		}
	}
}
