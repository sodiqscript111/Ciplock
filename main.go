package main

import (
	"Ciplock/cliplockdb"
	"Ciplock/middleware"
	"Ciplock/models"
	"Ciplock/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"log"
	"net/http"
	"time"
)

func main() {
	log.Println("Starting Ciplock...")

	cliplockdb.InitDB()
	log.Println("DB setup complete.")

	server := gin.Default()

	server.POST("/signup", Addusers)
	server.POST("/login", ValidateUser)
	server.POST("/auth/login/:project_id", ValidateCustomer)
	server.POST("/auth/signup/:project_id", SignupCustomer)
	authenticated := server.Group("/")
	authenticated.Use(middleware.Authorize())

	authenticated.POST("/addproject", AddProject)
	authenticated.GET("/projects", GetProjects)

	log.Println("Server running on :8080")
	if err := server.Run(":8080"); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
func Addusers(c *gin.Context) {
	var user models.Admin

	err := c.ShouldBindJSON(&user)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	password, err := utils.HashPassword(user.PasswordHash)
	if err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}
	user.CreatedAt = time.Now()
	user.UpdatedAt = time.Now()
	user.PasswordHash = password

	err = models.AddAdmin(user)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to add user"})
		return
	}

	c.JSON(200, gin.H{
		"message": "User added successfully",
		"user": gin.H{
			"id":         user.ID,
			"email":      user.Email,
			"created_at": user.CreatedAt.Format(time.RFC3339),
			"status":     "active",
		},
	})
}

func ValidateUser(c *gin.Context) {
	var loginPayload struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&loginPayload); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	user := models.Admin{
		Email: loginPayload.Email,
	}

	ok, err := user.ValidateUser(loginPayload.Password)
	if err != nil {
		c.JSON(401, gin.H{"error": "Authentication failed", "details": err.Error()})
		return
	}

	if !ok {
		c.JSON(401, gin.H{"error": "Invalid credentials"})
		return
	}

	token, err := utils.GenerateToken(user.Email, user.ID)
	if err != nil {
		c.JSON(500, gin.H{"error": "Could not generate token"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Login successful",
		"user_id": user.ID,
		"token":   token,
	})
}

func AddProject(c *gin.Context) {
	var project models.Project
	if err := c.ShouldBindJSON(&project); err != nil {
		c.JSON(400, gin.H{"error": "Invalid request payload"})
		return
	}

	adminID, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	apiKey, err := utils.GenerateAPIKey()
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to generate API key"})
		return
	}

	project.ID = uuid.New().String()
	project.APIKey = apiKey
	project.AdminID = adminID.(string)
	project.CreatedAt = time.Now()
	project.UpdatedAt = time.Now()
	project.IsActive = true

	if err := models.AddProjects(project); err != nil {
		c.JSON(500, gin.H{"error": "Failed to add project"})
		return
	}

	c.JSON(200, gin.H{
		"message": "Project created successfully",
		"project": gin.H{
			"id":      project.ID,
			"api_key": project.APIKey,
			"name":    project.Name,
		},
	})
}

func GetProjects(c *gin.Context) {
	userID, exists := c.Get("userId")
	if !exists {
		c.JSON(401, gin.H{"error": "Unauthorized"})
		return
	}

	uid, ok := userID.(string)
	if !ok {
		c.JSON(500, gin.H{"error": "Invalid user ID format"})
		return
	}

	projects, err := models.GetProjects(uid)
	if err != nil {
		c.JSON(500, gin.H{"error": "Failed to get projects"})
		return
	}

	c.JSON(200, gin.H{
		"projects": projects,
	})
}
func SignupCustomer(c *gin.Context) {
	projectIDParam := c.Param("project_id")

	projectID, err := uuid.Parse(projectIDParam)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid project_id format"})
		return
	}

	var input struct {
		FullName string `json:"full_name"`
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if input.Email == "" || input.Password == "" || input.FullName == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "All fields are required"})
		return
	}

	exists, err := models.CustomerExists(input.Email, projectID)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database error"})
		return
	}

	if exists {
		c.JSON(http.StatusConflict, gin.H{"error": "Customer already exists under this project"})
		return
	}

	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to hash password"})
		return
	}

	newUUID := uuid.New()

	customer := models.Customer{
		ID:        newUUID,
		FullName:  input.FullName,
		Email:     input.Email,
		Password:  string(hashedPassword),
		ProjectID: projectID,
		CreatedAt: time.Now(),
	}

	if err := models.CreateCustomer(customer); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusCreated, gin.H{
		"message": "Customer created successfully",
		"customer": gin.H{
			"id":         customer.ID.String(),
			"full_name":  customer.FullName,
			"email":      customer.Email,
			"project_id": customer.ProjectID.String(),
			"created_at": customer.CreatedAt.Format(time.RFC3339),
			"status":     "active",
		},
	})
}

func ValidateCustomer(c *gin.Context) {
	projectIDParam := c.Param("project_id")

	var input struct {
		Email    string `json:"email"`
		Password string `json:"password"`
	}

	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	if input.Email == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Email and password are required"})
		return
	}

	customer := models.Customer{Email: input.Email}

	valid, err := customer.ValidateCustomers(input.Password, projectIDParam)
	if err != nil || !valid {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid email or password"})
		return
	}

	customerIDStr := fmt.Sprintf("%v", customer.ID)

	token, err := utils.GenerateToken(
		customer.Email,
		customerIDStr,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to generate token"})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"message": "Customer validated successfully",
		"token":   token,
	})
}
