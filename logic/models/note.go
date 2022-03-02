package models

//////////////// ESTRUCTURA DE TIPO NOTA PARA GUARDAR LOS DATOS DE NEO4J

type Note struct {
	id            string
	name          string
	content       string
	relationships []string
}

// Simplemente crea una nueva nota
func NewNote(id string, name string, content string, relationships []string) Note {

	return Note{
		id:            id,
		name:          name,
		content:       content,
		relationships: relationships,
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

//Como las relationships vendrán como objetos aún quedan cosas por mirar
func (n Note) Relationships() []string {
	return n.relationships
}
