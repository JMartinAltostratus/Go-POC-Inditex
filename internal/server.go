package internal

import (
	"fmt"
	"log"

	handler "github.com/JMartinAltostratus/Go-POC-Inditex/internal/DB"
	//Aqui faltaría el import de Note si estuviese en un archivo DATA o STRUCTS o algo así
	"github.com/gin-gonic/gin"
)

//////////////// ESTRUCTURA DE TIPO NOTA PARA GUARDAR LOS DATOS DE NEO4J

type Note struct {
	id            string
	name          string
	content       string
	relationships []string
}

// Simplemente crea una nueva nota
func NewNote(id string, name string, content string, relationships []string) Note {

	return Note{
		id:            id,
		name:          name,
		content:       content,
		relationships: relationships,
	}
}

/////////// MOVIDAS DEL SERVIDOR PA REFACTORIZAR EN OTRO ARCHIVO
/////////// LOGICA Y ETC

type Server struct {
	httpAddr string
	engine   *gin.Engine
}

//Esto crea un servidor con gin, a esto se llama desde bootstrap
func New(host string, port uint) Server {
	srv := Server{
		engine:   gin.New(),                        //Aquí se usa gin para crear el handler
		httpAddr: fmt.Sprintf("%s:%d", host, port), //Aquí creamos la ADDR a la que lanzar la petición
	}

	srv.registerRoutes() //Al levantar el server le decimos qué endpoints queremos con esta función
	return srv
}

//No sé muy bien qué hace, revisar
func (s *Server) Run() error {
	log.Println("Server running on", s.httpAddr)
	return s.engine.Run(s.httpAddr)
}

func (s *Server) registerRoutes() {
	// TODO Crear el archivo searcher con los métodos para Elastic y Neo4J
	// TODO definir los parámetros de las búsquedas

	// Rutas para las búsquedas desde front
	s.engine.GET("/searchElastic", handler.SearchElastic())
	s.engine.GET("/searchNeo4J", handler.SearchNeo4J())

	// Rutas para editar el contenido de la nota
	s.engine.PUT("/editNote", handler.UpdateNote())
}
