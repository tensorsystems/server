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

// PaymentWaiver ...
type PaymentWaiver struct {
	gorm.Model
	ID        int     `gorm:"primaryKey"`
	PaymentID int     `json:"paymentId" gorm:"unique"`
	Payment   Payment `json:"payment"`
	UserID    int     `json:"userID"`
	User      User    `json:"user"`
	PatientID int     `json:"patientId"`
	Patient   Patient `json:"patient"`
	Approved  *bool   `json:"approved"`
	Count     int64   `json:"count"`
}

// Save ...
func (r *PaymentWaiver) Save() error {
	return DB.Save(&r).Error
}

// BatchSave ...
func (r *PaymentWaiver) BatchSave(waivers []PaymentWaiver) error {
	return DB.Save(&waivers).Error
}

// Get ...
func (r *PaymentWaiver) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetCount ...
func (r *PaymentWaiver) GetApprovedCount() (int, error) {
	var count int64
	err := DB.Model(&r).Where("approved IS NULL").Count(&count).Error
	return int(count), err
}

// GetAll ...
func (r *PaymentWaiver) GetAll(p PaginationInput) ([]PaymentWaiver, int64, error) {
	var result []PaymentWaiver

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Preload("Patient").Preload("User").Preload("Payment.Billing").Order("id DESC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Update ...
func (r *PaymentWaiver) Update() error {
	return DB.Updates(&r).Error
}

// ApproveWaiver ...
func (r *PaymentWaiver) ApproveWaiver(id int, approve bool) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Update payment waiver
		if err := tx.Where("id = ?", id).Preload("Payment").Take(&r).Error; err != nil {
			return err
		}

		r.Approved = &approve
		if err := tx.Updates(&r).Error; err != nil {
			return err
		}

		// Update payment status
		payment := r.Payment
		payment.Status = PaidPaymentStatus
		if err := tx.Updates(&payment).Error; err != nil {
			return err
		}

		return nil
	})
}

// Delete ...
func (r *PaymentWaiver) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
