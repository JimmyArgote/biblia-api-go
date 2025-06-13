package database

import (
	"database/sql"
	"fmt"
	"log"

	_ "github.com/go-sql-driver/mysql" // Driver do MySQL
)

var DB *sql.DB

// InitDB inicializa a conexão com o banco de dados e a armazena na variável global DB.
func InitDB(connectionString string) {
	var err error
	// Altera o nome do driver para "mysql"
	DB, err = sql.Open("mysql", connectionString)
	if err != nil {
		log.Fatalf("Erro ao abrir a conexão com o banco de dados: %v", err)
	}

	err = DB.Ping()
	if err != nil {
		log.Fatalf("Erro ao conectar com o banco de dados: %v", err)
	}

	fmt.Println("Conexão com o banco de dados MySQL estabelecida com sucesso!")
}
