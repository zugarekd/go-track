package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
)

type Record struct {
	AID   string
	GID   string
	ACPM  string
	USV   string
	Model string
	PCI   string
	CPM   string
}

type RadonProResponse struct {
	Message string `json:"message"`
}

// RadonProHandler is a simple handler that returns a JSON response
func RadonProHandler(w http.ResponseWriter, r *http.Request) {
	// "GET /log?AID=&GID=&CPM=10&ACPM=11.81&uSV=0.07&pci=0.45&model=RadonPro HTTP/1.1" 404 455 "-" "-"
	aid := r.URL.Query().Get("AID")
	println("AID: ", aid)
	gid := r.URL.Query().Get("GID")
	println("GID: ", gid)
	acpm := r.URL.Query().Get("ACPM")
	println("ACPM: ", acpm)
	usv := r.URL.Query().Get("uSV")
	println("uSV: ", usv)
	model := r.URL.Query().Get("model")
	println("Model: ", model)
	pci := r.URL.Query().Get("pci")
	println("PCI: ", pci)
	cpm := r.URL.Query().Get("CPM")
	println("CPM: ", cpm)

	record := Record{
		AID:   r.URL.Query().Get("AID"),
		GID:   r.URL.Query().Get("GID"),
		ACPM:  r.URL.Query().Get("ACPM"),
		USV:   r.URL.Query().Get("uSV"),
		Model: r.URL.Query().Get("model"),
		PCI:   r.URL.Query().Get("pci"),
		CPM:   r.URL.Query().Get("CPM"),
	}

	// Insert into the database
	if err := insertRecord(r.Context(), record); err != nil {
		http.Error(w, "Failed to insert record", http.StatusInternalServerError)
		log.Printf("Failed to insert record: %v", err)
		return
	}

	fmt.Fprintln(w, "Record inserted successfully")

	// insert data into postgres db using pgx driver
	// db := db.GetDB()
	// _, err := db.Exec("INSERT INTO radonpro (aid, gid, acpm, usv, model, pci, cpm) VALUES ($1, $2, $3, $4, $5, $6, $7)", aid, gid, acpm, usv, model, pci, cpm)
	// if err != nil {
	// 	log.Println(err)
	// 	w.WriteHeader(http.StatusInternalServerError)
	// 	return
	// }

	response := RadonProResponse{Message: "Ok"}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(response)
	w.WriteHeader(http.StatusOK)
}

func insertRecord(ctx context.Context, record Record) error {
	query := `
        INSERT INTO records (aid, gid, acpm, usv, model, pci, cpm)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := dbPool.Exec(ctx, query, record.AID, record.GID, record.ACPM, record.USV, record.Model, record.PCI, record.CPM)
	return err
}

/*
CREATE TABLE records (
    id SERIAL PRIMARY KEY,  -- Optional: Unique identifier for each record
    aid VARCHAR(255) NOT NULL,
    gid VARCHAR(255) NOT NULL,
    acpm VARCHAR(255) NOT NULL,
    usv VARCHAR(255) NOT NULL,
    model VARCHAR(255) NOT NULL,
    pci VARCHAR(255) NOT NULL,
    cpm VARCHAR(255) NOT NULL,
    created_at TIMESTAMP DEFAULT CURRENT_TIMESTAMP  -- Optional: Timestamp for when the record was created
);
*/
