package main

import (
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"
	"bytes"
	"github.com/gin-gonic/gin"
)

type SliceRequest struct {
	Material   string `json:"material"`
	Farbe      string `json:"farbe"`
	Quality    string `json:"quality"`
	Filling    string `json:"filling"`
	Fullpfad   string `json:"fullpfad"`
	Destination string `json:"destination"`
}

// Diese Funktion parst den Dateinamen und gibt die Druckzeit und das Gesamtgewicht zurück.
func parseFileName(name string) (string, string) {
	parts := strings.Split(name, "_")
	if len(parts) < 2 { // Überprüfen, ob wir mindestens zwei Teile haben
		return "", ""
	}
	totalWeight := strings.TrimSuffix(parts[len(parts)-1], ".gcode") // Letzter Teil ohne ".gcode"
	printTime := parts[len(parts)-2] // Vorletzter Teil
	return printTime, totalWeight
}

func main() {
	r := gin.Default()

	r.POST("/startslice", func(c *gin.Context) {
		var requestData SliceRequest
		if err := c.BindJSON(&requestData); err != nil {
			c.String(http.StatusBadRequest, "Fehler beim Parsen von JSON: %s", err)
			return
		}
	
		if requestData.Material == "" || requestData.Farbe == "" || requestData.Quality == "" || requestData.Filling == "" || requestData.Fullpfad == "" || requestData.Destination == "" {
			c.String(http.StatusBadRequest, "Fehlende Daten in der Anfrage")
			return
		}

		// Starten Sie hier Ihr Backend-Programm als CLI-Befehl.
		var stderr bytes.Buffer
		cmd := exec.Command("/slic3r/slic3r-dist/prusa-slicer", "/"+requestData.Fullpfad, "--load", "/slic3r/myconfig.ini", "--export-gcode", "--export-3mf")
		cmd.Stderr = &stderr
		
		err := cmd.Run()
		if err != nil {
			c.String(http.StatusInternalServerError, "Fehler beim Ausführen des prusa-slicer: %s, Fehlerausgabe: %s", err, stderr.String())
			return
		}

		files, err := ioutil.ReadDir("/"+requestData.Destination) // Annahme, dass der Ordner "uploads" im aktuellen Verzeichnis ist.
		if err != nil {
			c.String(http.StatusInternalServerError, "Fehler beim Lesen des Upload-Ordners: %s", err)
			fmt.Printf("Fehler beim Lesen des Upload-Ordners: %s", err)
			return
		}

	    /* for _, f := range files {
			if strings.HasSuffix(f.Name(), ".gcode") {
				c.String(http.StatusOK, f.Name())
				fmt.Printf(f.Name())
				return
			}
		} */

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".gcode") {
				printTime, totalWeight := parseFileName(f.Name())
				c.JSON(http.StatusOK, gin.H{
					"print_time":   printTime,
					"total_weight": totalWeight,
				})
				return
			}
		}

		c.String(http.StatusOK, "Vorgang erfolgreich gestartet")
	})

	r.Run(":3010")
}
