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

func (d DataEntries) Add(wg *sync.WaitGroup, mu *sync.Mutex, key string, value string) bool {
  mu.Lock()
  defer wg.Done()

  if d.Exists(key) {
    return false
  } else {
    entry := DataEntry{ Key: key, Value: value, Created: time.Now(), Updated: time.Now() }
    d[key] = &entry
    mu.Unlock()

    return true
  }
}

func (d DataEntries) Update(wg *sync.WaitGroup, mu *sync.Mutex, key string, value string) bool {
  mu.Lock()
  defer wg.Done()

  if d.Exists(key) {
    d[key].Value = value
    d[key].Updated = time.Now()
    m.Unlock()

    return true
  } else {
    return d.Add(key, value)
  }
}

func (d DataEntries) Delete(wg *sync.WaitGroup, mu *sync.Mutex, key string) bool {
  mu.Lock()
  defer wg.Done()

  if d.Exists(key) {
    delete(d, key)
    mu.Unlock()

    return true
  } else {
    return false
  }
}

func (d DataEntries) Exists(wg *sync.WaitGroup, mu *sync.Mutex, key string) bool {
  mu.Lock()
  defer wg.Done()

  _, ok := d[key]
  mu.Unlock()

  return ok
}

func NewDataEntries() DataEntries {
  return DataEntries{}
}
