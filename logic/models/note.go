package models

//////////////// ESTRUCTURA DE TIPO NOTA PARA GUARDAR LOS DATOS DE NEO4J

type Note struct {
	Id    string `json:"id"`
	Title string `json:"title"`
	//org.title textos orginales en caso de que hagan falta
	//org.text

	Name    string `json:"name"`
	Content string `json:"text"`

	Tags          []string `json:"tags"`          // las etiquetas con las que tiene relación, que son entidades a parte
	Related_notes []string `json:"related_notes"` //Los ID de las notas que tiene relacionadas
	Entities      []string `json:"entities"`      //Las entidades que tiene. Esto son palabras importantes que se encuentran con IA
}

// Simplemente crea una nueva nota
func NewNote(id string, name string, content string, tags []string, related_notes []string, entities []string) Note {

	return Note{
		Id:            id,
		Name:          name,
		Content:       content,
		Tags:          tags,
		Related_notes: related_notes,
		Entities:      entities,
	}
}

//////// ESTRUCTURA DE TIPO TAGS Y ENTIDADES
// USAR CUANDO LA V1 FUNCIONE

type Tag struct {
	Name string
}
