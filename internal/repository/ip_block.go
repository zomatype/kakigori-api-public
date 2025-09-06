package repository

import "time"

type GlobalIPBlock struct {
	requestCounts map[string]int       // ip -> count
	blockedIPs    map[string]time.Time // ip -> blockTime
}

func NewGlobalIPBlock() *GlobalIPBlock {
	return &GlobalIPBlock{
		requestCounts: make(map[string]int),
		blockedIPs:    make(map[string]time.Time),
	}
}

func (r *GlobalIPBlock) IncrementRequestCount(ip string) int {
	r.requestCounts[ip]++
	return r.requestCounts[ip]
}

func (r *GlobalIPBlock) ResetRequestCount(ip string) {
	delete(r.requestCounts, ip)
}

func (r *GlobalIPBlock) AddBlock(ip string, blockTime time.Time) {
	r.blockedIPs[ip] = blockTime
}

func (r *GlobalIPBlock) RemoveBlock(ip string) {
	delete(r.blockedIPs, ip)
}

func (r *GlobalIPBlock) GetBlockTime(ip string) (time.Time, bool) {
	t, ok := r.blockedIPs[ip]
	return t, ok
}

func (r *GlobalIPBlock) GetAllBlocks() map[string]time.Time {
	blocks := make(map[string]time.Time)
	for ip, t := range r.blockedIPs {
		blocks[ip] = t
	}
	return blocks
}

func (r *GlobalIPBlock) CleanupExpiredBlocks(threshold time.Time) {
	for ip, blockTime := range r.blockedIPs {
		if blockTime.Before(threshold) {
			delete(r.blockedIPs, ip)
		}
	}
}
