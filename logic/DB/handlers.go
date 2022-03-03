package DB

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
)

type createRequest struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Text string `json:"content"`
	//Relationships como un array de objetos nota??
}

// SearchElastic Hacer conexion con Elastic y buscar en funcion de una palabra o whatever
func SearchElastic() gin.HandlerFunc {
	//return func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "searchElastic ha funcionao")
	//}
	return func(ctx *gin.Context) {
		//Hacer cosas en funcion del archivo en el que esté
		fmt.Printf("CreateHandler correcto \n")    //LLEGA.
		var req createRequest                      //Me declaro una request con la forma del struct de arriba
		if err := ctx.BindJSON(&req); err != nil { //Aquí se usa gin para gestionar la petición y modifico el objeto anterior
			ctx.JSON(http.StatusBadRequest, err.Error()) //En caso de que no vaya, se devuelve un badrequest 400
			return
		}
	}
}

// SearchNeo4J Hacer conexion con Neo4J y llamar con mi objeto tipo nota
// a la función que haya que llamar, en este caso search
func SearchNeo4J() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "searchNeo4J ha funcionao")

		//EJEMPLO, ESTA NOTA HABRÍA QUE GUARDARLA EN OTRO CONTEXTO
		ctx.Status(http.StatusCreated) //Un 201 si va bien
		//fmt.Print(note.ID())

	}
}

// UpdateNote Update a la BD Neo4J con los cambios en la nota
func UpdateNote() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UpdateNote ha funcionao")
	}
}
