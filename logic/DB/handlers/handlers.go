package handlers

import (
	"encoding/json"
	"fmt"
	"github.com/JMartinAltostratus/Go-POC-Inditex/logic/models"
	"github.com/gin-gonic/gin"
	"github.com/neo4j/neo4j-go-driver/v4/neo4j"
	"io"
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

type requestUpdate struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	//org.title textos orginales en caso de que hagan falta
	//org.text

	Name    string `json:"name"`
	Content string `json:"text"`

	Tags          []string `json:"tags"`          // las etiquetas con las que tiene relación, que son entidades a parte
	Related_notes []string `json:"related_notes"` //Los ID de las notas que tiene relacionadas
	Entities      []string `json:"entities"`
}

type requestTag struct {
	Tag string `json:"tag"`
}

type requestNote struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type TempTag struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	Info  string `json:"info"`
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
		results, err := runQueryRetTag(dbURI, dbName, dbUser, dbPass, query)
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
		query := ""
		var req requestNote
		if err := ctx.BindJSON(&req); err != nil { //Aquí se usa gin para gestionar la petición y modifico el objeto anterior
			ctx.JSON(http.StatusBadRequest, err.Error()) //En caso de que no vaya, se devuelve un badrequest 400
			return
		}
		//query += fmt.Sprintf(`MATCH (n1:New)-[r:HAS_TAG]->(n2) WHERE n2.name = "%s" RETURN r, n1, n2 LIMIT 25`, req.Tag)
		query += fmt.Sprintf(`MATCH (note:New) WHERE note.title = "%s" RETURN note LIMIT 10`, req.Name)
		println(query) //Pa probá
		//Recién cambiado, cuidao con lo que hace
		results, err := runQueryRetNote(dbURI, dbName, dbUser, dbPass, query)
		if err != nil {
			panic(err)
		}
		if results != nil {
			for _, result := range results {
				fmt.Println(result)
				ctx.String(200, result) //DE VUELTA PAL FRONT
			}
		} else {
			ctx.String(204, "", "No content for this note")
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
		results, err := runQueryRetTag(dbURI, dbName, dbUser, dbPass, query)
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

func runQueryRetNote(uri, database, username, password string, query string) (result []string, err error) {
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
		var note = models.NewNote("", "", "", nil, nil, nil)
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
					//var title, _ = value.Props["name"].(string) //FORMA DE COGER UN CAMPO CONCRETO
					//fmt.Println(title, " ---> TITULO")
					note.Id, _ = value.Props["id"].(string)
					note.Name, _ = value.Props["name"].(string)
					note.Content, _ = value.Props["text"].(string)
					note.Tags, _ = value.Props["tags"].([]string)
					note.Related_notes, _ = value.Props["related"].([]string)
					note.Entities, _ = value.Props["entities"].([]string)
					//arrprueba := [...]string{"esto", "son", "relaciones entre notas"}
					//note := models.NewNote("1213412", "NotaDePrueba", "Esto es una nota de prueba", nil, nil, nil)
					note = models.NewNote(note.Id, note.Name, note.Content, note.Tags, note.Related_notes, note.Entities)
					var bytes []byte
					bytes, _ = json.Marshal(note)
					println("Nota", note.Name)
					println("JSON????", string(bytes))
					arr = append(arr, string(bytes))
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

func runQueryRetTag(uri, database, username, password string, query string) (result []string, err error) {
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
		var tag TempTag
		for result.Next() {
			//MIRAR A VER LA QUERY PARA VER QUÉ COJO EN ESTE CASO
			value, found := result.Record().Get("note")
			if found {
				value, ok := value.(neo4j.Node)
				if ok {
					tag.Id, _ = value.Props["id"].(string)
					tag.Title, _ = value.Props["name"].(string)
					tag.Info, _ = value.Props["text"].(string)

					//arrprueba := [...]string{"esto", "son", "relaciones entre notas"}
					//note := models.NewNote("1213412", "NotaDePrueba", "Esto es una nota de prueba", nil, nil, nil)
					var bytes []byte
					bytes, _ = json.Marshal(tag)
					println("Nota", tag.Id)
					println("JSON????", string(bytes))
					arr = append(arr, string(bytes))
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

func runQuery(uri, database, username, password string, query string) (result []string, err error) {
	driver, err := neo4j.NewDriver(uri, neo4j.BasicAuth(username, password, ""))
	if err != nil {
		return nil, err
	}
	defer func() { err = handleClose(driver, err) }()
	session := driver.NewSession(neo4j.SessionConfig{AccessMode: neo4j.AccessModeRead, DatabaseName: database})
	defer func() { err = handleClose(session, err) }()
	results, err := session.ReadTransaction(func(transaction neo4j.Transaction) (interface{}, error) {

		result, err := transaction.Run(query, map[string]interface{}{
			//"props": update,
		})
		fmt.Println("Resultado: ", result)
		if err != nil {
			return nil, err
		}
		var arr []string
		for result.Next() {
			value, found := result.Record().Get("note")
			if found {
				value, ok := value.(neo4j.Node)
				if ok {
					var title, _ = value.Props["name"].(string) //FORMA DE COGER UN CAMPO CONCRETO
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

/*func createNote(session neo4j.Session, note models.Note) {
	query := ""
	if note.Id == "" {
		query += fmt.Sprintf(`CREATE (:Note {idNote: "%s", name: "%s",content: "%s"})})`, note.Id, note.Name, note.Content)
	} else {
		print("La nota ya existe")
	}
	r, err := session.Run(query, map[string]interface{}{})
	fmt.Println(r)
	if err != nil {
		log.Fatal(err)
	}
}*/

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
		query := ""
		var req requestUpdate
		if err := ctx.BindJSON(&req); err != nil { //Aquí se usa gin para gestionar la petición y modifico el objeto anterior
			ctx.JSON(http.StatusBadRequest, err.Error()) //En caso de que no vaya, se devuelve un badrequest 400
			return
		}
		//query += fmt.Sprintf(`MATCH (n1:New)-[r:HAS_TAG]->(n2) WHERE n2.name = "%s" RETURN r, n1, n2 LIMIT 25`, req.Tag)
		//ESTA ES LA BUENA ??
		//query += fmt.Sprintf(`MATCH (note:News {id: '%s'}) SET note = {title: '%s', text: '%s' tags: '%s' related: '%s' entities: '%s'} RETURN note.title`, req.Id, req.Title, req.Content, req.Tags, req.Related_notes, req.Entities)
		//ESTA ES LA PRUEBA
		//query += fmt.Sprintf(`MATCH (note:Person {title: '%s'}) SET note = {title: '%s', text: '%s' tags: '%s' related: '%s' entities: '%s'} RETURN note.title`, req.Id, req.Title, req.Content, req.Tags, req.Related_notes, req.Entities)

		query += fmt.Sprintf(`MATCH (note:Person {name: '%s'}) SET note = $props RETURN note.title`, req.Id)
		fmt.Println(req)
		//HAY QUE FORMATEAR LOS ARRAYS AQUÍ PARA INSERTARLOS EN LA DB,
		//EN JSON ME LLEGA COMO [WHATEVER WHATEVER WHATEVER].

		println(query) //Pa probá
		results, err := runQuery(dbURI, dbName, dbUser, dbPass, query)
		if err != nil {
			panic(err)
		}
		if results != nil {
			for _, result := range results {
				fmt.Println("Note: ", result, "updated")
				ctx.String(200, fmt.Sprintf("Note: ", result, "updated")) //DE VUELTA PAL FRONT
			}
		} else {
			ctx.String(204, "", "No matches in the DB for this petition")
		}
	}
}
