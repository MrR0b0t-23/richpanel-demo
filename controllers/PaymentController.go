package controllers

import (
	"context"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/stripe/stripe-go/v72"
	"github.com/stripe/stripe-go/v72/price"
	"github.com/stripe/stripe-go/v72/paymentintent"
	"github.com/stripe/stripe-go/webhook"
)

var db = make(map[string]string)

func HandleWebhook() gin.HandlerFunc {
	return func(c *gin.Context) {
	// const MaxBodyBytes = int64(65536)
	// req.Body = http.MaxBytesReader(w, c.Request.Body, MaxBodyBytes)
	payload, err := ioutil.ReadAll(c.Request.Body)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Error reading request body: %v\n", err)
		c.JSON(http.StatusServiceUnavailable, gin.H{"error": err})
		return
	}

	// This is your Stripe CLI webhook secret for testing your endpoint locally.
	endpointSecret := "whsec_KZQKBJ2bHvFrdki3h3Cq3e40r45xg9e2"
	// Pass the request body and Stripe-Signature header to ConstructEvent, along
	// with the webhook signing key.
	event, err := webhook.ConstructEvent(payload, c.Request.Header.Get("Stripe-Signature"),
		endpointSecret)

	if err != nil {
		fmt.Fprintf(os.Stderr, "Error verifying webhook signature: %v\n", err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": err}) // Return a 400 error on a bad signature
		return
	}

	// Unmarshal the event data into an appropriate struct depending on its Type
	switch event.Type {
	case "payment_intent.succeeded":
		fmt.Print("🐮🐮🐮 Subscription schedule success\n")
		// Then define and call a function to handle the event subscription_schedule.canceled
	default:
		fmt.Fprintf(os.Stderr, "Unhandled event type: %s\n", event.Type)
	}

	c.JSON(http.StatusOK, gin.H{"status": "ok"})
}
}

func AmountCalucate(priceId string) int64{
	p, _ := price.Get(priceId, nil)
	return p.UnitAmount*100
}

func CreatePaymentIntent() gin.HandlerFunc {
	return func(c *gin.Context) {
		stripe.Key = "sk_test_51LeCTcSDyoGuF7sjCJmuY2dYTbdW7hCY0UCkeKYjWpQkPyTLxhLhHaZKNLhJ4iRLtC4Qq6fNMTjF9S8QhzEk3sTq00VQ1ogCKb"
		var _, cancel = context.WithTimeout(context.Background(), 100*time.Second)
		priceId := c.Param("priceId")
		defer cancel()
		params := &stripe.PaymentIntentParams{
			Amount:   stripe.Int64(AmountCalucate(priceId)),
			Currency: stripe.String(string(stripe.CurrencyINR)),
			AutomaticPaymentMethods: &stripe.PaymentIntentAutomaticPaymentMethodsParams{
			  Enabled: stripe.Bool(true),
			},
		  }
		
		pi, err := paymentintent.New(params)
		log.Printf("pi.New: %v", pi.ClientSecret)


		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"status":http.StatusBadRequest, "error": err.Error()})
			log.Printf("pi.New: %v", err)
			return 
		}
		c.JSON(http.StatusAccepted, gin.H{"status": http.StatusAccepted, "data":pi.ClientSecret})
	}
}

