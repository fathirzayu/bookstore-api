package controllers

import (
	"database/sql"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"bookstore-api/config"
	"bookstore-api/models"
)

type CreateCategoryRequest struct {
	Name string `json:"name" binding:"required,min=2"`
}

func GetCategories(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories ORDER BY id DESC")
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	defer rows.Close()
	var list []models.Category
	for rows.Next() {
		var cat models.Category
		if err := rows.Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.CreatedBy, &cat.ModifiedAt, &cat.ModifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"scan error"}); return
		}
		list = append(list, cat)
	}
	c.JSON(http.StatusOK, list)
}

func CreateCategory(c *gin.Context) {
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	_, err := config.DB.Exec("INSERT INTO categories (name, created_at, created_by, modified_at, modified_by) VALUES ($1,NOW(),'system',NOW(),'system')", req.Name)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to create"}); return }
	c.JSON(http.StatusCreated, gin.H{"message":"category created"})
}

func GetCategoryByID(c *gin.Context) {
	id := c.Param("id")
	var cat models.Category
	err := config.DB.QueryRow("SELECT id, name, created_at, created_by, modified_at, modified_by FROM categories WHERE id=$1", id).
		Scan(&cat.ID, &cat.Name, &cat.CreatedAt, &cat.CreatedBy, &cat.ModifiedAt, &cat.ModifiedBy)
	if err == sql.ErrNoRows { c.JSON(http.StatusNotFound, gin.H{"error":"category not found"}); return }
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	c.JSON(http.StatusOK, cat)
}

func UpdateCategory(c *gin.Context) {
	id := c.Param("id")
	var req CreateCategoryRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	res, err := config.DB.Exec("UPDATE categories SET name=$1, modified_at=NOW(), modified_by='system' WHERE id=$2", req.Name, id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	affected, _ := res.RowsAffected()
	if affected == 0 { c.JSON(http.StatusNotFound, gin.H{"error":"category not found for update"}); return }
	c.JSON(http.StatusOK, gin.H{"message":"category updated"})
}

func DeleteCategory(c *gin.Context) {
	id := c.Param("id")
	res, err := config.DB.Exec("DELETE FROM categories WHERE id=$1", id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	affected, _ := res.RowsAffected()
	if affected == 0 { c.JSON(http.StatusNotFound, gin.H{"error":"category not found for delete"}); return }
	c.JSON(http.StatusOK, gin.H{"message":"category deleted"})
}

func GetBooksByCategory(c *gin.Context) {
	id := c.Param("id")
	_, err := strconv.Atoi(id)
	if err != nil { c.JSON(http.StatusBadRequest, gin.H{"error":"invalid id"}); return }
	rows, err := config.DB.Query("SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM books WHERE category_id=$1 ORDER BY id DESC", id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	defer rows.Close()
	var list []models.Book
	for rows.Next() {
		var b models.Book
		if err := rows.Scan(&b.ID,&b.Title,&b.Description,&b.ImageURL,&b.ReleaseYear,&b.Price,&b.TotalPage,&b.Thickness,&b.CategoryID,&b.CreatedAt,&b.CreatedBy,&b.ModifiedAt,&b.ModifiedBy); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error":"scan error"}); return
		}
		list = append(list, b)
	}
	c.JSON(http.StatusOK, list)
}
