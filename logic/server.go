package logic

import (
	"fmt"
	"github.com/JMartinAltostratus/Go-POC-Inditex/logic/models"
	"log"
	"net/http"

	handler "github.com/JMartinAltostratus/Go-POC-Inditex/logic/DB"
	//Aqui faltaría el import de Note si estuviese en un archivo DATA o STRUCTS o algo así
	"github.com/gin-gonic/gin"
)

//////////////// ESTRUCTURA DE TIPO NOTA PARA GUARDAR LOS DATOS DE NEO4J

// Tengo que ver en qué nivel de indentación están las cosas por aquí
// y qué nombres tienen esas request que se están haciendo
type createRequest struct {
	ID   string `json:"id"`
	Name string `json:"name""`
	Text string `json:"content"`
	//Relationships como un array de objetos nota??
}

func CreateHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		//Hacer cosas en funcion del archivo en el que esté
		fmt.Printf("CreateHandler correcto \n")    //LLEGA.
		var req createRequest                      //Me declaro una request con la forma del struct de arriba
		if err := ctx.BindJSON(&req); err != nil { //Aquí se usa gin para gestionar la petición y modifico el objeto anterior
			ctx.JSON(http.StatusBadRequest, err.Error()) //En caso de que no vaya, se devuelve un badrequest 400
			return
		}

		fmt.Printf("------ HA LLEGADO A LA PETICION --------") //NO LLEGA
		//PONER EJEMPLO AQUI PARA VER QUE PASA

		//EJEMPLO, ESTA NOTA HABRÍA QUE GUARDARLA EN OTRO CONTEXTO
		note := models.NewNote(req.ID, req.Name, req.Text, nil) //
		handler.SearchElastic(note)                             //Para probar un ejemplo, hago una busqueda
		ctx.Status(http.StatusCreated)                          //Un 201 si todo va bien
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
	//log.Println(s.engine.Routes())
	return s.engine.Run(s.httpAddr)
}

func (s *Server) registerRoutes() {
	// TODO Crear el archivo searcher con los métodos para Elastic y Neo4J
	// TODO definir los parámetros de las búsquedas

	// Rutas para las búsquedas, desde aquí voy a separar dos archivos, uni
	// para rutas que me haga un create handler desde las busquedas
	// y otro para las ediciones que me cree un handler desde las ediciones
	s.engine.GET("/searchElastic", CreateHandler())
	s.engine.GET("/searchNeo4J", handler.SearchNeo4J())

	// Rutas para editar el contenido de la nota
	s.engine.PUT("/editNote", handler.UpdateNote())
}
