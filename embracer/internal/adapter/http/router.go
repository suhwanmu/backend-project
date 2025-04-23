package http

import (
	"embracer/internal/app"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func NewRouter(usecase *app.LoadTestUsecase) *gin.Engine {
	r := gin.Default()
	r.GET("/ping", func(c *gin.Context) {
		c.String(200, "pong")
	})

	load := r.Group("/load")
	{
		load.GET("/cpu", func(c *gin.Context) {
			duration, _ := strconv.Atoi(c.DefaultQuery("duration", "10"))
			if err := usecase.ExecuteCPULoad(duration); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "CPU load executed"})
		})

		load.GET("/memory", func(c *gin.Context) {
			duration, _ := strconv.Atoi(c.DefaultQuery("duration", "10"))
			mb, _ := strconv.Atoi(c.DefaultQuery("mb", "100"))
			if err := usecase.ExecuteMemoryLoad(duration, mb); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Memory load executed"})
		})

		load.GET("/disk", func(c *gin.Context) {
			duration, _ := strconv.Atoi(c.DefaultQuery("duration", "10"))
			if err := usecase.ExecuteDiskIOLoad(duration); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Disk IO load executed"})
		})

		load.GET("/network", func(c *gin.Context) {
			duration, _ := strconv.Atoi(c.DefaultQuery("duration", "10"))
			url := c.DefaultQuery("url", "https://httpbin.org/get")
			if err := usecase.ExecuteNetworkIOLoad(duration, url); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Network IO load executed"})
		})

		load.GET("/mix", func(c *gin.Context) {
			duration, _ := strconv.Atoi(c.DefaultQuery("duration", "10"))
			url := c.DefaultQuery("url", "https://httpbin.org/get")
			mb, _ := strconv.Atoi(c.DefaultQuery("mb", "100"))
			if err := usecase.ExecuteMixedLoad(duration, url, mb); err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
				return
			}
			c.JSON(http.StatusOK, gin.H{"message": "Mixed load executed"})
		})
	}

	return r
}
