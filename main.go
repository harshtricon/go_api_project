
package main

import (
    "log"
    "net/http"
)

func main() {
    // Serve static files
    fs := http.FileServer(http.Dir("static"))
    http.Handle("/", fs)

    // API routes
    http.HandleFunc("/api/notes", handleNotes) 

    log.Println("Server started at http://localhost:8000")
    log.Fatal(http.ListenAndServe(":8000", nil))
}
