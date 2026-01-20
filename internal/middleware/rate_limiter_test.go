package middleware

import (
	"testing"
	"time"
)

func TestRateLimiter_AllowsRequestsUnderLimit(t *testing.T) {
	limiter := NewRateLimiter(5, time.Minute)

	for i := 0; i < 5; i++ {
		if !limiter.isAllowed("192.168.1.1") {
			t.Errorf("Request %d should be allowed", i+1)
		}
	}
}

func TestRateLimiter_BlocksAfterLimit(t *testing.T) {
	limiter := NewRateLimiter(3, time.Minute)

	// First 3 requests should pass
	for i := 0; i < 3; i++ {
		limiter.isAllowed("192.168.1.1")
	}

	// 4th request should be blocked
	if limiter.isAllowed("192.168.1.1") {
		t.Error("4th request should be blocked")
	}
}

func TestRateLimiter_DifferentIPs(t *testing.T) {
	limiter := NewRateLimiter(2, time.Minute)

	// Max out IP 1
	limiter.isAllowed("192.168.1.1")
	limiter.isAllowed("192.168.1.1")

	// IP 2 should still be allowed
	if !limiter.isAllowed("192.168.1.2") {
		t.Error("Different IP should have its own limit")
	}
}

func TestRateLimiter_ResetsAfterWindow(t *testing.T) {
	// Use a very short window for testing
	limiter := NewRateLimiter(2, 100*time.Millisecond)

	// Max out the limit
	limiter.isAllowed("192.168.1.1")
	limiter.isAllowed("192.168.1.1")

	// Should be blocked now
	if limiter.isAllowed("192.168.1.1") {
		t.Error("Should be blocked at limit")
	}

	// Wait for window to reset
	time.Sleep(150 * time.Millisecond)

	// Should be allowed again
	if !limiter.isAllowed("192.168.1.1") {
		t.Error("Should be allowed after window reset")
	}
}
