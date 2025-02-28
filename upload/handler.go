package upload

import (
	"fmt"
	"log"
	"io"
	"net/http"
	"os"
	"gorm.io/gorm"
	"fortifyApp/utils"
	"fortifyApp/members"
)

// @Tags file-upload
// @Summary Upload a file
// @Description Upload a CSV file that contains fortify members
// @Accept multipart/form-data
// @Param file formData file true "CSV file to upload"
// @Success 200
// @Failure 400
// @Router /upload [post]
func UploadFileHandler(w http.ResponseWriter, r *http.Request, db *gorm.DB) {

	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	cohort := r.URL.Query().Get("cohort")
	year := r.URL.Query().Get("year")

	if cohort == "" || year == ""{
		http.Error(w, "Missing Required Parameters", http.StatusBadRequest)
		return
	}

	// Limit the size of the incoming request to 10MB
	r.Body = http.MaxBytesReader(w, r.Body, 10<<20)

	if err := r.ParseMultipartForm(10 << 20); err != nil {
		http.Error(w, "File too big or invalid form", http.StatusBadRequest)
		return
	}

	// Get the file from the form
	file, header, err := r.FormFile("uploadFile")
	if err != nil {
		http.Error(w, "Could not read file", http.StatusBadRequest)
		return
	}
	defer file.Close()
	
	// Create a new file in the server's filesystem
	dst, err := os.Create(fmt.Sprintf("./documentUploads/%s", header.Filename))
	if err != nil {
		http.Error(w, "Unable to save file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	// Copy the file's content to the server file
	if _, err := io.Copy(dst, file); err != nil {
		http.Error(w, "Failed to save file", http.StatusInternalServerError)
		return
	}

	err = members.StoreMembers(header.Filename, db, year, cohort)
	if err != nil{
		// log.Fatalf("Failed to store in the database: %v", err)
		response := utils.Response{
			Error: fmt.Sprintf("Error: %s", err),
		}
		utils.WriteJSONResponse(w, response, http.StatusInternalServerError)
		// http.Error(w, "Failed to store data", http.StatusInternalServerError)
		return
	}

	// Respond with success
	response := utils.Response{
		Message: "File uploaded and stored successfully",
	}

	utils.WriteJSONResponse(w, response, http.StatusOK)
	// w.Write([]byte(header.Filename))
	log.Printf("File uploaded successfully: %s", header.Filename)
	// fmt.Fprintf(w, "File uploaded successfully: %s", header.Filename)
}