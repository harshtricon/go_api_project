package main

import (
    "log"
    "github.com/gin-gonic/gin"
    "net/http"
)

func main() {
    r := gin.Default()

    // Serve static files
    r.Static("/static", "./static")

    // API routes with JWT middleware
    r.POST("/api/login", loginHandler)
    authorized := r.Group("/api")
    authorized.Use(AuthMiddleware())
    {
        authorized.GET("/notes", handleGetNotes)
        authorized.POST("/notes", handleCreateNote)
        authorized.DELETE("/notes", handleDeleteNote)
    }

    log.Println("Server started at http://localhost:8000")
    r.Run(":8000")
}
