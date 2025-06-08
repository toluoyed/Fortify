package users

import (
	"log"
	"gorm.io/gorm"
)

func CreateUser(user *User, db *gorm.DB) error{
	result := db.Create(user)
	return result.Error
}

func GetUsers(db *gorm.DB) ([]User,error){
	var users []User
	result := db.Find(&users)
	
	return users, result.Error
}

func UpdateUser(id int, requestUser User, db *gorm.DB) (User,error){
	
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Printf("Error finding user with id = %d: %v", id, result.Error)
		return user,result.Error
	}

	updates := map[string]interface{}{}

    if requestUser.FirstName != "" {
        updates["first_name"] = requestUser.FirstName
    }
	if requestUser.LastName != "" {
        updates["last_name"] = requestUser.LastName
    }
    if requestUser.Email != "" {
        updates["email"] = requestUser.Email
    }
	if requestUser.Password != "" {
        updates["password"] = requestUser.Password
    }
	role := requestUser.Role
	if role != "" && role == "SUPERUSER" || role == "USER"{
		updates["role"] = Role(role)
	}


	if len(updates) == 0 {
        return User{}, nil
    }

    // Update only the provided fields
    updateResult := db.Model(&User{}).Where("id = ?", id).Updates(updates)
    if updateResult.Error != nil {
        log.Printf("Error updating user with id = %d: %v", id, updateResult.Error)
        return User{}, updateResult.Error
    }

	// Retrieve the updated member
    var updatedUser User
    db.First(&updatedUser, id)

    return updatedUser, nil
	
}

func DeleteUser(id int, db *gorm.DB) error{
	var user User
	result := db.First(&user, id)
	if result.Error != nil {
		log.Printf("Error deleting user with id = %d: %v", id, result.Error)
		return result.Error
	}
	return db.Delete(&user).Error
}