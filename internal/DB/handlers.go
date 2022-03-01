package DB

import (
	"github.com/JMartinAltostratus/Go-POC-Inditex/internal/models"
	"github.com/gin-gonic/gin"
	"net/http"
)

// SearchElastic Hacer conexion con Elastic y buscar en funcion de una palabra o whatever
func SearchElastic(note models.Note) gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "searchElastic ha funcionao")
	}
}

// SearchNeo4J Hacer conexion con Neo4J y devolver un objeto de tipo Nota
func SearchNeo4J() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "searchNeo4J ha funcionao")
	}
}

// UpdateNote Update a la BD Neo4J con los cambios en la nota
func UpdateNote() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UpdateNote ha funcionao")
	}
}
