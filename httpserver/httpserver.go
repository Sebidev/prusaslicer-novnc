package main

import (
	"bytes"
	"fmt"
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

type SliceRequest struct {
	Fullpfad    string `json:"fullpfad"`
	Destination string `json:"destination"`
}

// Diese Funktion parst den Dateinamen und gibt die Druckzeit und das Gesamtgewicht zurück.
func parseFileName(name string) (string, string) {
	parts := strings.Split(name, "_")
	if len(parts) < 2 { // Überprüfen, ob wir mindestens zwei Teile haben
		return "", ""
	}
	totalWeight := strings.TrimSuffix(parts[len(parts)-1], ".gcode") // Letzter Teil ohne ".gcode"
	printTime := parts[len(parts)-2]                                 // Vorletzter Teil
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

		if requestData.Fullpfad == "" || requestData.Destination == "" {
			c.String(http.StatusBadRequest, "Fehlende Daten in der Anfrage")
			return
		}

		// Starten Sie hier Ihr Backend-Programm als CLI-Befehl.
		fullpaths := strings.Split(requestData.Fullpfad, ",")

		for _, fullpath := range fullpaths {
			trimmedPath := strings.TrimSpace(fullpath) // Entfernen Sie jeglichen Leerraum um den Pfad herum.
			
			var stderr bytes.Buffer
			cmd := exec.Command("/slic3r/slic3r-dist/prusa-slicer", "/"+trimmedPath, "--load", "/slic3r/myconfig.ini", "--export-gcode", "--export-3mf")
			cmd.Stderr = &stderr
		
			err := cmd.Run()
			if err != nil {
				c.String(http.StatusInternalServerError, "Fehler beim Ausführen des prusa-slicer für Datei %s: %s, Fehlerausgabe: %s", trimmedPath, err, stderr.String())
				return
			}
		}

		files, err := ioutil.ReadDir("/" + requestData.Destination) // Annahme, dass der Ordner "/uploads" im aktuellen Verzeichnis ist.
		if err != nil {
			c.String(http.StatusInternalServerError, "Fehler beim Lesen des Upload-Ordners: %s", err)
			fmt.Printf("Fehler beim Lesen des Upload-Ordners: %s", err)
			return
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".gcode") {
				printTimeStr, weightStr := parseFileName(f.Name())
		
				// Konvertieren Sie die zurückgegebenen Strings in Ints (oder Floats, wenn nötig)
				printTime, err1 := strconv.Atoi(printTimeStr)
				weight, err2 := strconv.Atoi(weightStr)
		
				// Überprüfen Sie auf Konvertierungsfehler
				if err1 != nil || err2 != nil {
					c.String(http.StatusInternalServerError, "Fehler bei der Umwandlung von print_time oder total_weight in Zahlen")
					return
				}
		
				totalPrintTime += printTime
				totalWeight += weight
			}
		}
		
		c.JSON(http.StatusOK, gin.H{
			"total_print_time": totalPrintTime,
			"total_weight": totalWeight,
		})
		
		c.String(http.StatusOK, "Vorgang erfolgreich gestartet")
	})

	r.Run(":3010")
}
