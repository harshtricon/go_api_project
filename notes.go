// notes.go
package main

import (
    "encoding/json"
    "net/http"
    "os"
)

type Note struct {
    ID   int    `json:"id"`
    Text string `json:"text"`
}

var notes []Note
var nextID int

func loadNotes() error {
    file, err := os.Open("storage.json")
    if err != nil {
        return err
    }
    defer file.Close()

    decoder := json.NewDecoder(file)
    return decoder.Decode(&notes)
}

func saveNotes() error {
    file, err := os.Create("storage.json")
    if err != nil {
        return err
    }
    defer file.Close()

    encoder := json.NewEncoder(file)
    return encoder.Encode(notes)
}

func handleNotes(w http.ResponseWriter, r *http.Request) {
    switch r.Method {
    case http.MethodGet:
        json.NewEncoder(w).Encode(notes)

    case http.MethodPost:
        var note Note
        json.NewDecoder(r.Body).Decode(&note)

        note.ID = nextID
        nextID++
        notes = append(notes, note)
        saveNotes()
        json.NewEncoder(w).Encode(note)

    case http.MethodDelete:
        var requestData struct {
            ID int `json:"id"`
        }
        json.NewDecoder(r.Body).Decode(&requestData)

        for i, note := range notes {
            if note.ID == requestData.ID {
                notes = append(notes[:i], notes[i+1:]...)
                saveNotes()
                w.WriteHeader(http.StatusOK)
                return
            }
        }
        http.Error(w, "Note not found", http.StatusNotFound)
    }
}
