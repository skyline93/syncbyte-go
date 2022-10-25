package cache

import (
	"testing"
	"time"
)

func TestCache(t *testing.T) {
	c := New(5, time.Second*60, time.Second)

	v := c.Get("a")
	if v != nil {
		t.Error("cache.Get('a') is error")
	}

	v = c.Get(1)
	if v != nil {
		t.Error("cache.Get(1) is error")
	}

	c.SetDefault("a", 1)
	c.SetDefault(1, "a")

	v = c.Get("a")
	if v == nil || v.(int) != 1 {
		t.Error("cache.Get('a') is error")
	}

	v = c.Get(1)
	if v == nil || v.(string) != "a" {
		t.Error("cache.Get(1) is error")
	}

	k1 := c.SetDefaultWithUuidKey(1)
	k2 := c.SetDefaultWithUuidKey("a")

	v = c.Get(k1)
	if v == nil || v.(int) != 1 {
		t.Error("cache.SetDefaultWithUuidKey(1) is error")
	}

	v = c.Get(k2)
	if v == nil || v.(string) != "a" {
		t.Error("cache.SetDefaultWithUuidKey('a') is error")
	}
}

func TestCacheExpiretion(t *testing.T) {
	c := New(5, time.Second*60, time.Second)
	c.Set("a", 1, 0)
	c.Set("b", 2, time.Second*20)
	c.Set("c", 3, time.Second*30)

	<-time.After(20 * time.Second)

	v := c.Get("a")
	if v == nil {
		t.Error("cache.Get('a') is error")
	}

	v = c.Get("b")
	if v != nil {
		t.Error("cache.Get('b') is error")
	}

	v = c.Get("c")
	if v == nil || v.(int) != 3 {
		t.Error("cache.Get('c') is error")
	}

	<-time.After(10 * time.Second)
	v = c.Get("c")
	if v != nil {
		t.Error("cache.Get('c') is error")
	}

	k1 := c.SetWithUuidKey(1, 10*time.Second)
	k2 := c.SetWithUuidKey("a", 20*time.Second)

	<-time.After(10 * time.Second)
	v = c.Get(k1)
	if v != nil {
		t.Error("cache.SetWithUuidKey(1) is error")
	}

	v = c.Get(k2)
	if v == nil || v.(string) != "a" {
		t.Error("cache.SetWithUuidKey('a') is error")
	}

	<-time.After(10 * time.Second)
	v = c.Get(k2)
	if v != nil {
		t.Error("cache.SetWithUuidKey('a') is error")
	}
}
