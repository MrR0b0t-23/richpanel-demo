package main

import (
	"os"
	database "server/database"
	routes "server/routes"

	"github.com/gin-contrib/cors"
	"github.com/gin-gonic/gin"
)

func main() {

	port := os.Getenv("PORT")
	router := gin.New()
	router.Use(gin.Logger())
	
	router.Use(cors.New(CORSConfig()))
  
	database.ConnectDB()
	routes.AllRoute(router)

	router.Run(":" + port)
}

func CORSConfig() cors.Config {
    corsConfig := cors.DefaultConfig()
    corsConfig.AllowOrigins = []string{"http://richpanel-react.herokuapp.com"}
    corsConfig.AllowCredentials = true
    corsConfig.AddAllowHeaders("Access-Control-Allow-Headers", "access-control-allow-origin, access-control-allow-headers", "Content-Type", "X-XSRF-TOKEN", "Accept", "Origin", "X-Requested-With", "Authorization")
    corsConfig.AddAllowMethods("GET", "POST", "PUT", "DELETE")
    return corsConfig
}
