package aws

import (
	"fmt"
	"io"
	"net/http"
)

// GeneratePresignedURL function fetches a URL from localhost:5000 with a key query parameter and writes it to the response.
func GeneratePresignedURL(w http.ResponseWriter, r *http.Request) {
	// Get the 'key' query parameter from the request
	key := r.URL.Query().Get("key")
	if key == "" {
		// If the key is missing, return an error
		http.Error(w, "Missing 'key' query parameter", http.StatusBadRequest)
		return
	}

	// Construct the URL with the key as a query parameter
	url := fmt.Sprintf("http://localhost:5000/geturl?key=%s", key)

	// Print to indicate that the process is starting
	fmt.Println("Generating presigned URL with key:", key)

	// Send GET request to the URL
	res, err := http.Get(url)
	if err != nil {
		// Handle error and send error message in HTTP response
		http.Error(w, "Failed to fetch URL: "+err.Error(), http.StatusInternalServerError)
		return
	}
	defer res.Body.Close()

	// Read the response body
	body, err := io.ReadAll(res.Body)
	if err != nil {
		// Handle error and send error message in HTTP response
		http.Error(w, "Failed to read response body: "+err.Error(), http.StatusInternalServerError)
		return
	}

	// Write the fetched URL content to the response
	w.Header().Set("Content-Type", "text/plain")
	w.Write(body)
	fmt.Println("Presigned URL content fetched successfully for key:", key)
}
