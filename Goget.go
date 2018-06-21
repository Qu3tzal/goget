package main

import (
    "errors"
    "net/http"
    "regexp"
    "strings"
    "encoding/json"
    "sync"
    "time"
    "log"
    "strconv"
)

// This is the key-value store.
var (
    valuesMap map[string]string
    valuesMapMutex sync.RWMutex
)

// Interface for the request's JSON body.
type RequestBody struct {
    Value string
}

// Retrieves the ID from the URL.
// Expected format is : /store/{id}
//  id = [a-zA-Z0-9\-_\.]+
func getIdFromURL(url string) (string, error) {
    // Trim the string.
    url = strings.TrimSpace(url)

    // Compile our regex.
    re := regexp.MustCompile("(/)?store/[a-zA-Z0-9\\-_\\.]+")

    // Check we have a correct format.
    match := re.MatchString(url)

    // Check we have a match and no errors.
    if !match {
        return "", errors.New("Could not found the base path. Expected format : /store/{id}")
    }

    // Get the ID.
    lastSlashIndex := strings.LastIndex(url, "/")

    if lastSlashIndex == -1 || lastSlashIndex == len(url) - 1 {
        return "", errors.New("Could not found the id. Expected format : /store/{id}")
    }

    id := url[lastSlashIndex + 1:len(url)]

    if id == "" {
        return "", errors.New("Could not found the id. Expected format : /store/{id}")
    }

    // We return the id.
    return id, nil
}

// Handles a GET request.
func handleGetRequest(writer http.ResponseWriter, request *http.Request) {
    // Get the id.
    id, err := getIdFromURL(request.URL.Path)

    if err != nil {
        // Error on the ID.
        log.Println("[GET] 400 Bad request : \"" + request.URL.Path + "\"")
        writer.WriteHeader(400)
        writer.Write([]byte(err.Error()))
    } else {
        // Lock & unlock the mutex to prevent data races.
        valuesMapMutex.Lock()
        value, found := valuesMap[id]
        valuesMapMutex.Unlock()

        // Check if we have the id in the store.
        if !found {
            log.Println("[GET] 404 ID not found : \"" + id + "\"")

            writer.WriteHeader(404)
            writer.Write([]byte("ID not found : '" + id + "'."))
        } else {
            log.Println("[GET] 200 served id : " + id)

            writer.Header().Set("Content-Type", "application/json")
            writer.WriteHeader(200)
            writer.Write([]byte("{\"" + id + "\": \"" + value + "\"}"))
        }
    }
}

// Handles a POST request.
func handlePostRequest(writer http.ResponseWriter, request *http.Request) {
    // Get the id.
    id, err := getIdFromURL(request.URL.Path)

    // Parses the body.
    var value RequestBody
    jsonErr := json.NewDecoder(request.Body).Decode(&value)

    if err != nil {
        // Error on the ID.
        log.Println("[POST] 400 Bad request : \"" + request.URL.Path + "\"")
        writer.WriteHeader(400)
        writer.Write([]byte(err.Error()))
    } else if jsonErr != nil || value.Value == "" {
        // Error on the value.
        writer.WriteHeader(400)

        if jsonErr != nil {
            log.Println("[POST] 400 error while parsing request body : " + jsonErr.Error())
            writer.Write([]byte("Error while parsing request body : " + jsonErr.Error()))
        } else {
            log.Println("[POST] 400 no value for id : " + id)
            writer.Write([]byte("Value is empty."))
        }
    } else {
        // No error.
        // Lock & unlock the mutex to prevent data races.
        valuesMapMutex.Lock()
        _, found := valuesMap[id]
        valuesMapMutex.Unlock()

        // Check the id is not already used
        if found {
            log.Println("[POST] 409 ID already exists : " + id)

            writer.WriteHeader(409)
            writer.Write([]byte("ID already exists."))
        } else {
            // Write the value in the store.
            // Lock & unlock the mutex to prevent data races.
            valuesMapMutex.Lock()
            valuesMap[id] = value.Value
            valuesMapMutex.Unlock()

            // Answer.
            log.Println("[POST] 200 created value for id : " + id)

            writer.Header().Set("Content-Type", "application/json")
            writer.WriteHeader(200)
        }
    }
}

// Handles a PUT request.
func handlePutRequest(writer http.ResponseWriter, request *http.Request) {
    // Get the id.
    id, err := getIdFromURL(request.URL.Path)

    // Parses the body.
    var value RequestBody
    jsonErr := json.NewDecoder(request.Body).Decode(&value)

    if err != nil {
        // Error on the ID.
        log.Println("[PUT] 400 Bad request : \"" + request.URL.Path + "\"")
        writer.WriteHeader(400)
        writer.Write([]byte(err.Error()))
    } else if jsonErr != nil || value.Value == "" {
        // Error on the value.
        writer.WriteHeader(400)

        if jsonErr != nil {
            log.Println("[PUT] 400 error while parsing request body : " + jsonErr.Error())
            writer.Write([]byte("Error while parsing request body : " + jsonErr.Error()))
        } else {
            log.Println("[PUT] 400 no value for id : " + id)
            writer.Write([]byte("Value is empty."))
        }
    } else {
        // No error.
        // Write the value in the store.
        // Lock & unlock the mutex to prevent data races.
        valuesMapMutex.Lock()
        valuesMap[id] = value.Value
        valuesMapMutex.Unlock()

        // Answer.
        log.Println("[PUT] 200 created or updated value for id : " + id)

        writer.Header().Set("Content-Type", "application/json")
        writer.WriteHeader(200)
    }
}

// Handles a PATCH request.
func handlePatchRequest(writer http.ResponseWriter, request *http.Request) {
    // Get the id.
    id, err := getIdFromURL(request.URL.Path)

    // Parses the body.
    var value RequestBody
    jsonErr := json.NewDecoder(request.Body).Decode(&value)

    if err != nil {
        // Error on the ID.
        log.Println("[PATCH] 400 Bad request : \"" + request.URL.Path + "\"")
        writer.WriteHeader(400)
        writer.Write([]byte(err.Error()))
    } else if jsonErr != nil || value.Value == "" {
        // Error on the value.
        writer.WriteHeader(400)

        if jsonErr != nil {
            log.Println("[PATCH] 400 error while parsing request body : " + jsonErr.Error())
            writer.Write([]byte("Error while parsing request body : " + jsonErr.Error()))
        } else {
            log.Println("[PATCH] 400 no value for id : " + id)
            writer.Write([]byte("Value is empty."))
        }
    } else {
        // No error.
        // Lock & unlock the mutex to prevent data races.
        valuesMapMutex.Lock()
        _, found := valuesMap[id]
        valuesMapMutex.Unlock()

        // Check the id exists
        if !found {
            log.Println("[GET] 404 ID not found : \"" + id + "\"")

            writer.WriteHeader(404)
            writer.Write([]byte("ID not found : '" + id + "'."))
        } else {
            // Update the value in the store.
            // Lock & unlock the mutex to prevent data races.
            valuesMapMutex.Lock()
            valuesMap[id] = value.Value
            valuesMapMutex.Unlock()

            // Answer.
            log.Println("[PATCH] 200 updated value for id : " + id)

            writer.Header().Set("Content-Type", "application/json")
            writer.WriteHeader(200)
        }
    }
}

// Handles a DELETE request.
func handleDeleteRequest(writer http.ResponseWriter, request *http.Request) {
    // Get the id.
    id, err := getIdFromURL(request.URL.Path)

    if err != nil {
        // Error on the ID.
        log.Println("[DELETE] 400 Bad request : \"" + request.URL.Path + "\"")
        writer.WriteHeader(400)
        writer.Write([]byte(err.Error()))
    } else {
        // Lock & unlock the mutex to prevent data races.
        valuesMapMutex.Lock()
        _, found := valuesMap[id]
        valuesMapMutex.Unlock()

        // Check if we have the id in the store.
        if !found {
            log.Println("[DELETE] 404 ID not found : \"" + id + "\"")

            writer.WriteHeader(404)
            writer.Write([]byte("ID not found : '" + id + "'."))
        } else {
            // Delete the value.
            // Lock & unlock the mutex to prevent data races.
            valuesMapMutex.Lock()
            delete(valuesMap, id)
            valuesMapMutex.Unlock()

            log.Println("[DELETE] 200 deleted id : " + id)

            writer.Header().Set("Content-Type", "application/json")
            writer.WriteHeader(200)
        }
    }
}

// The HTTP Handler structure.
type GogetHandler struct{}

// Handles the requests and call the corresponding function depending on the request method.
func (gh GogetHandler) ServeHTTP(writer http.ResponseWriter, request *http.Request) {
    startTime := time.Now()

    switch request.Method {
        case "GET":
            handleGetRequest(writer, request)
        case "POST":
            handlePostRequest(writer, request)
        case "PUT":
            handlePutRequest(writer, request)
        case "PATCH":
            handlePatchRequest(writer, request)
        case "DELETE":
            handleDeleteRequest(writer, request)
    }

    endTime := time.Now()
    log.Println("\tRequest took " + strconv.FormatFloat(endTime.Sub(startTime).Seconds() * 1000.0, 'f', 2, 64) + "ms to run.")
}

// Main function.
func main() {
    // Init the key-value store and its mutex.
    valuesMap = make(map[string]string)
    valuesMapMutex = sync.RWMutex{}

    // Start the web server.
    http.ListenAndServe(":8080", GogetHandler{})
}
