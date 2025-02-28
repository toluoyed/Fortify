package members

import (
	"fmt"
	"os"
	"strconv"
	"strings"
	"errors"
	"encoding/csv"
	"log"
	"gorm.io/gorm"
	"gorm.io/gorm/clause"
	"fortifyApp/utils"
)

func StoreMembers(fileName string, db *gorm.DB, year string, cohort string) error{

	file, err := os.Open(fmt.Sprintf("./documentUploads/%s", fileName))

	if err != nil{
		log.Fatal("Failed to open file: ", err)
		return fmt.Errorf("Failed to open file: %v", err)
	}

	defer file.Close()

	reader := csv.NewReader(file)
	records, err := reader.ReadAll()
	
	if err != nil {
		log.Fatalf("Failed to read CSV file: %v", err)
		return fmt.Errorf("Failed to read CSV file: %v", err)
	}

	for i, record := range records {
		if i == 0 {
			// Skip the header row
			continue
		}
		year, _ := strconv.Atoi(year)
		member := Member{FirstName: record[0], LastName: record[1], Email: record[2], PhoneNumber: &record[3], Cohort: cohort, Year: year}

		
		
		err = db.Clauses(clause.OnConflict{DoNothing: true}).Create(&member).Error; 
		if err != nil{
			if errors.Is(err, gorm.ErrDuplicatedKey) {
				log.Printf("Duplicate record detected for record: %v", record)
			} else {
				log.Printf("Failed to insert record: %v, error: %v", record, err)
			}
			return err
		}
		
	}

	return nil
}

func StoreMember(member *Member, db *gorm.DB) error{
	result := db.Create(member)
	return result.Error
}

func GetMembers(cohortFilter string, yearFilter string, statusFilter string, db *gorm.DB) ([]Member,error){
	var members []Member
	query := db.Model(&Member{})

	if cohortFilter != "" {
		cohorts := strings.Split(cohortFilter, ",")
		query = query.Where("cohort IN ?", cohorts)
	}
	if yearFilter != ""{

		year, err := strconv.Atoi(yearFilter)
		if err != nil {
			log.Printf("Error converting Year string to Int: %v", err)
			return members,&utils.ValidationError{Parameter: "year"}
		}

		query = query.Where("year = ?", year)
	}
	if statusFilter != "" {
		query = query.Where("status = ?", statusFilter)
	}
	if cohortFilter == "" && statusFilter == "" && yearFilter == ""{
		result := db.Find(&members)
		return members, result.Error
	}
	result := query.Find(&members)
	return members, result.Error

	
}

func UpdateMember(id int, requestMember Member, db *gorm.DB) (Member,error){
	
	var member Member
	result := db.First(&member, id)
	if result.Error != nil {
		log.Printf("Error finding member with id = %d: %v", id, result.Error)
		return member,result.Error
	}

	if requestMember.FirstName != "" {
		member.FirstName = requestMember.FirstName
	}
	if requestMember.LastName != "" {
		member.LastName = requestMember.LastName
	}
	if requestMember.Email != "" {
		member.Email = requestMember.Email
	}
	if requestMember.PhoneNumber != nil {
		member.PhoneNumber = requestMember.PhoneNumber
	}
	if requestMember.Cohort != "" {
		member.Cohort = requestMember.Cohort
	}
	if requestMember.Year != 0 {
		member.Year = requestMember.Year
	}
	if requestMember.Session1 != nil {
        *member.Session1 = *requestMember.Session1
    }
    if requestMember.Session2 != nil {
        *member.Session2 = *requestMember.Session2
    }
    if requestMember.Session3 != nil {
        *member.Session3 = *requestMember.Session3
    }
    if requestMember.Session4 != nil {
        *member.Session4 = *requestMember.Session4
    }

	return member, db.Save(&member).Error
	
}

func DeleteMember(id int, db *gorm.DB) error{
	var member Member
	result := db.First(&member, id)
	if result.Error != nil {
		log.Printf("Error deleting member with id = %d: %v", id, result.Error)
		return result.Error
	}
	return db.Delete(&member).Error
}

