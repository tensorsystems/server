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

package models

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
func (r *User) CheckPassword(m *User, providedPassword string) error {
	err := bcrypt.CompareHashAndPassword([]byte(m.Password), []byte(providedPassword))
	if err != nil {
		return err
	}

	return nil
}
