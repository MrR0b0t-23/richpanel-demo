package routes

import (
	"server/controllers"

	"github.com/gin-gonic/gin"
)

func AllRoute(router *gin.Engine){
	router.POST("/signup", controllers.Signup())
	router.POST("/login/:emailId/:password", controllers.Login())
	router.POST("/onload", controllers.OnLoad())
	router.POST("/userplan", controllers.UserPlan())
	router.GET("/getplans", controllers.AllPlans())
	router.POST("/cancelplan", controllers.CancelPlan())
	router.POST("/activeplan", controllers.ActivePlan())
	router.POST("/create-payment-intent/:priceId", controllers.CreatePaymentIntent())
	router.POST("/webhook", controllers.HandleWebhook())
	}
