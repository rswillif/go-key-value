package kvdata

import (
  "time"
  "sync"
)

type DataEntries map[string]*DataEntry

type DataEntry struct {
  Key string
  Value string
  Created time.Time
  Updated time.Time
}

var mu sync.Mutex

func (d DataEntries) Add(key string, value string) bool {
  if d.Exists(key) {
    return false
  } else {
    entry := DataEntry{ Key: key, Value: value, Created: time.Now(), Updated: time.Now() }
    mu.Lock()
    d[key] = &entry
    mu.Unlock()
    return true
  }
}

func (d DataEntries) Update(key string, value string) bool {
  if d.Exists(key) {
    mu.Lock()
    d[key].Value = value
    d[key].Updated = time.Now()
    mu.Unlock()
    return true
  } else {
    return d.Add(key, value)
  }
}

func (d DataEntries) Delete(key string) bool {
  if d.Exists(key) {
    mu.Lock()
    delete(d, key)
    mu.Unlock()
    return true
  } else {
    return false
  }
}

func (d DataEntries) Exists(key string) bool {
  _, ok := d[key]
  return ok
}
