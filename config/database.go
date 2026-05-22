package config

import (
	"fmt"
	"log"

	"gorm.io/driver/sqlite"
	"gorm.io/gorm"
)

func ConnectDataBase() *gorm.DB {
	dbFile := "local_database.db"
	db, err := gorm.Open(sqlite.Open(dbFile), &gorm.Config{})
	if err != nil {
		log.Fatalf("Échec de la connexion à la base : %v", err)
	}

	// Optionnel mais recommandé : On vérifie que la base répond bien
	sqlDB, err := db.DB()
	if err != nil {
		log.Fatalf("Impossible de récupérer l'instance SQL sous-jacente : %v", err)
	}

	err = sqlDB.Ping()
	if err != nil {
		log.Fatalf("La base de données ne répond pas au Ping : %v", err)
	}

	fmt.Println("Connexion établie avec succès à SQLite !")
	return db
}
