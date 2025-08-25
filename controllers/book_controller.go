package controllers

import (
	"database/sql"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"bookstore-api/config"
	"bookstore-api/models"
)

type CreateBookRequest struct {
	Title       string `json:"title" binding:"required,min=2"`
	Description string `json:"description"`
	ImageURL    string `json:"image_url"`
	ReleaseYear int    `json:"release_year" binding:"required"`
	Price       int    `json:"price" binding:"required"`
	TotalPage   int    `json:"total_page" binding:"required"`
	CategoryID  int    `json:"category_id" binding:"required"`
}

func thicknessFromPage(p int) string {
	if p >= 100 { return "tebal" }
	return "tipis"
}

func yearValid(y int) bool { return y >= 1980 && y <= 2024 }

func GetBooks(c *gin.Context) {
	rows, err := config.DB.Query("SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM books ORDER BY id DESC")
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

func CreateBook(c *gin.Context) {
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	if !yearValid(req.ReleaseYear) { c.JSON(http.StatusBadRequest, gin.H{"error":"release_year must be between 1980 and 2024"}); return }
	th := thicknessFromPage(req.TotalPage)

	// ensure category exists
	var exists int
	if err := config.DB.QueryRow("SELECT COUNT(1) FROM categories WHERE id=$1", req.CategoryID).Scan(&exists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return
	}
	if exists == 0 { c.JSON(http.StatusBadRequest, gin.H{"error":"category_id not found"}); return }

	_, err := config.DB.Exec(`INSERT INTO books (title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by)
		VALUES ($1,$2,$3,$4,$5,$6,$7,$8,NOW(),'system',NOW(),'system')`,
		req.Title, req.Description, req.ImageURL, req.ReleaseYear, req.Price, req.TotalPage, th, req.CategoryID)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"failed to create"}); return }
	c.JSON(http.StatusCreated, gin.H{"message":"book created"})
}

func GetBookByID(c *gin.Context) {
	id := c.Param("id")
	var b models.Book
	err := config.DB.QueryRow("SELECT id, title, description, image_url, release_year, price, total_page, thickness, category_id, created_at, created_by, modified_at, modified_by FROM books WHERE id=$1", id).
		Scan(&b.ID,&b.Title,&b.Description,&b.ImageURL,&b.ReleaseYear,&b.Price,&b.TotalPage,&b.Thickness,&b.CategoryID,&b.CreatedAt,&b.CreatedBy,&b.ModifiedAt,&b.ModifiedBy)
	if err == sql.ErrNoRows { c.JSON(http.StatusNotFound, gin.H{"error":"book not found"}); return }
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	c.JSON(http.StatusOK, b)
}

func UpdateBook(c *gin.Context) {
	id := c.Param("id")
	var req CreateBookRequest
	if err := c.ShouldBindJSON(&req); err != nil { c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()}); return }
	if !yearValid(req.ReleaseYear) { c.JSON(http.StatusBadRequest, gin.H{"error":"release_year must be between 1980 and 2024"}); return }
	th := thicknessFromPage(req.TotalPage)

	// validate category
	var exists int
	if err := config.DB.QueryRow("SELECT COUNT(1) FROM categories WHERE id=$1", req.CategoryID).Scan(&exists); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return
	}
	if exists == 0 { c.JSON(http.StatusBadRequest, gin.H{"error":"category_id not found"}); return }

	res, err := config.DB.Exec(`UPDATE books SET title=$1, description=$2, image_url=$3, release_year=$4, price=$5, total_page=$6, thickness=$7, category_id=$8, modified_at=$9, modified_by='system' WHERE id=$10`,
		req.Title, req.Description, req.ImageURL, req.ReleaseYear, req.Price, req.TotalPage, th, req.CategoryID, time.Now(), id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	affected, _ := res.RowsAffected()
	if affected == 0 { c.JSON(http.StatusNotFound, gin.H{"error":"book not found for update"}); return }
	c.JSON(http.StatusOK, gin.H{"message":"book updated"})
}

func DeleteBook(c *gin.Context) {
	id := c.Param("id")
	res, err := config.DB.Exec("DELETE FROM books WHERE id=$1", id)
	if err != nil { c.JSON(http.StatusInternalServerError, gin.H{"error":"db error"}); return }
	affected, _ := res.RowsAffected()
	if affected == 0 { c.JSON(http.StatusNotFound, gin.H{"error":"book not found for delete"}); return }
	c.JSON(http.StatusOK, gin.H{"message":"book deleted"})
}
