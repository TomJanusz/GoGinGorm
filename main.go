package main

import (
	"goservice-web/config"
	"goservice-web/handlers"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func main() {
	// 1. On récupère la connexion DB (assurez-vous que votre fonction renvoie bien *gorm.DB)
	db := config.ConnectDataBase()

	// OBLIGATOIRE : On passe la DB à la fonction start
	start(db)
}

func start(db *gorm.DB) {
	router := gin.Default()

	// 2. On initialise le handler en lui injectant la DB
	productHandler := &handlers.ProductHandler{DB: db}

	router.GET("/ping", func(c *gin.Context) {
		c.JSON(200, gin.H{"message": "pong"})
	})

	router.LoadHTMLGlob("templates/*")

	router.GET("/", func(c *gin.Context) {
		c.HTML(http.StatusOK, "index.tmpl", gin.H{
			"title": "Main website",
		})
	})

	// 3. On appelle les méthodes sur notre variable 'productHandler' (SANS parenthèses à la fin)
	router.GET("/products/new", productHandler.NewProductForm)
	router.POST("/products", productHandler.CreateProduct)
	router.GET("/products", productHandler.GetProducts)
	router.GET("/products/:id", productHandler.GetProductById)
	router.POST("/products/:id/delete", productHandler.DeleteProduct)
	// Formulaire d'édition (GET)
	router.GET("/products/:id/edit", productHandler.EditProductForm)

	// Traitement de la modification (POST)
	router.POST("/products/:id/update", productHandler.UpdateProduct)

	router.Run()
}
