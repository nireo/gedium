package posts

import (
	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nireo/go-blog-api/database/models"
	"github.com/nireo/go-blog-api/lib/common"
)

// Post type alias
type Post = models.Post

// User type alias
type User = models.User

// JSON type alias
type JSON = common.JSON

func create(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// interface fro request
	type RequestBody struct {
		Text string `json:"text" binding:"required"`
	}

	var requestBody RequestBody
	// if the body is same as the RequestBody interface
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatus(400)
		return
	}

	// get the user for posts
	user := c.MustGet("user").(User)
	post := Post{Text: requestBody.Text, User: user}
	db.NewRecord(post)
	db.Create(&post)
	// send the new post and a success status code
	c.JSON(200, post.Serialize())
}

func list(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	cursor := c.Query("cursor")
	recent := c.Query("recent")
	var posts []Post

	if cursor == "" {
		if err := db.Preload("User").Limit(10).Order("id desc").Find(&posts).Error; err != nil {
			c.AbortWithStatus(500)
			return
		}
	} else {
		condition := "id < ?"
		if recent == "1" {
			condition = "id > ?"
		}
		if err := db.Preload("User").Limit(10).Order("id desc").Where(condition, cursor).Find(&posts).Error; err != nil {
			c.AbortWithStatus(500)
			return
		}
	}

	length := len(posts)
	serialized := make([]JSON, length, length)
	for i := 0; i < length; i++ {
		serialized[i] = posts[i].Serialize()
	}

	c.JSON(200, serialized)
}

func postFromID(c *gin.Context) {
	// get the database from gin context
	db := c.MustGet("db").(*gorm.DB)

	// get the id parameter
	id := c.Param("id")
	var post Post

	// preload related model
	if err := db.Set("gorm:auto_preload", true).Where("id = ?", id).First(&post).Error; err != nil {
		// if the has not been found
		c.AbortWithStatus(404)
		return
	}

	// return post json data and successfull data
	c.JSON(200, post.Serialize())
}