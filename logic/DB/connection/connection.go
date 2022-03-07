package connection

import (
	"github.com/gin-gonic/gin"
	"net/http"
)

/* import (
	"fmt"
	"github.com/JMartinAltostratus/Go-POC-Inditex/logic/models"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io"
	"log"
	"net/http"
	"reflect"
) */

type request struct {
	ID            string   `json:"id"`
	Name          string   `json:"name"`
	Text          string   `json:"content"`
	Relationships []string `json:"relations"` //Aquí está la lista de las tags con las que está relacionada esta nota
}

// ------- CONSTANTES DE LA BD
const (
	dbUser = "neo4j"
	dbPass = "tones-sample-experts"
	dbURI  = "bolt://44.199.246.59:7687"
	dbPort = ":7687"
	dbName = "neo4j"
)

func SearchNeo4J() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req request
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusOK, "searchNeo4J esta funcionando")
		//EJEMPLO, ESTA NOTA HABRÍA QUE GUARDARLA EN OTRO CONTEXTO
		//ctx.Status(http.StatusCreated) //Un 201 si va bien
		//fmt.Print(note.ID())

		// -------- Conexion con la neo4J V1 NO FUNCIONA ERROR 1 VARIABLE BUT DRIVER.NEWSESSION RETURNS 2 VALUES
		//nota := models.NewNote(req.ID, req.Name, req.Text, nil)
		/*driver, err := neo4j.NewDriver(dbURI, neo4j.BasicAuth(dbUser, dbPass, ""))
		defer func() { err = handleClose(driver, err) }() //defer para que se haga al final
		session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeWrite, DatabaseName: dbName})
		defer func() { err = handleClose(session, err) }() //defer para que se haga al final
		createNote(session, models.Note{}) */
		//createRelation(session, models.Relation{})

		// ------- conexion con la Neo4J v2 PRUEBA

	}
}
