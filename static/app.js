document.addEventListener("DOMContentLoaded", () => {
    const noteForm = document.getElementById("note-form");
    const noteInput = document.getElementById("note-input");
    const noteList = document.getElementById("note-list");

    async function fetchNotes() {
        const response = await fetch("/api/notes");
        const notes = await response.json();
        noteList.innerHTML = "";
        notes.forEach(note => addNoteToList(note));
    }

    function addNoteToList(note) {
        const listItem = document.createElement("li");
        listItem.textContent = note.text;

        const deleteButton = document.createElement("button");
        deleteButton.textContent = "Delete";
        deleteButton.onclick = async () => {
            await fetch("/api/notes", {
                method: "DELETE",
                headers: { "Content-Type": "application/json" },
                body: JSON.stringify({ id: note.id })
            });
            fetchNotes();
        };

        listItem.appendChild(deleteButton);
        noteList.appendChild(listItem);
    }

    noteForm.addEventListener("submit", async (event) => {
        event.preventDefault();
        const text = noteInput.value;
        await fetch("/api/notes", {
            method: "POST",
            headers: { "Content-Type": "application/json" },
            body: JSON.stringify({ text })
        });
        noteInput.value = "";
        fetchNotes();
    });

    fetchNotes();
});
