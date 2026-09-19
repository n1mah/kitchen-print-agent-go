package main
import (
	"fmt"
	"net/http"
)
func main() {
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
    	fmt.Fprintf(w, `{"status":"ok"}`)
	})
	http.ListenAndServe("127.0.0.1:9123", nil)
}