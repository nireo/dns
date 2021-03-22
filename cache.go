package main

import (
	"log"
	"sync"
	"time"

	"golang.org/x/net/dns/dnsmessage"
)

type Cache struct {
	sync.RWMutex
	domains map[string]cacheEntry
}

type cacheEntry struct {
	timestamp, ttl int64
	message        dnsmessage.Message
}

// remember to initialize the cache
var cache *Cache

// Get finds a domain from the cache
func (c *Cache) Get(domain string) (dnsmessage.Message, bool) {
	c.RLock()
	data, ok := c.domains[domain]
	c.RUnlock()

	if time.Now().Unix() > data.timestamp+data.ttl && ok {
		log.Printf("cache expired for domain: %s", domain)
		c.Delete(domain)
		return dnsmessage.Message{}, false
	}

	return data.message, ok
}

// Delete deletes a domain from the cache
func (c *Cache) Delete(domain string) {
	c.Lock()
	delete(c.domains, domain)
	c.Unlock()
}

// Purge removes all domain entries
func (c *Cache) Purge() {
	c.Lock()
	c.domains = make(map[string]cacheEntry)
	c.Unlock()
}

func (c *Cache) Set(domain string, msg dnsmessage.Message) {
	c.Lock()
	c.domains[domain] = cacheEntry{
		message:   msg,
		ttl:       ttl,
		timestamp: time.Now().Unix(),
	}
	c.Unlock()
}
