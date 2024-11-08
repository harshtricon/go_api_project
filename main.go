package main

import (
    "log"
    "github.com/gin-gonic/gin"
)

func main() {
    r := gin.Default()

    // Serve static files
    r.Static("/static", "./static")

    // Public route for login
    r.POST("/api/login", loginHandler)

    // Protected routes for notes with JWT middleware
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
