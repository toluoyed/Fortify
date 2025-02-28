package main

import (

	"fmt"
	"net/http"
	"os"
	"fortifyApp/upload"
	"fortifyApp/database"
	"fortifyApp/members"
	"fortifyApp/users"
	"fortifyApp/auth"
	"fortifyApp/utils"
	_ "fortifyApp/docs"
	"github.com/swaggo/http-swagger"
	// "log"
)

func main() {

	if _, err := os.Stat("./documentUploads"); os.IsNotExist(err) {
		os.Mkdir("./documentUploads", os.ModePerm)
	}

	db := database.GetConnection()
	
 	defer database.Close(db)

	port := os.Getenv("PORT")
	if port == "" {
        port = "8080" // Default port for local development
    }

	// db.Debug().Exec("CREATE TYPE status AS ENUM ('COMPLETE','INCOMPLETE')")

	// db.Debug().Exec("CREATE TYPE role AS ENUM ('USER','SUPERUSER')")

	// err := db.AutoMigrate(&members.Member{})
	// if err != nil {
	// 	log.Fatalf("Failed to migrate the database: %v", err)
	// }

	// err := db.AutoMigrate(&user.User{})
	// if err != nil {
	// 	log.Fatalf("Failed to migrate the database: %v", err)
	// }


	// Swagger UI
	http.HandleFunc("/swagger/", httpSwagger.WrapHandler)

	http.HandleFunc("/login", func(w http.ResponseWriter, r *http.Request) {
		users.LoginHandler(w,r,db)
	})

	http.HandleFunc("/register", func(w http.ResponseWriter, r *http.Request) {
		users.CreateUserHandler(w,r,db)
	})

	http.Handle("/upload", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		upload.UploadFileHandler(w,r,db)
	})))

	http.Handle("/members", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		switch r.Method {
		case http.MethodPost:
			members.PostMemberHandler(w, r, db)
		case http.MethodGet:
			members.GetMembersHandler(w, r, db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})))

	http.Handle("/members/", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := utils.GetIDFromPath(r.URL.Path)
		if id == "" {
			http.Error(w, "Invalid member ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPost:
			members.UpdateMemberHandler(w,id,r,db)
		case http.MethodDelete:
			members.DeleteMemberHandler(w,id,r,db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		
	})))

	http.Handle("/users", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		
		if r.Method == http.MethodGet {
			users.GetUserHandler(w,r,db)
		}else{
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	
	})))
	
	http.Handle("/users/", auth.AuthMiddleware(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

		id := utils.GetIDFromPath(r.URL.Path)
		if id == "" {
			http.Error(w, "Invalid user ID", http.StatusBadRequest)
			return
		}
		switch r.Method {
		case http.MethodPost:
			users.UpdateUserHandler(w,id,r,db)
		case http.MethodDelete:
			users.DeleteUserHandler(w,id,r,db)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
		
	})))

	// Start the server
	fmt.Println("Starting server on :8080")
	if err := http.ListenAndServe(":"+port , nil); err != nil {
		fmt.Printf("Error starting server: %v", err)
	}
}
