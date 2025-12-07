package web

import (
	"net/http"

	"github.com/KurosawaAngel/kindler/internal/application/interactors"
)

type Handler struct {
	sendFile *interactors.SendFile
}

func NewHandler(sendFile *interactors.SendFile) *Handler {
	return &Handler{
		sendFile: sendFile,
	}
}

func (h *Handler) SendFile(w http.ResponseWriter, r *http.Request) {
	toEmail := r.FormValue("email")

	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "File not found", http.StatusBadRequest)
		return
	}
	defer file.Close()

	if err := h.sendFile.Execute(file, header.Filename, toEmail); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	http.Redirect(w, r, "/", http.StatusSeeOther)
}
