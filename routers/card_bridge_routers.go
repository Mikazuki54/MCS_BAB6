package routers

import (
	"bab6/controllers"

	"github.com/gin-gonic/gin"
)

func StartServer() *gin.Engine {
	router := gin.Default()

	router.GET("/cards", controllers.GetCards)
	router.POST("/card/input/:id", controllers.InsertCard)
	router.DELETE("/card/delete/:id", controllers.DeleteCard)

	return router
}
