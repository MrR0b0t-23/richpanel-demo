package controllers

import (
	"context"
	"net/http"
	"server/database"
	helpers "server/helpers"
	models "server/models"
	responses "server/responses"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)

var userCollection *mongo.Collection = database.GetCollection(database.DB, "userDetails")
var validate = validator.New()

func Signup() gin.HandlerFunc{ 
	return func(c *gin.Context){
		ctx, cancel := context.WithTimeout(context.Background(), 100*time.Second)
		var user models.User
		var checkUser models.User

		defer cancel()

		//validate the request body
		if err := c.BindJSON(&user); err != nil{
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		//using validator library to validate required fields
		if validationErr := validate.Struct(&user); validationErr != nil{
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": validationErr.Error()}})
			return
		}

		err := userCollection.FindOne(ctx, bson.M{"emailid": user.EmailId}).Decode(&checkUser)
		if err == nil{
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": "User Already Exist"}})
			return
		}
		
		
		hashPassword, err := helpers.HashPassword(user.Password)
		if err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"Error": err.Error()})
			return
		}
		
		newUser := models.User{
			Id: primitive.NewObjectID(),
			UserName: user.UserName,
			EmailId: user.EmailId,
			Password: hashPassword,
		}

		result, err := userCollection.InsertOne(ctx, newUser)
		if err != nil{
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}

		c.JSON(http.StatusCreated, responses.UserResponse{Status: http.StatusCreated, Message: "success", Data: map[string]interface{}{"data": result}})
    }
}


func Login() gin.HandlerFunc{
	return func(c *gin.Context){
		var ctx, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		EmailId := c.Param("emailId")
		Password := c.Param("password")
		var user models.User
		defer cancel()

		err := userCollection.FindOne(ctx, bson.M{"emailid": EmailId}).Decode(&user)
		if err != nil{
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		
		if ! helpers.CheckPasswordHash(Password, user.Password){
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Incorrect Password"}})
			return
		}

		jwtToken := GenerateJwtToken(EmailId)

		c.JSON(http.StatusOK, responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": user, "token": jwtToken}})
}
}


type JwtResponse struct{
	JwtToken string `json:"token"`
}

func OnLoad() gin.HandlerFunc{
	return func(c *gin.Context){
		var _, cancel = context.WithTimeout(context.Background(), 10000*time.Second)
		defer cancel()
		
		var jwtRes JwtResponse

		//validate the request body
		if err := c.BindJSON(&jwtRes); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"status":http.StatusBadRequest, "message":"failure", "error": err.Error()})
			return
		}

		isValid, emailId, _ := VerifyJwtToken(jwtRes.JwtToken)

		if isValid{
			c.JSON(http.StatusOK, gin.H{"status":http.StatusOK, "message": "success", "data": emailId})
		return
		}
		c.JSON(http.StatusOK, gin.H{"status": http.StatusOK, "message": "failure", "data": "Unauthorized Acess"})

	}
}