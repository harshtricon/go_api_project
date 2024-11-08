package main

import (
    "database/sql"
    "encoding/json"
    "github.com/gin-gonic/gin"
    "net/http"
    _ "github.com/go-sql-driver/mysql"
)

type Note struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
}

var db *sql.DB

func init() {
    var err error
    db, err = sql.Open("mysql", "user:admin@tcp(127.0.0.1:3306)/mynotes")
    if err != nil {
        panic("Failed to connect to MySQL: " + err.Error())
    }
}
func handleGetNotes(c *gin.Context) {
    rows, err := db.Query("SELECT id, text FROM notes")
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve notes"})
        return
    }
    defer rows.Close()

    var notes []Note
    for rows.Next() {
        var note Note
        if err := rows.Scan(&note.ID, &note.Text); err != nil {
            c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to parse note"})
            return
        }
        notes = append(notes, note)
    }

    c.JSON(http.StatusOK, notes)
}
func handleCreateNote(c *gin.Context) {
    var note Note
    if err := c.BindJSON(&note); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    _, err := db.Exec("INSERT INTO notes (text) VALUES (?)", note.Text)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create note"})
        return
    }

    c.JSON(http.StatusCreated, note)
}

func handleDeleteNote(c *gin.Context) {
    var request struct {
        ID int `json:"id"`
    }
    if err := c.BindJSON(&request); err != nil {
        c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid input"})
        return
    }

    result, err := db.Exec("DELETE FROM notes WHERE id = ?", request.ID)
    if err != nil {
        c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to delete note"})
        return
    }

    rowsAffected, _ := result.RowsAffected()
    if rowsAffected == 0 {
        c.JSON(http.StatusNotFound, gin.H{"error": "Note not found"})
        return
    }

    c.JSON(http.StatusOK, gin.H{"message": "Note deleted successfully"})
}
