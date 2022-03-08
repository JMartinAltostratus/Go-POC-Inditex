package handlers

import (
	"fmt"
	"github.com/JMartinAltostratus/Go-POC-Inditex/logic/models"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io"
	"log"
	"net/http"
	"reflect"
)

//Este es el modelo de datos de la request que me va a LLEGAR,
//habría que factorizarlo para el objeto tipo NOTA (models note) que voy a guardar
//en la DB, para así poder hacer un parse y convertir las relaciones
//en lo que tiene que ser y etc etc
type request struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Text string `json:"content"`
}

type requestTag struct {
	Tag string `json:"tag"`
}

type requestNote struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type jsonResponse struct {
	title         string
	text          string
	tags          []string
	related_notes []string
	entities      []string
} //Esto se parece PELIGROSAMENTE a un tipo nota. Mirar a ver.

// ------- CONSTANTES DE LA BD
const (
	dbUser = "neo4j"
	dbPass = "tones-sample-experts"
	dbURI  = "bolt://44.199.246.59:7687"
	dbPort = ":7687"
	dbName = "neo4j"
)

// SearchElastic Hacer conexion con Elastic y buscar en funcion de una palabra o whatever
func SearchElastic() gin.HandlerFunc {
	//return func(ctx *gin.Context) {
	//	ctx.String(http.StatusOK, "searchElastic ha funcionao")
	//}
	return func(ctx *gin.Context) {
		//Hacer cosas en funcion del archivo en el que esté
		fmt.Printf("CreateHandler correcto \n")    //LLEGA.
		var req request                            //Me declaro una request con la forma del struct de arriba
		if err := ctx.BindJSON(&req); err != nil { //Aquí se usa gin para gestionar la petición y modifico el objeto anterior
			ctx.JSON(http.StatusBadRequest, err.Error()) //En caso de que no vaya, se devuelve un badrequest 400
			return
		}
	}
}

// Recibe un string en el contexto y devuelve un mapa de ID-Titulo
func SearchByTag() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		query := ""
		var req requestTag
		if err := ctx.BindJSON(&req); err != nil { //Aquí se usa gin para gestionar la petición y modifico el objeto anterior
			ctx.JSON(http.StatusBadRequest, err.Error()) //En caso de que no vaya, se devuelve un badrequest 400
			return
		}
		//query += fmt.Sprintf(`MATCH (n1:New)-[r:HAS_TAG]->(n2) WHERE n2.name = "%s" RETURN r, n1, n2 LIMIT 25`, req.Tag)
		query += fmt.Sprintf(`MATCH (note:Person) WHERE note.name = "%s" RETURN note LIMIT 10`, req.Tag)
		println(query) //Pa probá
		results, err := runQuery(dbURI, dbName, dbUser, dbPass, query)
		if err != nil {
			panic(err)
		}
		if results != nil {
			for _, result := range results {
				fmt.Println(result)
				ctx.String(200, result) //DE VUELTA PAL FRONT
			}
		} else {
			ctx.String(204, "", "No content for this tag")
		}
	}
}
func SearchByNote() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		var req request //Creo la request, que sale de ctx.bindJSON
		if err := ctx.BindJSON(&req); err != nil {
			ctx.JSON(http.StatusBadRequest, err.Error())
			return
		}
		ctx.String(http.StatusOK, "searchByNote esta funcionando")
		query := ""
		query += fmt.Sprintf(`MATCH (note) RETURN (note) AS note`)
		results, err := runQuery(dbURI, dbName, dbUser, dbPass, query)
		if err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Println(result + "1") //TODO probar esto
		}
	}
}

// SearchNeo4J Hacer conexion con Neo4J y llamar con mi objeto tipo nota
// a la función que haya que llamar, en este caso search
func SearchNeo4J() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		// ------- conexion con la Neo4J v2 PRUEBA
		query := ""
		query += fmt.Sprintf(`MATCH (note:Person) RETURN note LIMIT 10`)
		results, err := runQuery(dbURI, dbName, dbUser, dbPass, query)
		if err != nil {
			panic(err)
		}

		for _, result := range results {
			fmt.Println(result)
			ctx.String(200, result) //ESTO DEVUELVE LOS NOMBRES
		}
		//return results
	}
}

//FALTA NADA; CAMBIAR ESTO PARA QUE SE CREE Y DEVUELVA UN OBJETO TIPO
//NOTA Y MANDARLO PAL FRONT CON TREMENDO JSON.MARSHAL. RECORDAR
//QUE SOLO LOS ATRIBUTOS PÚBLICOS SE MARSHALEAN Y LISTO.

func runQuery(uri, database, username, password string, query string) (result []string, err error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	defer func() { err = handleClose(driver, err) }()
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: database})
	defer func() { err = handleClose(session, err) }()
	results, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		result, err := transaction.Run(query, map[string]interface{}{})
		if err != nil {
			return nil, err
		}
		var arr []string
		for result.Next() {
			//Lo que hago con el resultado, en este caso espero
			//que sean string así que los recojo en un array y apaño
			value, found := result.Record().Get("note")
			if found {
				value, ok := value.(neo4j.Node)
				if ok {
					//KETER LAS MOVIDAS EN LA NOTA QUE SEA
					fmt.Println(value.Id, " ---> ID")
					fmt.Println(value.Labels, " ---> LABELS")
					fmt.Println(value.Props, " ---> PROPS")
					var title, _ = value.Props["name"].(string) //FORMA DE COGER UN CAMPO CONCRETO
					//fmt.Println(title, " ---> TITULO")
					arr = append(arr, title)
				}
			}
		}
		if err = result.Err(); err != nil {
			return nil, err
		}
		return arr, nil
	})
	if err != nil {
		return nil, err
	}
	result = results.([]string) //SEGURO QUE ESTO FUNCIONA??
	return result, err
}

func createNote(session neo4j.Session, note models.Note) {
	query := ""
	if note.ID() == "" {
		query += fmt.Sprintf(`CREATE (:Note {idNote: "%s", name: "%s",content: "%s"})})`, note.ID(), note.Name(), note.Content())
	} else {
		print("La nota ya existe")
	}
	r, err := session.Run(query, map[string]interface{}{})
	fmt.Println(r)
	if err != nil {
		log.Fatal(err)
	}
}

func handleClose(closer io.Closer, previousError error) error {
	err := closer.Close()
	if err == nil {
		return previousError
	}
	if previousError == nil {
		return err
	}
	return fmt.Errorf("%v closure error occurred:\n%s\ninitial error was:\n%w", reflect.TypeOf(closer), err.Error(), previousError)
}

// UpdateNote Update a la BD Neo4J con los cambios en la nota
func UpdateNote() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.String(http.StatusOK, "UpdateNote ha funcionao")
	}
}
