package main

import (
    "bytes"
    "encoding/json"
    "fmt"
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
    if cachedPermissions, found := permissionsCache[cacheKey]; found {
        cacheLock.RUnlock()
        return cachedPermissions, nil
    }
    cacheLock.RUnlock()

    var permissions []UserPermissions
    requestBody, err := json.Marshal(usernames)
    if err != nil {
        return nil, err
    }

    request, err := http.NewRequest("POST", "https://example.com/api/permissions/batch", bytes.NewBuffer(requestBody))
    if err != nil {
        return nil, err
    }
    request.Header.Set("Content-Type", "application/json")

    client := &http.Client{}
    response, err := client.Do(request)
    if err != nil {
        return nil, err
    }
    defer response.Body.Close()

    if response.StatusCode == http.StatusOK {
        err = json.NewDecoder(response.Body).Decode(&permissions)
        if err != nil {
            return nil, err
        }
        cacheLock.Lock()
        permissionsCache[cacheKey] = permissions
        cacheLock.Unlock()
    } else {
        return nil, fmt.Errorf("Bad response: %d %s", response.StatusCode, http.StatusText(response.StatusCode))
    }

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