package handlers

import (
	"goservice-web/models"
	"net/http"

	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

type ProductHandler struct {
	DB *gorm.DB
}

// Affiche la page du formulaire (GET /products/new)
func (h *ProductHandler) NewProductForm(c *gin.Context) {
	c.HTML(http.StatusOK, "create.tmpl", gin.H{})
}

// Traite la création (POST /products)
func (h *ProductHandler) CreateProduct(c *gin.Context) {
	var newProduct models.Product

	// ShouldBind remplit automatiquement le Name et Statut depuis le formulaire HTML ou le JSON
	if err := c.ShouldBind(&newProduct); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// Sauvegarde dans SQLite via GORM
	result := h.DB.Create(&newProduct)
	if result.Error != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de créer le produit"})
		return
	}

	// Si l'utilisateur a soumis le formulaire depuis son navigateur, on le redirige
	if c.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		c.Redirect(http.StatusSeeOther, "/products")
		return
	}

	// Sinon (Postman / API), on répond en JSON
	c.JSON(http.StatusCreated, newProduct)
}

// GetProducts récupère les produits de la DB et les affiche dans le template HTML
func (h *ProductHandler) GetProducts(c *gin.Context) {
	var products []models.Product

	// GORM va chercher tous les produits dans SQLite
	if err := h.DB.Find(&products).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Erreur lors de la récupération"})
		return
	}

	// On affiche le template "products.tmpl"
	c.HTML(http.StatusOK, "products.tmpl", gin.H{
		"products": products, // On passe la liste sous la clé "products" pour le HTML
	})
}

func (h *ProductHandler) GetProductById(c *gin.Context) {
	// 1. On récupère l'ID écrit dans l'URL
	id := c.Param("id")

	var product models.Product

	// 2. GORM cherche le premier produit qui correspond à cet ID
	// .First() renvoie une erreur si l'ID n'existe pas en base de données
	if err := h.DB.First(&product, id).Error; err != nil {
		// Si le produit n'existe pas, on renvoie une erreur 404
		c.HTML(http.StatusNotFound, "error.tmpl", gin.H{
			"message": "Désolé, ce produit n'existe pas.",
		})
		return
	}

	// 3. On affiche le produit dans un template dédié (ou en JSON si vous préférez)
	c.HTML(http.StatusOK, "product.tmpl", gin.H{
		"product": product,
	})
}

// DeleteProduct gère la suppression d'un produit (POST /products/:id/delete)
func (h *ProductHandler) DeleteProduct(c *gin.Context) {
	// 1. On récupère l'ID du produit à supprimer depuis l'URL
	id := c.Param("id")

	var product models.Product

	// 2. On demande à GORM de supprimer le produit qui a cet ID
	// GORM cible automatiquement la bonne table grâce au type de la variable 'product'
	if err := h.DB.Delete(&product, id).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de supprimer le produit"})
		return
	}

	// 3. Une fois supprimé, on redirige instantanément l'utilisateur vers la liste rafraîchie
	c.Redirect(http.StatusSeeOther, "/products")
}

// EditProductForm affiche le formulaire de modification pré-rempli (GET /products/:id/edit)
func (h *ProductHandler) EditProductForm(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// On cherche le produit pour pouvoir pré-remplir le formulaire HTML
	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit introuvable"})
		return
	}

	c.HTML(http.StatusOK, "update.tmpl", gin.H{
		"product": product,
	})
}

// UpdateProduct enregistre les modifications en base (POST /products/:id/update)
func (h *ProductHandler) UpdateProduct(c *gin.Context) {
	id := c.Param("id")
	var product models.Product

	// 1. On vérifie d'abord que le produit existe bien en base
	if err := h.DB.First(&product, id).Error; err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "Produit introuvable"})
		return
	}

	// 2. On récupère les nouvelles données envoyées par le formulaire HTML
	if err := c.ShouldBind(&product); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Données invalides"})
		return
	}

	// 3. On demande à GORM de sauvegarder toutes les modifications en base
	// La méthode .Save() met à jour la ligne entière correspondant à l'ID de 'product'
	if err := h.DB.Save(&product).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Impossible de mettre à jour le produit"})
		return
	}

	// 4. Redirection vers la liste complète ou gestion API
	if c.Request.Header.Get("Content-Type") == "application/x-www-form-urlencoded" {
		c.Redirect(http.StatusSeeOther, "/products")
		return
	}

	c.JSON(http.StatusOK, product)
}
