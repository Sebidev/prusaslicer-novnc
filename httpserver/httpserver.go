package main

import (
	"io/ioutil"
	"net/http"
	"os/exec"
	"strings"

	"github.com/gin-gonic/gin"
)

func main() {
	r := gin.Default()

	destination := "uploads"

	r.POST("/startslice", func(c *gin.Context) {
		name := c.PostForm("name")
		email := c.PostForm("email")
		adresse := c.PostForm("adresse")
		material := c.PostForm("material")
		farbe := c.PostForm("farbe")
		quality := c.PostForm("quality")
		filling := c.PostForm("filling")
		fullpfad := c.PostForm("fullpfad")
		destination := c.PostForm("destination")

		if name == "" || email == "" || adresse == "" || material == "" || farbe == "" || quality == "" || filling == "" || fullpfad == "" || destination == "" {
			c.String(http.StatusBadRequest, "Fehlende Daten in der Anfrage")
			return
		}

		// Starten Sie hier Ihr Backend-Programm als CLI-Befehl.
		cmd := exec.Command("prusa-slicer", fullpfad, "--load", "myconfig.ini", "--export-gcode", "--export-3mf")
		err := cmd.Run()
		if err != nil {
			c.String(http.StatusInternalServerError, "Fehler beim Ausführen des prusa-slicer: %s", err)
			return
		}

		c.String(http.StatusOK, "Vorgang erfolgreich gestartet")
	})

	r.GET("/getData", func(c *gin.Context) {
		files, err := ioutil.ReadDir(destination) // Annahme, dass der Ordner "uploads" im aktuellen Verzeichnis ist.
		if err != nil {
			c.String(http.StatusInternalServerError, "Fehler beim Lesen des Upload-Ordners: %s", err)
			return
		}

		for _, f := range files {
			if strings.HasSuffix(f.Name(), ".gcode") {
				c.String(http.StatusOK, f.Name())
				return
			}
		}

		c.String(http.StatusNotFound, "Keine .gcode Datei gefunden")
	})

	r.Run(":3010")
}
