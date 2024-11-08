package handlers

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"
)

type RadonRecord struct {
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

func RadonProGaugeHandler(w http.ResponseWriter, r *http.Request) {
	rec, err := queryData(r.Context())
	_ = err
	w.WriteHeader(http.StatusOK)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(rec)
}

// RadonProHandler is a simple handler that returns a JSON response
func RadonProHandler(w http.ResponseWriter, r *http.Request) {
	// "GET /log?AID=&GID=&CPM=10&ACPM=11.81&uSV=0.07&pci=0.45&model=RadonPro HTTP/1.1" 404 455 "-" "-"
	record := RadonRecord{
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

	response := RadonProResponse{Message: "Ok"}
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	json.NewEncoder(w).Encode(response)

}

func insertRecord(ctx context.Context, record RadonRecord) error {
	query := `
        INSERT INTO RadonRecord (aid, gid, acpm, usv, model, pci, cpm)
        VALUES ($1, $2, $3, $4, $5, $6, $7)
    `
	_, err := dbPool.Exec(ctx, query, record.AID, record.GID, record.ACPM, record.USV, record.Model, record.PCI, record.CPM)
	return err
}

func queryData(ctx context.Context) (RadonRecord, error) {
	rec := RadonRecord{}
	// Sample query (replace with your actual query)
	rows, err := dbPool.Query(ctx, "SELECT aid, gid, acpm, usv, model, pci, cpm FROM radonrecord order by id desc limit 1")
	if err != nil {
		return rec, fmt.Errorf("query failed: %w", err)
	}
	defer rows.Close()

	// Iterate through the result set
	for rows.Next() {
		var id int
		var name string
		var createdAt time.Time

		// Scan each row into variables
		if err := rows.Scan(&rec.AID, &rec.GID, &rec.ACPM, &rec.USV, &rec.Model, &rec.PCI, &rec.CPM); err != nil {
			return rec, fmt.Errorf("failed to scan row: %w", err)
		}

		fmt.Printf("ID: %d, Name: %s, Created At: %s\n", id, name, createdAt)
	}

	// Check if any error occurred during iteration
	if rows.Err() != nil {
		return rec, fmt.Errorf("row iteration error: %w", rows.Err())
	}

	return rec, nil
}

/*
CREATE TABLE RadonRecord (
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
