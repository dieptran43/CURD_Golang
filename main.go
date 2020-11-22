package main

import (
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	//"gorm.io/driver/postgres"
	//"gorm.io/driver/sqlite"
	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

type Post struct {
	ID        int       `json:"id" gorm:"primary key: auto_increment"`
	Title     string    `json:"title" gorm:"size 255"`
	Content   string    `json:"content" gorm:"size 255"`
	CreatedAt time.Time `json:"created_at" gorm:"CURRENT_TIMESTAMP"`
	UpdatedAt time.Time `json:"updated_at" gorm:"CURRENT_TIMESTAMP"`
}

type CreatePostParams struct {
	Title   string `json:"title"`
	Content string `json:"content"`
}

type UpdatePostParams struct {
	ID      int    `json:"id, string, int" binding:"required"`
	Title   string `json:"title" binding:"required"`
	Content string `json:"content" binding:"required"`
}
type RQHeader struct {
	Authorization string
}

var db *gorm.DB

const API_KEY = "API KEY"

func main() {
	router := gin.Default()
	router.Use(midleLogging)
	db, _ = seeding()

	v1 := router.Group("/v1")
	{
		v1.POST("/post", ensureAuthorized, createPost)
		v1.GET("/post/:id", readPost)
		v1.PUT("/post", ensureAuthorized, updatePost)
		v1.DELETE("/post/:id", ensureAuthorized, deletePost)
	}
	log.Fatal(router.Run(":8080"))
}
func createPost(context *gin.Context) {
	fmt.Println("Call Create Func")
	var params CreatePostParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}

	post := Post{
		Title:   params.Title,
		Content: params.Content,
	}
	err = db.Create(&post).Error
	if err != nil {
		context.JSON(http.StatusUnprocessableEntity, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"post": post})
	return
}

func readPost(context *gin.Context) {
	fmt.Println("Call Read Func")
	id := context.Param("id")
	post := Post{}
	err := db.First(&post, id).Error
	if err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": err})
		return
	}
	context.JSON(http.StatusOK, gin.H{"post": post})
	return
}

func updatePost(context *gin.Context) {
	fmt.Println("Call Update Func")
	var params UpdatePostParams
	err := context.ShouldBindJSON(&params)
	if err != nil {
		context.AbortWithError(http.StatusBadRequest, err)
		return
	}
	post := Post{
		ID:        params.ID,
		Title:     params.Title,
		Content:   params.Content,
		UpdatedAt: time.Time{},
	}
	err = db.Model(&post).Updates(post).Error
	if err != nil {
		context.AbortWithError(http.StatusInsufficientStorage, err)
		return
	}
	context.JSON(http.StatusOK, gin.H{"post": post})
	return
}

func deletePost(context *gin.Context) {
	fmt.Println("Call Delete Func")
	id := context.Param("id")
	var post Post
	err := db.First(&post, id).Error
	if err != nil {
		context.JSON(http.StatusOK, gin.H{"message": "Not found"})
		return
	}
	err = db.Delete(&Post{}, id).Error
	if err != nil {
		context.AbortWithStatus(http.StatusUnprocessableEntity)
		return
	}
	context.JSON(http.StatusOK, gin.H{"message": "Deleted........."})
	return

}

func midleLogging(context *gin.Context) {
	context.Next()
	log.Println(context.Request)
}

func ensureAuthorized(context *gin.Context) {
	var header RQHeader
	err := context.ShouldBindHeader(&header)
	if err != nil {
		context.AbortWithStatus(http.StatusBadRequest)
		return
	}
	if header.Authorization != API_KEY {
		context.AbortWithStatus(http.StatusNotAcceptable)
		return
	}
	context.Next()
}

func seeding() (*gorm.DB, error) {
	//db, err := gorm.Open(sqlite.Open("gorm.db"), &gorm.Config{})
	dsn := "root:root@tcp(127.0.0.1:3306)/DemoAccountOwner?charset=utf8mb4&parseTime=True&loc=Local"
	db, err := gorm.Open(mysql.Open(dsn), &gorm.Config{})
	if err != nil {
		return nil, err
	}

	err = db.AutoMigrate(&Post{})
	if err != nil {
		return nil, err
	}
	return db, nil
}
