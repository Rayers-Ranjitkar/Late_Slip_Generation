package controllers

import (
	"context"
	"lateslip/initialializers"
	"lateslip/models"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/go-playground/validator/v10"
	"github.com/xuri/excelize/v2"
	"go.mongodb.org/mongo-driver/bson"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"go.mongodb.org/mongo-driver/mongo/options"
)

func UploadStudentData(c *gin.Context) {
	//TODO: need to update the excel cells to match the model fields
	// Right now, all the fields are assumed

	//get the file from the request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File not found"})
		return
	}

	//check if the file is an Excel file
	if !strings.HasSuffix(fileHeader.Filename, ".xlsx") {
		c.JSON(400, gin.H{"error": "Invalid file type. Please upload an Excel file."})
		return
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// Read the Excel file
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "Failed to parse Excel file",
		})
		return
	}
	defer xlsx.Close()

	// Get all rows from Sheet1
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil || len(rows) < 2 {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "Empty or invalid Excel sheet",
		})
		return
	}

	var newStudents []models.Student
	validator := validator.New()

	// First, get all existing students' emails in one query
	studentsCollection := initialializers.DB.Collection("students")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	// Get all existing emails
	cursor, err := studentsCollection.Find(ctx, bson.M{}, options.Find().SetProjection(bson.M{"email": 1}))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"error":   "Failed to fetch existing students",
		})
		return
	}
	defer cursor.Close(ctx)

	existingEmails := make(map[string]primitive.ObjectID)
	var existingStudent struct {
		ID    primitive.ObjectID `bson:"_id"`
		Email string             `bson:"email"`
	}
	for cursor.Next(ctx) {
		if err := cursor.Decode(&existingStudent); err != nil {
			continue
		}
		existingEmails[existingStudent.Email] = existingStudent.ID
	}

	// Process Excel rows
	for i, row := range rows[1:] { // Skip header row
		if len(row) < 4 {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Invalid data format in row " + strconv.Itoa(i+2),
				"row":     row,
			})
			return
		}

		//TODO: can use switch case for name, full name , email, email address
		// 		ID            primitive.ObjectID `bson:"_id,omitempty" json:"id"`
		// StudentID     string             `bson:"student_id" json:"student_id"`
		// Name          string             `bson:"name" json:"name"`
		// Email         string             `bson:"email" json:"email"`
		// Semester      string             `bson:"semester" json:"semester"`
		// Level         string             `bson:"level" json:"level"`
		// LateSlipCount int                `bson:"late_slip_count" json:"late_slip_count"`

		student := models.Student{
			StudentID:     row[0],
			Name:          row[1],
			Email:         row[2],
			Semester:      row[3],
			Level:         row[4],
			LateSlipCount: 0, // Default value
		}

		// Validate the student struct
		if err := validator.Struct(student); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Invalid data in row " + strconv.Itoa(i+2) + ": " + err.Error(),
				"row":     row,
			})
			return
		}

		// Check if email exists in our map
		if _, exists := existingEmails[student.Email]; !exists {
			student.ID = primitive.NewObjectID()
			newStudents = append(newStudents, student)
		}
	}

	// Bulk insert new students
	if len(newStudents) > 0 {
		var documents []interface{}
		for _, student := range newStudents {
			documents = append(documents, student)
		}

		_, err = studentsCollection.InsertMany(ctx, documents)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"error":   "Failed to insert new students",
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Excel file processed successfully",
		"stats": gin.H{
			"new":      len(newStudents),
			"existing": len(existingEmails),
		},
	})
}

func UploadScheduleData(c *gin.Context) {
	//get the file from the request
	fileHeader, err := c.FormFile("file")
	if err != nil {
		c.JSON(400, gin.H{"error": "File not found"})
		return
	}

	//check if the file is an Excel file
	if !strings.HasSuffix(fileHeader.Filename, ".xlsx") {
		c.JSON(400, gin.H{"error": "Invalid file type. Please upload an Excel file."})
		return
	}

	// Open the uploaded file
	file, err := fileHeader.Open()
	if err != nil {
		c.JSON(400, gin.H{"error": "Failed to open file"})
		return
	}
	defer file.Close()

	// Read the Excel file
	xlsx, err := excelize.OpenReader(file)
	if err != nil {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "Failed to parse Excel file",
		})
		return
	}
	defer xlsx.Close()

	// Get all rows from Sheet1
	rows, err := xlsx.GetRows("Sheet1")
	if err != nil || len(rows) < 2 {
		c.JSON(400, gin.H{
			"success": false,
			"error":   "Empty or invalid Excel sheet",
		})
		return
	}

	var newSchedules []models.Schedule
	validator := validator.New()

	// First, get all existing schedules' names in one query
	schedulesCollection := initialializers.DB.Collection("schedules")
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()

	//process excel rows
	for i, row := range rows[1:] { // Skip header row
		if len(row) < 4 {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Invalid data format in row " + strconv.Itoa(i+2),
				"row":     row,
			})
			return
		}

		schedule := models.Schedule{
			ID:             primitive.NewObjectID(),
			ModuleCode:     row[0],
			ModuleName:     row[1],
			StartTime:      row[2],
			EndTime:        row[3],
			Day:            row[4],
			RoomName:       row[5],
			InstructorName: row[6],
			Semester:       row[7],
			CreatedAt:      time.Now(),
			UpdatedAt:      time.Now(),
		}

		// Validate the schedule struct
		if err := validator.Struct(schedule); err != nil {
			c.JSON(400, gin.H{
				"success": false,
				"error":   "Invalid data in row " + strconv.Itoa(i+2) + ": " + err.Error(),
				"row":     row,
			})
			return
		}
		newSchedules = append(newSchedules, schedule)

	}
	if len(newSchedules) > 0 {
		var documents []any
		for _, schedule := range newSchedules {
			documents = append(documents, schedule)
		}

		_, err = schedulesCollection.InsertMany(ctx, documents)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{
				"success": false,
				"mesasge": "Failed to insert new schedules",
				"error":   err.Error(),
			})
			return
		}
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "Excel file processed successfully",
		"stats": gin.H{
			"new": len(newSchedules),
		},
	})

}
