package lists

import (
	"ShoppingList-Backend/pkg/application"
	"fmt"
	"net/http"
	"time"
)

func ListEvents(app *application.Application) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/event-stream")
		w.Header().Set("Connection", "keep-alive")
		w.Header().Set("Cache-Control", "no-cache")

		flusher, ok := w.(http.Flusher)
		if !ok {
			http.Error(w, "Streaming not supported", 500)
			return
		}

		for i := 0; i < 10; i++ {
			// fmt.Fprintf(w, "event: ping\n")
			fmt.Fprintf(w, "data: {\"i\": %v}", i)
			fmt.Fprintf(w, "\n\n")
			flusher.Flush()
			time.Sleep(1 * time.Second)
		}
	}
}
