package users

import (
	"fmt"
	"log"
	"net/http"
	"strconv"
	"errors"
	"encoding/json"
	"gorm.io/gorm"
	"fortifyApp/utils"
)

// @Tags users
// @Summary Login
// @Description Login
// @Accept json
// @Produce json
// @Success  200
// @Failure 400
// @Router /login [post]
func LoginHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB){
	
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	var req struct {
		Username string `json:"username"`
		Password string `json:"password"`
	}
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		http.Error(w, "Invalid request body", http.StatusBadRequest)
		return
	}

	var user User
	result := db.Where("name = ?", req.Username).First(&user)
	if result.Error != nil || !utils.CheckPasswordHash(req.Password, user.Password) {
		http.Error(w, "Invalid credentials", http.StatusUnauthorized)
		return
	}

	accessToken, _ := utils.CreateToken(user.Name)
	refreshToken, _ := utils.GenerateRefreshToken(user.Name)

	json.NewEncoder(w).Encode(map[string]string{
		"access_token":  accessToken,
		"refresh_token": refreshToken,
	})
}


func CreateUserHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	var user User
	err := json.NewDecoder(r.Body).Decode(&user)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	fmt.Sprintf("User: %v", user)

	// Validate input
	if user.Name == "" || user.Email == "" || user.Password == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	hashedPassword, err := utils.HashPassword(user.Password)
	if err != nil {
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}
	user.Password = hashedPassword

	err = CreateUser(&user, db)
	if err != nil{
		http.Error(w, fmt.Sprintf("Failed to create user: %v", err), http.StatusInternalServerError)
		return
	}

	response := utils.Response{
		Message: "User created and saved",
	}

	log.Printf("User Created and saved successfully")
	utils.WriteJSONResponse(w, response, http.StatusOK)
}

// @Tags users
// @Summary Update a user
// @Description Update user details and information
// @Accept json
// @Param id path int true "User ID"
// @Produce json
// @Success  200
// @Failure 400
// @Router /users/{id} [post]
func UpdateUserHandler(w http.ResponseWriter, idStr string, r *http.Request, db *gorm.DB) {

	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	var inputUser User
	if err := json.NewDecoder(r.Body).Decode(&inputUser); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	user, err := UpdateUser(id, inputUser, db)
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			http.Error(w, fmt.Sprintf("User with id = %d not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update User: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully updated user with id %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(user)
}

// @Tags users
// @Summary Get all users
// @Description Get all users
// @Accept json
// @Produce json
// @Success  200 {array} User
// @Failure 400
// @Router /users [get]
func GetUserHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {
	users, err := GetUsers(db)
	if err != nil{
		if _, ok := err.(*utils.ValidationError); ok{
			http.Error(w, "Invalid parameter", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}
	if len(users) == 0{
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(users)
}

// @Tags users
// @Summary Delete a user
// @Description Delete a user using member ID
// @Param id path int true "User ID"
// @Success  200
// @Failure 400
// @Router /users/{id} [delete]
func DeleteUserHandler(w http.ResponseWriter, idStr string, r *http.Request, db *gorm.DB) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid user ID", http.StatusBadRequest)
		return
	}
	
	err = DeleteUser(id, db)
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			http.Error(w, fmt.Sprintf("User with id = %d not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete user: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted user with id %d", id)
	w.WriteHeader(http.StatusOK)
}

