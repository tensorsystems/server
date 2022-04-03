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
	"database/sql/driver"

	"gorm.io/gorm"
)

// PaymentStatus ...
type PaymentStatus string

// Payment statuses ...
const (
	PaidPaymentStatus            PaymentStatus = "PAID"
	NotPaidPaymentStatus         PaymentStatus = "NOTPAID"
	WaiverRequestedPaymentStatus PaymentStatus = "PAYMENT_WAIVER_REQUESTED"
)

// Scan ...
func (p *PaymentStatus) Scan(value interface{}) error {
	*p = PaymentStatus(value.(string))
	return nil
}

// Value ...
func (p PaymentStatus) Value() (driver.Value, error) {
	return string(p), nil
}

// Payment ...
type Payment struct {
	gorm.Model
	ID        int           `gorm:"primaryKey"`
	InvoiceNo string        `json:"invoiceNo"`
	Status    PaymentStatus `json:"status" sql:"payment_status"`
	BillingID int           `json:"billingId"`
	Billing   Billing       `json:"billing"`
}

// Save ...
func (r *Payment) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *Payment) Get(ID int) error {
	err := DB.Where("id = ?", ID).Take(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// Get ...
func (r *Payment) GetByIds(ids []int) ([]Payment, error) {
	var result []Payment
	err := DB.Where("id IN ?", ids).Find(&result).Error
	if err != nil {
		return nil, err
	}
	return result, nil
}

// Update ...
func (r *Payment) Update() error {
	return DB.Updates(&r).Error
}

// BatchUpdate ...
func (r *Payment) BatchUpdate(ids []int, e Payment) error {
	return DB.Model(&r).Where("id IN ?", ids).Updates(&e).Error
}

// RequestWaiver ...
func (r *Payment) RequestWaiver(paymentID int, patientID int, userID int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Update payment
		r.ID = paymentID
		r.Status = WaiverRequestedPaymentStatus
		if err := tx.Updates(&r).Error; err != nil {
			return err
		}

		// Save payment waiver
		var paymentWaiver PaymentWaiver
		paymentWaiver.PatientID = patientID
		paymentWaiver.PaymentID = paymentID
		paymentWaiver.UserID = userID

		if err := tx.Save(&paymentWaiver).Error; err != nil {
			return err
		}

		return nil
	})
}

// RequestWaiverBatch ...
func (r *Payment) RequestWaiverBatch(paymentIds []int, patientId int, userId int) error {
	return DB.Transaction(func(tx *gorm.DB) error {
		// Update payments
		if err := tx.Where("id IN ?", paymentIds).Updates(&Payment{Status: WaiverRequestedPaymentStatus}).Error; err != nil {
			return err
		}

		// Save payment waivers
		var paymentWaivers []PaymentWaiver
		for i := range paymentIds {
			waiver := PaymentWaiver{
				PatientID: patientId,
				PaymentID: paymentIds[i],
				UserID:    userId,
			}

			paymentWaivers = append(paymentWaivers, waiver)
		}

		if err := tx.Save(&paymentWaivers).Error; err != nil {
			return err
		}

		return nil
	})
}

// Delete ...
func (r *Payment) Delete(ID int) error {
	var e Payment
	err := DB.Where("id = ?", ID).Delete(&e).Error
	return err
}
