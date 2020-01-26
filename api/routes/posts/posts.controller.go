package posts

import (
	"fmt"

	"github.com/gin-gonic/gin"
	"github.com/jinzhu/gorm"
	"github.com/nireo/go-blog-api/database/models"
	"github.com/nireo/go-blog-api/lib/common"
)

// Define topics, so we can check if the topic given in request is valid
var topics = [5]string{"programming", "ai", "fitness", "self-improvement", "technology"}

// Post type alias
type Post = models.Post

// User type alias
type User = models.User

// JSON type alias
type JSON = common.JSON

// function for checking if topic is valid
func checkIfValid(topic string) bool {
	valid := false
	for _, v := range topics {
		if v == topic {
			valid = true
		}
	}
	return valid
}

func create(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)

	// interface fro request
	type RequestBody struct {
		Text        string `json:"text" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
		Topic       string `json:"topic" binding:"required"`
		ImageURL    string `json:"imageURL" binding:"required"`
	}

	var requestBody RequestBody
	// if the body is same as the RequestBody interface
	if err := c.BindJSON(&requestBody); err != nil {
		fmt.Println(err)
		c.AbortWithStatus(400)
		return
	}

	// check if topic is valid
	if !checkIfValid(requestBody.Topic) {
		c.AbortWithStatus(400)
		return
	}

	// get the user for posts
	user := c.MustGet("user").(User)

	// take title, description and text from body, but set likes to 0
	post := Post{
		Text:        requestBody.Text,
		User:        user,
		Title:       requestBody.Title,
		Likes:       0,
		Description: requestBody.Description,
		Topic:       requestBody.Topic,
		ImageURL:    requestBody.ImageURL,
	}

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
	//  get the posts from topic
	topic := c.DefaultQuery("topic", "none")
	if topic != "none" {
		if err := db.Preload("User").Where("topic = ?", topic).Limit(10).Find(&posts).Error; err != nil {
			// if no posts in category are found
			c.AbortWithStatus(404)
			return
		}

		// serialize data
		serialized := make([]JSON, len(posts), len(posts))
		for index := range posts {
			serialized[index] = posts[index].Serialize()
		}

		c.JSON(200, serialized)

	} else {
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

		serialized := make([]JSON, len(posts), len(posts))
		for index := range posts {
			serialized[index] = posts[index].Serialize()
		}

		c.JSON(200, serialized)
	}
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

func update(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")

	user := c.MustGet("user").(User)

	type RequestBody struct {
		Text        string `json:"text" binding:"required"`
		Title       string `json:"title" binding:"required"`
		Description string `json:"description" binding:"required"`
	}

	var requestBody RequestBody

	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatus(400)
		return
	}

	var post Post
	if err := db.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	// check if the user owns the post
	if post.UserID != user.ID {
		// return status 403 - forbidden
		c.AbortWithStatus(403)
		return
	}

	// reassign post values
	post.Text = requestBody.Text
	post.Title = requestBody.Title
	post.Description = requestBody.Description

	// save to database
	db.Save(&post)
	c.JSON(200, post.Serialize())
}

func handleLike(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("postId")

	type RequestBody struct {
		Like string `json:"text" binding:"required"`
	}

	var requestBody RequestBody
	if err := c.BindJSON(&requestBody); err != nil {
		c.AbortWithStatus(400)
		return
	}

	var post Post
	if err := db.Preload("User").Where("id = ?", id).First(&post).Error; err != nil {
		c.AbortWithStatus(403)
		return
	}

	post.Likes = post.Likes + 1
	db.Save(&post)
	c.JSON(200, post.Serialize())
}

func remove(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	id := c.Param("id")
	user := c.MustGet("user").(User)

	var post Post
	if err := db.Where("id = ?", id).First(&post).Error; err != nil {
		// post not found
		c.AbortWithStatus(404)
		return
	}

	if post.UserID != user.ID {
		// user doesn't own the post
		c.AbortWithStatus(403)
		return
	}

	db.Delete(&post)
	c.Status(204)
}

func yourBlogs(c *gin.Context) {
	db := c.MustGet("db").(*gorm.DB)
	user := c.MustGet("user").(User)

	var posts []Post
	if err := db.Model(&user).Related(&posts).Error; err != nil {
		c.AbortWithStatus(404)
		return
	}

	serialized := make([]JSON, len(posts), len(posts))
	for index := range posts {
		serialized[index] = posts[index].Serialize()
	}

	c.JSON(200, serialized)
}
