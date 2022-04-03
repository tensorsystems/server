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

// SurgicalProcedureType ...
type SurgicalProcedureType struct {
	gorm.Model
	ID       int       `gorm:"primaryKey" json:"id"`
	Title    string    `json:"title"`
	Active   bool      `json:"active"`
	Billings []Billing `json:"billings" gorm:"many2many:surgical_procedure_type_billings"`
	Supplies []Supply  `json:"supplies" gorm:"many2many:surgical_procedure_supplies"`
	Count    int64     `json:"count"`
}

// Save ...
func (r *SurgicalProcedureType) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *SurgicalProcedureType) GetAll(p PaginationInput, searchTerm *string) ([]SurgicalProcedureType, int64, error) {
	var result []SurgicalProcedureType

	dbOp := DB.Scopes(Paginate(&p)).Select("*, count(*) OVER() AS count")

	if searchTerm != nil {
		dbOp.Where("title ILIKE ?", "%"+*searchTerm+"%")
	}

	dbOp.Order("id ASC").Preload("Billings").Preload("Supplies.Billings").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// Get ...
func (r *SurgicalProcedureType) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *SurgicalProcedureType) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *SurgicalProcedureType) Update() error {
	return DB.Transaction(func(tx *gorm.DB) error {
		tx.Model(&r).Association("Billings").Replace(&r.Billings)
		tx.Model(&r).Association("Supplies").Replace(&r.Supplies)

		return tx.Updates(&r).Error
	})
}

// Delete ...
func (r *SurgicalProcedureType) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
