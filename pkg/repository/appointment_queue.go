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

import "gorm.io/gorm"

// AppointmentQueue ...
type AppointmentQueue struct {
	gorm.Model
	ID                 int              `gorm:"primaryKey" json:"id"`
	AppointmentID      int              `json:"appointmentId"`
	Appointment        Appointment      `json:"appointment"`
	QueueDestinationID int              `json:"queueDestinationId"`
	QueueDestination   QueueDestination `json:"queueDestination"`
	Count              int64            `json:"count"`
}

// Save ...
func (r *AppointmentQueue) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *AppointmentQueue) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// Update ...
func (r *AppointmentQueue) Update() (*AppointmentQueue, error) {
	err := DB.Save(&r).Error
	if err != nil {
		return nil, err
	}
	return r, nil
}

// FindByAppointment ...
func (r *AppointmentQueue) FindByAppointment(p PaginationInput, userID int) ([]AppointmentQueue, int64, error) {
	var result []AppointmentQueue
	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where("appointment_id = ?", userID).Preload("Appointment").Preload("QueueDestination").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// FindTodaysAppointments ...
func (r *AppointmentQueue) FindTodaysAppointments(appointmentID int) ([]AppointmentQueue, error) {
	var result []AppointmentQueue

	err := DB.Joins("left join user_queues on user_queues.queue_id = queues.id").Where("user_queues.appointment_id = ?", appointmentID).Find(&result).Error

	if err != nil {
		return result, err
	}

	return result, nil
}


// Delete ...
func (r *AppointmentQueue) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
