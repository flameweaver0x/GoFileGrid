package main

import (
    "bytes"
    "encoding/json"
    "fmt"
    "io/ioutil"
    "net/http"
    "sync"
)

type UserPermissions struct {
    Username    string
    Permissions []string
}

var (
    permissionsCache = make(map[string][]UserPermissions)
    cacheLock        = sync.RWMutex{}
)

func generateCacheKey(usernames []string) string {
    return fmt.Sprintf("%v", usernames)
}

func fetchUserPermissionsBatch(usernames []string) ([]UserPermissions, error) {
    cacheKey := generateCacheKey(usernames)

    cacheLock.RLock()
    cachedPermissions, found := permissionsCache[cacheKey]
    cacheLock.RUnlock()
    if found {
        return cachedPermissions, nil
    }

    // Prepare the request body for the HTTP POST request
    requestBody, err := json.Marshal(usernames)
    if err != nil {
        return nil, fmt.Errorf("error marshaling usernames: %v", err)
    }

    request, err := http.NewRequest("POST", "https://example.com/api/permissions/batch", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, fmt.Errorf("error creating new HTTP request: %v", err)
    }
    request.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, fmt.Errorf("error executing HTTP request: %v", err)
    }
    defer response.Body.Close()

    if response.StatusCode != http.StatusOK {
        // Reading the response body for more detailed error context
        responseBody, _ := ioutil.ReadAll(response.Body) // Ignoring error on purpose to not override the main error
        return nil, fmt.Errorf("bad response: %d %s, details: %s", response.StatusCode, http.StatusText(response.StatusCode), responseBody)
    }

    var permissions []UserPermissions
    err = json.NewDecoder(response.Body).Decode(&permissions)
    if err != nil {
        return nil, fmt.Errorf("error decoding response JSON: %v", err)
    }

    // Updating the cache
    cacheLock.Lock()
    permissionsCache[cacheKey] = permissions
    cacheLock.Unlock()

    return permissions, nil
}

func main() {
    users := []string{"user1", "user2"}
    permissions, err := fetchUserPermissionsBatch(users)
    if err != nil {
        fmt.Println("Error fetching permissions:", err)
        return
    }

    fmt.Println("Permissions:", permissions)
}