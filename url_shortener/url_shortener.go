package main
import (
	"net/http"
	"math/rand"
	"time"
	"github.com/gin-gonic/gin"
)

var urlMap = make(map[string]string)
func main() {
	rand.Seed(time.Now().UnixNano())
	r := gin.Default()
	r.POST("/shorten", shortenURL)
	r.GET("/:shortURL", redirectURL)
	r.Run(":8080")
}

func shortenURL(c *gin.Context) {
	longURL := c.PostForm("url")
	shortURL := generateShortURL()
	urlMap[shortURL] = longURL
	c.JSON(http.StatusOK, gin.H{"shortURL": shortURL})
}

func redirectURL(c *gin.Context) {
	shortURL := c.Param("shortURL")
	longURL, ok := urlMap[shortURL]
	if !ok {
		c.JSON(http.StatusNotFound, gin.H{"message": "Short URL not found"})
		return
	}
	c.Redirect(http.StatusMovedPermanently, longURL)
}

func generateShortURL() string {
	const charset = "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	shortURL := make([]byte, 6)
	for i := range shortURL {
		shortURL[i] = charset[rand.Intn(len(charset))]
	}
	return string(shortURL)
}