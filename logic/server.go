package logic

import (
	"fmt"
	handler "github.com/JMartinAltostratus/Go-POC-Inditex/logic/DB/handlers"
	"github.com/gin-gonic/gin"
	"log"
)

type Server struct {
	httpAddr string
	engine   *gin.Engine
}

// New Crea un nuevo servidor con GIN como engine, y parsea la dirección
//httpAddr que luego utilizará el método Run(). Además, registra los
//endpoints de entrada de la API en el registerRoutes
func New(host string, port uint) Server {
	srv := Server{
		engine:   gin.New(),                        //Aquí se usa gin para crear el handler
		httpAddr: fmt.Sprintf("%s:%d", host, port), //Aquí creamos la ADDR a la que lanzar la petición
	}
	srv.registerRoutes() //Al levantar el server le decimos qué endpoints queremos con esta función
	return srv
}

// Run Corre el servidor en la dirección s.httpAddr, definida en el archivo
// bootstrap del paquete cmd/api
func (s *Server) Run() error {
	log.Println("Server running on", s.httpAddr)
	//log.Println(s.engine.Routes())
	return s.engine.Run(s.httpAddr)
}

// Registra las rutas a las que se podrá atacar desde URI:PUERTO/ ENDPOINT, y
// lanza las funciones en handlers. Se trabaja a través del body JSON.
func (s *Server) registerRoutes() {

	// --- RUTAS PARA LAS BÚSQUEDAS EN BD ---
	//s.engine.GET("/searchElastic", handler.SearchElastic())
	s.engine.GET("/searchNeo4J", handler.SearchNeo4J())

	// --- RUTAS PARA LAS EDICIONES EN BD ---
	s.engine.PUT("/editNote", handler.UpdateNote())
}
