package controllers

import (
	"context"
	"fmt"
	"net/http"
	"server/database"
	"server/models"
	responses "server/responses"
	"time"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo"
)


var subcriptionCollection *mongo.Collection = database.GetCollection(database.DB, "subscriptionDetails")

func UserPlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		defer cancel()
		var plan models.UserPlan
		var jwtRes JwtResponse

		//validate the request body
		if err := c.BindJSON(&jwtRes); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"status":http.StatusBadRequest,  "message": "Invalid resquest", "error": err.Error()})
			return
		}
		
		isValid, emailId, _ := VerifyJwtToken(jwtRes.JwtToken)

		if !isValid{
			c.JSON(http.StatusInternalServerError, gin.H{"status":http.StatusInternalServerError, "message": "Unauthorised Access"})
			return
		}

		err := subcriptionCollection.FindOne(c, bson.M{"emailid": emailId}).Decode(&plan)
	
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"status":http.StatusInternalServerError, "message": "No user Found", "error": err.Error() })
			return
		}

		c.JSON(http.StatusOK, gin.H{"status":http.StatusOK,  "message": "success", "data": plan})
	}
}

var plansCollection *mongo.Collection = database.GetCollection(database.DB, "planDetails")

func AllPlans()gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		var plans []models.Plans

		defer cancel()
      
        results, err := plansCollection.Find(ctx, bson.M{})
      
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }
      
        //reading from the db in an optimal way
        defer results.Close(ctx)
        for results.Next(ctx) {
            var singlePlan models.Plans
            if err = results.Decode(&singlePlan); err != nil {
                c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            }
          
            plans = append(plans, singlePlan)
        }
      
        c.JSON(http.StatusOK,
            responses.UserResponse{Status: http.StatusOK, Message: "success", Data: map[string]interface{}{"data": plans}},
        )

	}
}

type ActivePlanResponse struct{
	JwtToken string `json:"token"`
	PlanName string `json:"planName"`
	VideoQuality string `json:"videoQuality"`
	Resolution string `json:"resolution"`
	Device []string `json:"device"`
	Price int `json:"price"`
	Type string `json:"type"`
	ActiveScreen string `json:"activeScreens"`
}

func ActivePlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
	
		var resp ActivePlanResponse
	
		//validate the request body
		if err := c.BindJSON(&resp); err != nil{
			c.JSON(http.StatusBadRequest, gin.H{"data": "Invalid Request", "error": err.Error()})
			return
		}
			
		isValid, emailId, _ := VerifyJwtToken(resp.JwtToken)
	
		// invalid jwt token 
		if !isValid{
			c.JSON(http.StatusInternalServerError, gin.H{"data": "Unauthorised Access" })
			return
		}
		

		// getting user	 plan details
		var plan models.UserPlan
		err := subcriptionCollection.FindOne(c, bson.M{"emailid": emailId}).Decode(&plan)

		//result, err := paylogCollection.InsertOne(ctx, plan)

		//if already subscriped update the subscription, else create new subscription
		
		renew := time.Now().Add( 356*(24*time.Hour))
		
		if resp.Type == "Monthly"{
			renew = time.Now().Add( 30*(24*time.Hour))
		}
		
		
		
		if err != nil {
			var newplan = models.UserPlan{
				Id: primitive.NewObjectID(),
				EmailId: emailId,
				PlanName: resp.PlanName,
				Device: resp.Device,
				Price: resp.Price,
				Type: resp.Type,
				SubcripedAt: time.Now().Format("Jan 02, 2006"), 
				RenewAt: renew.Format("Jan 02, 2006"),
				Status: "active",
			}
	
			result, err := subcriptionCollection.InsertOne(ctx, newplan)
			if err != nil{
				c.JSON(http.StatusInternalServerError, gin.H{"data": err.Error()})
				return
			}
	
			c.JSON(http.StatusCreated, gin.H{"data": "New Subcription Created", "result": result})	

		}
		update := bson.M{
			"emailid": emailId,
			"planname": resp.PlanName,
			"device": resp.Device,
			"price": resp.Price,
			"type": resp.Type,
			"status": "active",
		}
		
		result, err := subcriptionCollection.UpdateOne(ctx, bson.M{"_id": plan.Id}, bson.M{"$set": update})
		if err != nil {
			c.JSON(http.StatusInternalServerError,gin.H{"data": err.Error()})
			return
		}
		
		if result.MatchedCount == 1{

			c.JSON(http.StatusOK, gin.H{"data": "Subscription Updated"})
			return 
		}

		c.JSON(http.StatusBadRequest, gin.H{"data": err.Error()})
}}

func CancelPlan() gin.HandlerFunc {
	return func(c *gin.Context) {
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()

		var plan models.UserPlan
		var jwtRes JwtResponse

		//validate the request body
		if err := c.BindJSON(&jwtRes); err != nil{
			c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
			return
		}
		
		isValid, emailId, _ := VerifyJwtToken(jwtRes.JwtToken)

		// invalid jwt token 
		if !isValid{
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": "Unauthorised Access" }})
		 return
		}

		// getting user plan details
		fmt.Println(jwtRes)
		err := subcriptionCollection.FindOne(c, bson.M{"emailid": emailId}).Decode(&plan)
		if err != nil {
			c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err }})
			return
		}

		update := bson.M{
			"emailid": plan.EmailId,
			"planname": plan.PlanName,
			"device": plan.Device,
			"price": plan.Price,
			"type": plan.Type,
			"subcripedat": plan.SubcripedAt,
			"renewat": plan.RenewAt,
			"status": "cancelled",
		}

		result, err := subcriptionCollection.UpdateOne(ctx, bson.M{"_id": plan.Id}, bson.M{"$set": update})
        if err != nil {
            c.JSON(http.StatusInternalServerError, responses.UserResponse{Status: http.StatusInternalServerError, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
            return
        }
		
		if result.MatchedCount == 1 {
			c.Redirect(http.StatusSeeOther, "http://localhost:3000/")
		}

		c.JSON(http.StatusBadRequest, responses.UserResponse{Status: http.StatusBadRequest, Message: "error", Data: map[string]interface{}{"data": err.Error()}})
}
}