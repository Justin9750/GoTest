package utils

import (
	"encoding/json"
	"net/http"
)

func JsonResponse(w http.ResponseWriter, data interface{}) {
	json.NewEncoder(w).Encode(data)
}

func JsonDeserialize(data []byte, obj interface{}) {
	json.Unmarshal(data, &obj)
}
