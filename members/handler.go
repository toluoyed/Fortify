package members

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

// @Tags members
// @Summary Create a member
// @Description Create a single fortify member
// @Accept json
// @Produce json
// @Success  200
// @Failure 400
// @Router /members [post]
func PostMemberHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {

	var member Member
	err := json.NewDecoder(r.Body).Decode(&member)
	if err != nil {
		http.Error(w, "Invalid JSON", http.StatusBadRequest)
		return
	}

	// Validate input
	if member.FirstName == "" || member.LastName == "" || member.Email == "" || member.Cohort == "" {
		http.Error(w, "Missing required fields", http.StatusBadRequest)
		return
	}

	err = StoreMember(&member, db)
	if err != nil{
		http.Error(w, fmt.Sprintf("Failed to save data: %v", err), http.StatusInternalServerError)
		return
	}

	response := utils.Response{
		Message: "Member created and saved",
	}

	log.Printf("Member Created and saved successfully")
	utils.WriteJSONResponse(w, response, http.StatusOK)
}

// @Tags members
// @Summary Get all members
// @Description Get all fortify member
// @Accept json
// @Produce json
// @Param year query int false "Cohort Year"
// @Param cohort query string false "Cohort"
// @Param status query string false "Completion Status"
// @Success  200 {array} Member
// @Failure 400
// @Router /members [get]
func GetMembersHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB){
	
	cohortFilter := r.URL.Query().Get("cohort")
	statusFilter := r.URL.Query().Get("status")
	yearFilter := r.URL.Query().Get("year")
	
	members, err := GetMembers(cohortFilter, yearFilter, statusFilter, db)
	if err != nil{
		if _, ok := err.(*utils.ValidationError); ok{
			http.Error(w, "Invalid parameter", http.StatusBadRequest)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to fetch data: %v", err), http.StatusInternalServerError)
		return
	}
	if len(members) == 0{
		w.WriteHeader(http.StatusNoContent)
		return
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(members)
}

// @Tags members
// @Summary Update a member
// @Description Update a fortify member using member ID
// @Accept json
// @Param id path int true "Member ID"
// @Produce json
// @Success  200
// @Failure 400
// @Router /members/{id} [post]
func UpdateMemberHandler(w http.ResponseWriter, idStr string, r *http.Request, db *gorm.DB){
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}
	
	var inputMember Member
	if err := json.NewDecoder(r.Body).Decode(&inputMember); err != nil {
		http.Error(w, "Invalid JSON data", http.StatusBadRequest)
		return
	}
	member, err := UpdateMember(id, inputMember, db)
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			http.Error(w, fmt.Sprintf("Member with id = %d not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to update member: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully update member with id %d", id)
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(member)
}

// @Tags members
// @Summary Delete a member
// @Description Delete a fortify member using member ID
// @Param id path int true "Member ID"
// @Success  200
// @Failure 400
// @Router /members/{id} [delete]
func DeleteMemberHandler(w http.ResponseWriter, idStr string, r *http.Request, db *gorm.DB){
	
	id, err := strconv.Atoi(idStr)
	if err != nil {
		http.Error(w, "Invalid member ID", http.StatusBadRequest)
		return
	}
	
	err = DeleteMember(id, db)
	if err != nil{
		if errors.Is(err, gorm.ErrRecordNotFound){
			http.Error(w, fmt.Sprintf("Member with id = %d not found", id), http.StatusNotFound)
			return
		}
		http.Error(w, fmt.Sprintf("Failed to delete member: %v", err), http.StatusInternalServerError)
		return
	}

	log.Printf("Successfully deleted member with id %d", id)
	w.WriteHeader(http.StatusOK)
}