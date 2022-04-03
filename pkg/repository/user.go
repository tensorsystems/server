/*
  Copyright 2021 Kidus Tiliksew

  This file is part of Tensor EMR.

  Tensor EMR is free software: you can redistribute it and/or modify
  it under the terms of the version 2 of GNU General Public License as published by
  the Free Software Foundation.

  Tensor EMR is distributed in the hope that it will be useful,
  but WITHOUT ANY WARRANTY; without even the implied warranty of
  MERCHANTABILITY or FITNESS FOR A PARTICULAR PURPOSE.  See the
  GNU General Public License for more details.

  You should have received a copy of the GNU General Public License
  along with this program.  If not, see <http://www.gnu.org/licenses/>.
*/

package repository

import (
	"time"

	"golang.org/x/crypto/bcrypt"
	"gorm.io/datatypes"
	"gorm.io/gorm"
)

// User ...
type User struct {
	gorm.Model
	ID int `gorm:"primaryKey"`

	FirstName string `json:"firstName"`
	LastName  string `json:"lastName"`

	Email    string `json:"email"`
	Password string `json:"password"`

	OldUserName string `json:"oldUserName"`

	UserTypes []UserType `gorm:"many2many:user_type_roles;" json:"userTypes"`

	Appointments []Appointment `json:"appointments"`

	Active bool `json:"active"`

	SignatureID *int  `json:"signatureId"`
	Signature   *File `json:"signature"`

	ProfilePicID *int  `json:"profilePicId"`
	ProfilePic   *File `json:"profilePic"`

	// Confirm
	ConfirmSelector string
	ConfirmVerifier string
	Confirmed       bool

	// Lock
	AttemptCount int
	LastAttempt  *time.Time
	Locked       *time.Time

	// Recover
	RecoverSelector    string
	RecoverVerifier    string
	RecoverTokenExpiry *time.Time

	// OAuth2
	OAuth2UID          string
	OAuth2Provider     string
	OAuth2AccessToken  string
	OAuth2RefreshToken string
	OAuth2Expiry       *time.Time

	// 2fa
	TOTPSecretKey      string
	SMSPhoneNumber     string
	SMSSeedPhoneNumber string
	RecoveryCodes      string

	Document string `gorm:"type:tsvector"`
	Count    int64  `json:"count"`
}

// Seed ...
func (r *User) Seed() {
	var userType UserType
	userType.GetByTitle("Admin")

	var user User
	user.FirstName = "Admin"
	user.LastName = "Admin"
	user.Email = "info@tensorsystems.net"
	user.UserTypes = append(user.UserTypes, userType)
	user.Password = "changeme"
	user.Active = true
	user.HashPassword()

	DB.Create(&user)
}

// AfterCreate ...
func (r *User) AfterCreate(tx *gorm.DB) error {
	var user User
	if err := tx.Where("id = ?", r.ID).Preload("UserTypes").Error; err != nil {
		return err
	}

	isPhysician := false
	for _, e := range r.UserTypes {
		if e.Title == "Physician" {
			isPhysician = true
		}
	}

	if isPhysician {
		var patientEncounterLimit PatientEncounterLimit
		patientEncounterLimit.MondayLimit = 150
		patientEncounterLimit.TuesdayLimit = 150
		patientEncounterLimit.WednesdayLimit = 150
		patientEncounterLimit.ThursdayLimit = 150
		patientEncounterLimit.FridayLimit = 150
		patientEncounterLimit.SaturdayLimit = 150
		patientEncounterLimit.SundayLimit = 150
		patientEncounterLimit.Overbook = 5

		if err := tx.Create(&patientEncounterLimit).Error; err != nil {
			return err
		}

		queue := datatypes.JSON([]byte("[" + "]"))

		var patientQueue PatientQueue
		patientQueue.QueueName = "Dr. " + user.FirstName + " " + user.LastName
		patientQueue.Queue = queue
		patientQueue.QueueType = "USER"

		if err := tx.Create(&patientQueue).Error; err != nil {
			return err
		}

		var queueSubscription QueueSubscription
		queueSubscription.UserID = user.ID
		queueSubscription.Subscriptions = append(queueSubscription.Subscriptions, patientQueue)

		if err := tx.Create(&queueSubscription).Error; err != nil {
			return err
		}
	}

	return nil
}

// HashPassword encrypts user password
func (r *User) HashPassword() error {
	bytes, err := bcrypt.GenerateFromPassword([]byte(r.Password), 14)
	if err != nil {
		return err
	}

	r.Password = string(bytes)
	return nil
}

// HashPassword encrypts user password
func (r *User) CheckPasswordEquality(password1 string, password2 string) bool {
	bytes1, err := bcrypt.GenerateFromPassword([]byte(password1), 14)
	if err != nil {
		return false
	}

	bytes2, err := bcrypt.GenerateFromPassword([]byte(password2), 14)
	if err != nil {
		return false
	}

	return string(bytes1) == string(bytes2)
}

// CheckPassword checks user password
func (r *User) CheckPassword(providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(r.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}

// Save ...
func (r *User) Save(userTypes []UserType) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		if err := tx.Create(&r).Error; err != nil {
			return err
		}

		if userTypes != nil {
			tx.Model(&r).Association("UserTypes").Replace(&userTypes)

			isPhysician := false
			for _, e := range userTypes {
				if e.Title == "Physician" {
					isPhysician = true
				}
			}

			if isPhysician {

			}
		}

		return nil
	})
}

// Search ...
func (r *User) SearchPhysicians(searchTerm string) ([]*User, error) {
	var result []*User

	var userType UserType
	if err := userType.GetByTitle("Physician"); err != nil {
		return nil, err
	}

	if len(searchTerm) > 0 {
		DB.Model(&userType).Where("first_name ILIKE ?", "%"+searchTerm+"%").Or("last_name ILIKE ?", "%"+searchTerm+"%").Association("Users").Find(&result)
	}

	return result, nil
}

// GetAll ...
func (r *User) GetAll(p PaginationInput) ([]User, int64, error) {
	var result []User

	err := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Order("id ASC").Preload("UserTypes").Preload("Signature").Preload("ProfilePic").Find(&result).Error
	if err != nil {
		return result, 0, err
	}

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// Search ...
func (r *User) Search(p PaginationInput, filter *User, searchTerm *string) ([]User, int64, error) {
	var result []User

	tx := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Order("id ASC").Preload("UserTypes").Preload("Signature").Preload("ProfilePic")

	if searchTerm != nil {

		tx.Where("document @@ plainto_tsquery(?)", *searchTerm)
	}

	err := tx.Find(&result).Error

	if err != nil {
		return result, 0, err
	}

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	return result, count, err
}

// Get ...
func (r *User) Get(ID int) error {
	return DB.Where("id = ?", ID).Preload("Signature").Preload("ProfilePic").Preload("UserTypes").Take(&r).Error
}

// GetByEmail ...
func (r *User) GetByEmail(email string) error {
	return DB.Where("email = ?", email).Preload("UserTypes").Take(&r).Error
}

// GetByOldUserName ...
func (r *User) GetByOldUserName(userName string) error {
	return DB.Where("old_user_name ILIKE ?", userName).Preload("UserTypes").Take(&r).Error
}

// CheckIfUserLegacy ...
func (r *User) CheckIfUserLegacy(oldUserName string) error {
	return DB.Where("old_user_name ILIKE ?", oldUserName).Where("email != ''").Take(&r).Error
}

// GetByUserType ...
func (r *User) GetByUserType(userTypeID int) (users []User, err error) {
	err = DB.Where("user_type_id = ?", userTypeID).Find(&users).Error
	return
}

// GetByUserTypeTitle ...
func (r *User) GetByUserTypeTitle(userTypeTitle string) ([]*User, error) {
	var userType UserType
	err := DB.Model(&UserType{}).Where("title = ?", userTypeTitle).Preload("Users").Take(&userType).Error
	if err != nil {
		return nil, err
	}

	var result []*User
	for i, v := range userType.Users {
		if v.Active {
			result = append(result, &userType.Users[i])
		}
	}

	return result, err
}

// Update ...
func (r *User) Update(userTypes []UserType) error {
	err := DB.Updates(&r).Error
	if err != nil {
		return err
	}

	if err := DB.Select("active").Where("id = ?", r.ID).Updates(User{Active: r.Active}).Error; err != nil {
		return err
	}

	if userTypes != nil {
		DB.Model(&r).Association("UserTypes").Replace(&userTypes)
	}

	return nil
}

// Ping ...
func (r *User) Ping() error {
	db, err := DB.DB()
	if err != nil {
		return err
	}

	return db.Ping()
}
