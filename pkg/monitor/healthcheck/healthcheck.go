package healthcheck

import (
	"encoding/json"
	"github.com/dalmarcogd/gwp/internal"
	"net/http"
)

//Handler that return a health check response
func Handler(writer http.ResponseWriter, request *http.Request) {
	if http.MethodGet != request.Method {
		writer.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	writer.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(writer).Encode(map[string]interface{}{
		"status": internal.GetServerRun().Healthy(),
	})
}
