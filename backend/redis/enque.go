package redis
import (
	"encoding/json"
	
	"net/http"
	
)

type EnqueueRequest struct {
	Queue string `json:"queue"`
	Item  string `json:"item"`
}
func EnqueueHandler(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		return
	}

	var req EnqueueRequest
	err := json.NewDecoder(r.Body).Decode(&req)
	if err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	if req.Queue == "" || req.Item == "" {
		http.Error(w, "Queue and Item are required", http.StatusBadRequest)
		return
	}

	err = Enqueue(req.Queue, req.Item)
	if err != nil {
		http.Error(w, "Failed to enqueue item: "+err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(map[string]string{"message": "Successfully enqueued item"})
}