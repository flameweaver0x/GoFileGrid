// Simulating batch processing for fetching user permissions
func fetchUserPermissionsBatch(usernames []string) ([]UserPermissions, error) {
    // Assume there's an external API that can handle batch requests for user permissions
    // This function constructs a single request to fetch permissions for multiple users at once
    // Reducing the number of API calls overall
    
    // This is pseudo-code and needs to be adapted to match the API you're interacting with
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
    } else {
        // Handle non-OK responses
    }

    return permissions, nil
}