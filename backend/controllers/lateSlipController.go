package controllers

import (
	"context"
	"fmt"

	"lateslip/events"
	"lateslip/initialializers"
	"lateslip/models"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/sendgrid/sendgrid-go"
	"github.com/sendgrid/sendgrid-go/helpers/mail"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
)

// sendEmail sends an email using SendGrid asynchronously
func sendEmail(toEmail, subject, content string) {
	go func() {
		// Check if API key is set
		apiKey := os.Getenv("SENDGRID_API_KEY")
		if apiKey == "" {
			log.Printf("SendGrid API key not found")
			return
		}

		fromEmail := os.Getenv("SENDGRID_FROM_EMAIL")
		if fromEmail == "" {
			log.Printf("SendGrid from email not found")
			return
		}

		from := mail.NewEmail("HeraldSync Late Slip System", fromEmail)
		to := mail.NewEmail("Admin", toEmail)

		// Create HTML content
		htmlContent := fmt.Sprintf(`
            <div style="font-family: Arial, sans-serif; padding: 20px;">
                <h2>New Late Slip Request</h2>
                <p>%s</p>
                <p>This is an automated message from the HeraldSync Late Slip System.</p>
            </div>
        `, content)

		message := mail.NewSingleEmail(from, subject, to, content, htmlContent)

		// // Add custom headers
		// message.SetHeader("List-Unsubscribe", "<mailto:unsubscribe@heraldsync.com>")
		// message.SetHeader("Precedence", "bulk")

		client := sendgrid.NewSendClient(apiKey)

		log.Printf("Attempting to send email to: %s", toEmail)
		response, err := client.Send(message)
		if err != nil {
			log.Printf("Failed to send email: %v", err)
			return
		}

		if response.StatusCode >= 300 {
			log.Printf("Email failed with status %d: %s", response.StatusCode, response.Body)
			return
		}

		log.Printf("Email sent successfully! Status: %d", response.StatusCode)
	}()
}

func RequestLateSlip(c *gin.Context) {
	//TODO: need to check if student's late slip request limit is reached or not (max is 4 per semester)
	//If the limit is reached, return an error response

	//get student ID from context and reason from request body
	userId, exists := c.Get("user_id")
	requestID := c.GetString("request_id")
	if !exists {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "User not authenticated"})
		return
	}

	studentID, err := primitive.ObjectIDFromHex(userId.(string))
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid user ID"})
		return
	}

	type requestBody struct {
		Reason string `json:"reason" binding:"required"`
	}
	var body requestBody
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	//create a new late slip request
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlip := models.LateSlip{
		ID:        primitive.NewObjectID(),
		RequestID: requestID,
		StudentID: studentID,
		Reason:    body.Reason,
		Status:    "pending",
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	//insert the late slip into the database
	lateSlipCollection := initialializers.DB.Collection("lateslips")
	_, err = lateSlipCollection.InsertOne(ctx, lateSlip)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create late slip"})
		return
	}

	//TODO: send notification to admin
	// This could be done via email, push notification, etc.
	sendEmail(
		"samreedmaharjan1899@gmail.com",
		"New Late Slip Request Notification",
		fmt.Sprintf(`
        A new late slip request has been submitted with the following details:
        
        Student ID: %s
        Reason: %s
        Status: %s
        Submitted: %s
    `, studentID.Hex(), body.Reason, lateSlip.Status, lateSlip.CreatedAt.Format("Jan 2, 2006 3:04 PM")),
	)
	events.NotifyAdmins(
		fmt.Sprintf("New late slip request from %s", studentID.Hex()),
		lateSlip,
	)

	//return the late slip
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Late slip created successfully", "lateSlip": lateSlip})

}

func ApproveLateSlip(c *gin.Context) {
	type Body struct {
		LateSlipID string `json:"lateSlipId" binding:"required"`
		StudentID  string `json:"studentId" binding:"required"`
	}
	var body Body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lateSlipID, err := primitive.ObjectIDFromHex(body.LateSlipID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}
	studentID, err := primitive.ObjectIDFromHex(body.StudentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlipCollection := initialializers.DB.Collection("lateslips")
	var lateSlip models.LateSlip
	err = lateSlipCollection.FindOne(ctx, bson.M{"_id": lateSlipID}).Decode(&lateSlip)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch late slip",
		})
		return
	}

	if lateSlip.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Late slip is already " + lateSlip.Status,
		})
		return
	}
	//TODO: Replace the User model with the Student model
	UserCollection := initialializers.DB.Collection("users")
	var student models.User
	err = UserCollection.FindOne(ctx, bson.M{"_id": studentID}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch student details",
		})
		return
	}

	//TODO: need to update the lateslip count logic rignt now
	// --- should use student model instead of lateslip model

	// // Increment late slip count
	// _, err = StudentCollection.UpdateOne(
	// 	ctx,
	// 	bson.M{"_id": lateSlip.StudentID},
	// 	bson.M{"$inc": bson.M{"lateSlipCount": 1}},
	// )
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"success": false,
	// 		"error":   "Failed to update student late slip count",
	// 	})
	// 	return
	// }

	// Update late slip status
	lateSlip.Status = "approved"
	lateSlip.UpdatedAt = time.Now()

	//update the late slip in the database
	_, err = lateSlipCollection.UpdateOne(ctx, bson.M{"_id": lateSlipID}, bson.M{"$set": lateSlip})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update late slip"})
		return
	}

	//TODO: send notification to student
	// This could be done via email, push notification, etc.
	sendEmail(
		"bivekshrestha239@gmail.com",
		"Late Slip Approval Notification",
		fmt.Sprintf(`
        Your late slip request has been approved.
        
        Late Slip ID: %s
        Reason: %s
        Status: %s
        Approved: %s
    `, lateSlip.ID.Hex(), lateSlip.Reason, lateSlip.Status, lateSlip.UpdatedAt.Format("Jan 2, 2006 3:04 PM")),
	)
	events.NotifyStudent(
		studentID.Hex(),
		fmt.Sprintf("Your late slip request has been %s", lateSlip.Status),
	)

	//return the late slip
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Late slip approved successfully", "lateSlip": lateSlip})
}

func GetAllLateSlips(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlipCollection := initialializers.DB.Collection("lateslips")
	cursor, err := lateSlipCollection.Find(ctx, bson.M{})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch late slips"})
		return
	}
	defer cursor.Close(ctx)

	var lateSlips []models.LateSlip
	if err = cursor.All(ctx, &lateSlips); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode late slips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "lateSlips": lateSlips})
}

// GET /admin/lateslip/requests
func GetAllPendingLateSlip(c *gin.Context) {
	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlipCollection := initialializers.DB.Collection("lateslips")
	cursor, err := lateSlipCollection.Find(ctx, bson.M{"status": "pending"})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch late slips"})
		return
	}
	defer cursor.Close(ctx)

	var lateSlips []models.LateSlip
	if err = cursor.All(ctx, &lateSlips); err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to decode late slips"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"success": true, "lateSlips": lateSlips})

}

// reject late slip
func RejectLateSlip(c *gin.Context) {
	type Body struct {
		LateSlipID string `json:"lateSlipId" binding:"required"`
		StudentID  string `json:"studentId" binding:"required"`
	}
	var body Body
	if err := c.BindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	lateSlipID, err := primitive.ObjectIDFromHex(body.LateSlipID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid late slip ID"})
		return
	}
	studentID, err := primitive.ObjectIDFromHex(body.StudentID)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid student ID"})
		return
	}

	ctx, cancel := context.WithTimeout(c.Request.Context(), 10*time.Second)
	defer cancel()

	lateSlipCollection := initialializers.DB.Collection("lateslips")
	var lateSlip models.LateSlip
	err = lateSlipCollection.FindOne(ctx, bson.M{"_id": lateSlipID}).Decode(&lateSlip)
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch late slip",
		})
		return
	}

	if lateSlip.Status != "pending" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"error":   "Late slip is already " + lateSlip.Status,
		})
		return
	}
	//TODO: Replace the User model with the Student model
	UserCollection := initialializers.DB.Collection("users")
	var student models.User
	err = UserCollection.FindOne(ctx, bson.M{"_id": studentID}).Decode(&student)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch student details",
		})
		return
	}

	//TODO: need to update the lateslip count logic rignt now
	// --- should use student model instead of lateslip model

	// // Increment late slip count
	// _, err = StudentCollection.UpdateOne(
	// 	ctx,
	// 	bson.M{"_id": lateSlip.StudentID},
	// 	bson.M{"$inc": bson.M{"lateSlipCount": 1}},
	// )
	// if err != nil {
	// 	c.JSON(http.StatusInternalServerError, gin.H{
	// 		"success": false,
	// 		"error":   "Failed to update student late slip count",
	// 	})
	// 	return
	// }

	// Update late slip status
	lateSlip.Status = "rejected"
	lateSlip.UpdatedAt = time.Now()

	//update the late slip in the database
	_, err = lateSlipCollection.UpdateOne(ctx, bson.M{"_id": lateSlipID}, bson.M{"$set": lateSlip})
	if err != nil {
		c.Error(err)
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update late slip"})
		return
	}

	//TODO: send notification to student
	// This could be done via email, push notification, etc.
	sendEmail(
		"bivekshrestha239@gmail.com",
		"Late Slip Rejection Notification",
		fmt.Sprintf(`
        Your late slip request has been rejected.
        
        Late Slip ID: %s
        Reason: %s
        Status: %s
        Approved: %s
    `, lateSlip.ID.Hex(), lateSlip.Reason, lateSlip.Status, lateSlip.UpdatedAt.Format("Jan 2, 2006 3:04 PM")),
	)
	events.NotifyStudent(
		studentID.Hex(),
		fmt.Sprintf("Your late slip request has been %s", lateSlip.Status),
	)

	//return the late slip
	c.JSON(http.StatusOK, gin.H{"success": true, "message": "Late slip rejected successfully", "lateSlip": lateSlip})
}
