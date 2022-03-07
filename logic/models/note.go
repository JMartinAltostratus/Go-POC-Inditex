package models

//////////////// ESTRUCTURA DE TIPO NOTA PARA GUARDAR LOS DATOS DE NEO4J

type Note struct {
	id    string
	title string
	//org.title textos orginales en caso de que hagan falta
	//org.text

	name    string
	content string

	tags          []string // las etiquetas con las que tiene relación, que son entidades a parte
	related_notes []string //Los ID de las notas que tiene relacionadas
	entities      []string //Las entidades que tiene. Esto son palabras importantes que se encuentran con IA
}

// Simplemente crea una nueva nota
func NewNote(id string, name string, content string, tags []string, related_notes []string, entities []string) Note {

	return Note{
		id:            id,
		name:          name,
		content:       content,
		tags:          tags,
		related_notes: related_notes,
		entities:      entities,
	}
}

//Como la estructura está en minúsculas, unos getter setter bien guapos

func (n Note) ID() string {
	return n.id
}

func (n Note) Name() string {
	return n.name
}

func (n Note) Content() string {
	return n.content
}

func (n Note) Relationships() []string {
	return n.related_notes
}

func (n Note) Tags() []string {
	return n.tags
}

//////// ESTRUCTURA DE TIPO TAGS Y ENTIDADES
// USAR CUANDO LA V1 FUNCIONE

type Tag struct {
	Name string
}
