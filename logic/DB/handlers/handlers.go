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
	//Relationships como un array de objetos nota??
}

// ------- CONSTANTES DE LA BD
const (
	dbUser = "neo4j"
	dbPass = "posts-preposition-heats"
	dbURI  = "bolt://3.239.185.213:7687"
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

// SearchNeo4J Hacer conexion con Neo4J y llamar con mi objeto tipo nota
// a la función que haya que llamar, en este caso search
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
		results, err := runQuery(dbURI, dbName, dbUser, dbPass)
		if err != nil {
			panic(err)
		}
		for _, result := range results {
			fmt.Println(result)
		}
	}

}

func runQuery(uri, database, username, password string) (result []string, err error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	defer func() { err = handleClose(driver, err) }()
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: database})
	defer func() { err = handleClose(session, err) }()
	results, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {
		println("------- HASTA AQUÍ LLEGA -------")
		//DA UN ERROR CON EL LIMIT AQUÍ PERO PARECE QUE FUNCIONA HASTA AQUÍ
		result, err := transaction.Run(
			`
			MATCH (n)
			RETURN COUNT(n) AS count
			LIMIT $limit
			`, map[string]interface{}{
				"limit": 10,
			})
		if err != nil {
			return nil, err
		}
		var arr []string
		for result.Next() {
			value, found := result.Record().Get("count")
			if found {
				arr = append(arr, value.(string))
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
	result = results.([]string)
	return result, err
}

func createNote(session neo4j.Session, note models.Note) {
	query := ""
	if note.ID() == "" {
		query += fmt.Sprintf(`CREATE (:Note {idNote: "%s", name: "%s",content: "%s"})})`, note.ID(), note.Name(), note.Content())
	} else {
		print("La nota ya existe")
	}
	//map[string]interface{}{}??????????????
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
