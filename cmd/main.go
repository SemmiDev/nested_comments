package main

import (
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"log"
	ns "nested-comments"
	"net/http"
)

func main() {
	config, err := ns.LoadConfig()
	if err != nil {
		log.Fatal(err)
	}

	db, err := ns.LoadMysqlConnection(config.MysqlConn)
	if err != nil {
		log.Fatal(err)
	}
	defer ns.CloseMysqlConnection(db)

	postRepo := ns.NewPostRepository(db)
	commentRepo := ns.NewCommentRepository(db)

	server := CreateNewServer(postRepo, commentRepo)
	server.MountHandlers()

	log.Printf("Server started on port %d\n", config.Port)
	log.Fatal(http.ListenAndServe(fmt.Sprintf(":%d", config.Port), server))
}
