// Producto 1 - FP.050 DevOps y cloud computing (UOC)
// Web en Go que muestra el mensaje "Soy alumno de la UOC" y una imagen.
// Producto 2: se añade la ruta /pod y el entorno (producción o pruebas)
// para comprobar el balanceo de carga en Kubernetes.
package main

import (
	"embed"
	"fmt"
	"html/template"
	"io/fs"
	"log"
	"net/http"
	"os"
	"time"
)

// La plantilla y la imagen se guardan dentro del ejecutable.
//
//go:embed templates static
var ficheros embed.FS

var tpl = template.Must(template.ParseFS(ficheros, "templates/index.html"))

// Datos que se pintan en la plantilla.
type Pagina struct {
	Mensaje  string
	Servidor string
	Hora     string
	Entorno  string
}

// Versión de la aplicación. Se cambia para probar el despliegue automático.
const version = "2"

func getenv(clave, defecto string) string {
	if v := os.Getenv(clave); v != "" {
		return v
	}
	return defecto
}

func main() {
	puerto := getenv("PORT", "8080")
	img, _ := fs.Sub(ficheros, "static")

	// Página principal con el mensaje y la imagen.
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		host, _ := os.Hostname()
		datos := Pagina{
			Mensaje:  "Soy alumno de la UOC",
			Servidor: host,
			Hora:     time.Now().Format("02/01/2006 15:04:05"),
			Entorno:  getenv("ENTORNO", "local"),
		}
		tpl.Execute(w, datos)
	})

	// Imagen de la web.
	http.Handle("/static/", http.StripPrefix("/static/", http.FileServer(http.FS(img))))

	// Comprobación de salud (la usa Kubernetes para saber si el pod está vivo).
	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintln(w, "ok")
	})

	// Producto 2: indica qué pod ha respondido. Sirve para ver el balanceo.
	http.HandleFunc("/pod", func(w http.ResponseWriter, r *http.Request) {
		host, _ := os.Hostname()
		fmt.Fprintf(w, "Pod: %s | Entorno: %s | Version: %s\n",
			host, getenv("ENTORNO", "local"), version)
	})

	log.Printf("Servidor escuchando en el puerto %s", puerto)
	log.Fatal(http.ListenAndServe(":"+puerto, nil))
}
