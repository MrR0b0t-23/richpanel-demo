package models

import (
	"go.mongodb.org/mongo-driver/bson/primitive"
)

type User struct {
	Id				primitive.ObjectID	 `bson:"_id"`
	UserName		string				 `json:"userName" binding:"required,max=30"`
	EmailId         string				 `json:"emailId" binding:"email,required"`
	Password   		string				 `json:"password" binding:"required,min=6,max=18"`
}

type UserPlan struct{
	Id		primitive.ObjectID	 `bson:"_id"`
	EmailId string `json:"emailId"`
	PlanName string `json:"planName"`
	Device []string `json:"device"`
	Price int `json:"price"`
	Type string `json:"type"`
	SubcripedAt string `json:"subcripedAt"`
	RenewAt string `json:"renewAt"`
	Status string `json:"status"`

}

type Plans struct{
	PlanName string `json:"planName"`
	VideoQuality string `json:"videoQuality"`
	Resolution string `json:"resolution"`
	Devices []string `json:"devices"`
	ActiveScreens string `json:"activeScreens"`
	MonthlyPrice string `json:"monthlyPrice"`
	YearlyPrice string `json:"yearlyPrice"`
	ProductID string `json:"productID"`
	MonthlyPriceAPI string `json:"monthlyPriceAPI"`
	YearlyPriceAPI string `json:"yearlyPriceAPI"`
}