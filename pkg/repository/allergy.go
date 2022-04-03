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

// Allergy ...
type Allergy struct {
	gorm.Model
	ID               int    `gorm:"primaryKey" json:"id"`
	Title            string `json:"title"`
	IssueSeverity    string `json:"issueSeverity"`
	IssueReaction    string `json:"issueReaction"`
	IssueOutcome     string `json:"issueOutcome"`
	IssueOccurrence  string `json:"issueOccurrence"`
	PatientHistoryID int    `json:"patientHistoryId"`
	Count            int64  `json:"count"`
}

// Save ...
func (r *Allergy) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *Allergy) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// Update ...
func (r *Allergy) Update() error {
	return DB.Updates(&r).Error
}

// GetAll ...
func (r *Allergy) GetAll(p PaginationInput, filter *Allergy) ([]Allergy, int64, error) {
	var result []Allergy

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count").Where(filter).Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Delete ...
func (r *Allergy) Delete(ID int) error {
	err := DB.Where("id = ?", ID).Delete(&r).Error
	return err
}
