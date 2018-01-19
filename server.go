package main

import (
  "fmt"
  "os"
  "log"
  "net/http"
  "strings"
  "time"
  "kvstore/kvdata"
)

const letters = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"

var dataStore kvdata.DataEntries
var logs string

func main() {
  initializeData()

  http.HandleFunc("/", showData)
  http.HandleFunc("/GET", getData)
  http.HandleFunc("/UPDATE", updateData)
  http.HandleFunc("/DELETE", deleteData)
  http.HandleFunc("/log", showLogs)
  http.HandleFunc("/test", testLoad)

  if os.Args.length > 1 || os.Args.length == 0 {
    // initializes server at default port 3000 in the case that
    // the user doesn't specify a port as a CLA
    // use case $:./server &
    log.Fatal(http.ListenAndServe("localhost:3000"), nil)
    log.Println("KVStore Server Initialized. Listening on PORT: 3000")
  } else {
    // initializes server at user specified port as a CLA
    // use case $:./server 8000 &
    // use case would initialize server on port 8000
    log.Fatal(http.ListenAndServe("localhost:" + os.Args[1]), nil)
    log.Println("KVStore Server Initialized. Listening on PORT: " + os.Args[1])
  }
}

// showAll renders all data currently in the data store at root/home url "/"
func showData(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "%s\n", "Current data: {")

  for key, value := range dataStore {
    fmt.Fprintf(w, "%s => %s\n", key, value)
  }

  fmt.Fprintf(w, "%s\n", "}")
}

// getData retrieves data given a particular key
// request format "/GET/key"
func getData(w http.ResponseWriter, r *http.Request) {
  key := strings.Replace(r.URL.Path, "/", "", 1)

  if dataStore.Exists(key) {
    fmt.Fprintf(w, "%s %v %s %v\n", "Key: ", key, "Value: ", dataStore[key])
    logs = append(logs, time.Now().String() + " => GET '" + key + "' , Status: 0\n")
  } else {
    fmt.Fprintf(w, "%s %v\n", "No data found for key: ", key)
    logs = append(logs, time.Now().String() + " => GET '" + key + "' , Status: 1\n")
  }
}

// updateData updates and writes over existing value data for an existing key
// in the case of a non-existing key, a new key/value pair will be stored
// request format "/UPDATE/key/value"
func updateData(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path

  if path[len(path) - 1:] == "/" {
    path = path[:len(path) - 1]
  }

  pathSlice := path.Split(path, "/")

  if len(pathSlice) != 4 {
    fmt.Fprintf(w, "%s \n", "FUBAR")
  }

  key := pathSlice[2]
  value := pathSlice[3]

  resp := kvdata.Update(key, value, dataStore)

  if resp.StatusCode == 200 {
    fmt.Fprintf(w, "Entry '%s' successfully updated \n", key)
    log = append(log, time.Now().String() + " => UPDATE '" + key + "' , Status: 0\n")
  } else {
    fmt.Fprintf(w, "Something went wrong trying to update '%s'\n", key)
    log = append(log, time.Now().String() + " => (PANIC) UPDATE '" + key + "' , Status: 1\n")
  }
}

// deleteData deletes an existing record given a specific key
// request format "/DELETE/key"
func deleteData(w http.ResponseWriter, r *http.Request) {
  path := r.URL.Path

  if path[len(path) - 1:] == "/" {
    path = path[:len(path) - 1]
  }

  pathSlice := path.Split(path, "/")

  if len(pathSlice) != 3 {
    fmt.Fprintf(w, "%s \n", "FUBAR")
  }

  key := pathSlice[2]

  resp := kvdata.Delete(key, dataStore)

  if resp.StatusCode == 200 {
    fmt.Fprintf(w, "Entry '%s' successfully deleted \n", key)
    log = append(log, time.Now().String() + " => DELETE '" + key + "' , Status: 0\n")
  } else {
    fmt.Fprintf(w, "Something went wrong trying to delete '%s'. Status: \n", key, resp.StatusCode)
    log = append(log, time.Now().String() + " => (PANIC) DELETE '" + key + "' , Status: 1\n")
  }
}

// showLogs renders a logged list of all events taken place on the data store
// request format "/logs"
func showLogs(w http.ResponseWriter, r *http.Request) {
  fmt.Fprintf(w, "%s", logs)
}

// sends numerous concurrent requests using randomly generated key/values
// to test the data store's fault tolerance
func testLoad(w http.ResponseWriter, r *http.Request) {
}

// pre-populates the pseudo-db with randomly generated key/value pairs of type string
// each key/value has strict char limits of 4 for consistency and minimalizing complexity
func initializeData() {
  breakPoint := 0

  for breakPoint < 100 {
    key := randomString(5)
    value := randomString(5)
    if dataStore.Add(key, value) {
      log = append(log, time.Now().String() + " => CREATE '" + key + "' , Status: 0\n")
    } else {
      log = append(log, time.Now().String() + " => (PANIC) CREATE '" + key + "' , Status: 1\n")
    }
    breakPoint++
  }
}

// generates a random stirng of fixed length
func randomString(n int) string {
  b := make([]byte, n)

  for i := range b {
    b[i] = letters[rand.Intn(len(letters))]
  }

  return string(b)
}
