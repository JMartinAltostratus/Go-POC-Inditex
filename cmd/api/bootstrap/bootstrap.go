package bootstrap

/*
	Este fichero levanta el servidor y sirve para definir tanto el host
	como el puerto de manera que sea escalable más adelante.
*/

import (
	server "github.com/JMartinAltostratus/Go-POC-Inditex/internal"
	_ "github.com/go-sql-driver/mysql"
)

const (
	host = "localhost"
	port = 8080
)

func Run() error {
	srv := server.New(host, port)
	return srv.Run()
}
