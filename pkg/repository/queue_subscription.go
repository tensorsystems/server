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
	"gorm.io/gorm"
)

// QueueSubscription ....
type QueueSubscription struct {
	gorm.Model
	ID            int            `gorm:"primaryKey" json:"id"`
	UserID        int            `json:"userId" gorm:"uniqueIndex"`
	User          User           `json:"user"`
	Subscriptions []PatientQueue `json:"subscriptions" gorm:"many2many:user_queue_subscriptions;"`
}

// Save
func (r *QueueSubscription) Save() error {
	return DB.Create(&r).Error
}

// GetByUserId ...
func (r *QueueSubscription) GetByUserId(userID int) error {
	return DB.Where("user_id = ?", userID).Preload("Subscriptions").Find(&r).Error
}

// Subscribe ...
func (r *QueueSubscription) Subscribe(userId int, patientQueueId int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := DB.Where("id = ?", userId).Take(&user).Error; err != nil {
			return err
		}

		if err := DB.Where("user_id = ?", user.ID).Take(&r).Error; err != nil {
			r.UserID = user.ID

			if err := DB.Create(&r).Error; err != nil {
				return err
			}
		}

		var patientQueue PatientQueue
		if err := DB.Where("id = ?", patientQueueId).Take(&patientQueue).Error; err != nil {
			return err
		}

		DB.Model(&r).Association("Subscriptions").Append(&patientQueue)

		return nil
	})
}

// Unsubscribe ...
func (r *QueueSubscription) Unsubscribe(userId int, patientQueueId int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		var user User
		if err := DB.Where("id = ?", userId).Take(&user).Error; err != nil {
			return err
		}

		if err := DB.Where("user_id = ?", user.ID).Take(&r).Error; err != nil {
			return err
		}

		var patientQueue PatientQueue
		if err := DB.Where("id = ?", patientQueueId).Take(&patientQueue).Error; err != nil {
			return err
		}

		DB.Model(&r).Association("Subscriptions").Delete(&patientQueue)
		return nil
	})
}
