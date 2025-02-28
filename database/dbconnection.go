package database

import(
	"database/sql"
	"os"
	"gorm.io/driver/postgres"
  	"gorm.io/gorm"
 	"time"
	"fmt"
	"log"
	_ "github.com/lib/pq"
)


func GetConnection() *gorm.DB {

	// host := "localhost"
	// port := "5432"
	// user := "fortifyuser"
	// password := "mypassword"
	// dbname := "fortifydb"
	dburl := os.Getenv("DATABASE_URL")

	// Format connection string
	// psqlInfo := fmt.Sprintf("host=%s port=%s user=%s password=%s dbname=%s sslmode=disable",
	// 	host, port, user, password, dbname)
	
	// fmt.Printf("Psql: %v", psqlInfo)
	
	// db, err := sql.Open("postgres", psqlInfo)

	db, err := sql.Open("postgres", dburl)

	if err != nil {
		panic(err)
	}

	db.SetMaxIdleConns(10)
	db.SetMaxOpenConns(100)
	db.SetConnMaxIdleTime(5 * time.Minute)
	db.SetConnMaxLifetime(60 * time.Minute)

	fmt.Printf("Type DB: %T\n", db)

	gormDB, err := gorm.Open(postgres.New(postgres.Config{Conn: db,}), &gorm.Config{})

	fmt.Printf("Gorm DB: %T\n", gormDB)
	  
	return gormDB

}

func Close(db *gorm.DB){
	
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("failed to get sql.DB: %v", err)
	}

	// Close the database connection when no longer needed
	sqlDB.Close()
}