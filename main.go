package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/saransh-khobragade/golang-redis/cache"
)

var (
	redisCache = cache.NewRedisCache(os.Getenv("REDIS_CONNECTION_STRING"), 0, 1)
)

func CORSMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", "*")
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Content-Type, Content-Length, Accept-Encoding, X-CSRF-Token, Authorization, accept, origin, Cache-Control, X-Requested-With")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT")

		if c.Request.Method == "OPTIONS" {
			c.AbortWithStatus(204)
			return
		}

		c.Next()
	}
}

func main() {
	r := gin.Default()

	r.GET("/healthcheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{
			"status": "ok",
		})
	})

	r.POST("/movies", func(ctx *gin.Context) {
		var movie cache.Movie
		if err := ctx.ShouldBind(&movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		fmt.Println(movie)

		res, err := redisCache.CreateMovie(&movie)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res,
		})

	})
	r.GET("/movies", func(ctx *gin.Context) {
		movies, err := redisCache.GetMovies()
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movies": movies,
		})
	})
	r.GET("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		movie, err := redisCache.GetMovie(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"message": "movie not found",
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": movie,
		})
	})
	r.PUT("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		res, err := redisCache.GetMovie(id)
		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		var movie cache.Movie

		if err := ctx.ShouldBind(&movie); err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}

		res.Title = movie.Title
		res.Description = movie.Description
		res, err = redisCache.UpdateMovie(res)

		if err != nil {
			ctx.JSON(http.StatusBadRequest, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"movie": res,
		})
	})
	r.DELETE("/movies/:id", func(ctx *gin.Context) {
		id := ctx.Param("id")
		err := redisCache.DeleteMovie(id)
		if err != nil {
			ctx.JSON(http.StatusNotFound, gin.H{
				"error": err.Error(),
			})
			return
		}
		ctx.JSON(http.StatusOK, gin.H{
			"message": "movie deleted successfully with id: " + id,
		})
	})
	r.Use(CORSMiddleware())
	value := os.Getenv("PORT")
	fmt.Println(r.Run(":" + value))
}
