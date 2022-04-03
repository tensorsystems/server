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

// ChiefComplaintType ...
type ChiefComplaintType struct {
	gorm.Model
	ID     int    `gorm:"primaryKey" json:"id"`
	Title  string `json:"title"`
	Count  int64  `json:"count"`
	Active bool   `json:"active"`
}

// Save ...
func (r *ChiefComplaintType) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *ChiefComplaintType) GetAll(p PaginationInput, searchTerm *string) ([]ChiefComplaintType, int64, error) {
	var result []ChiefComplaintType

	dbOp := DB.Scopes(Paginate(&p))

	if searchTerm != nil {
		dbOp.Where("title ILIKE ?", "%"+*searchTerm+"%")
	}

	dbOp.Order("id ASC").Find(&result)

	var count int64
	if len(result) > 0 {
		count = result[0].Count
	}

	if dbOp.Error != nil {
		return result, 0, dbOp.Error
	}

	return result, count, dbOp.Error
}

// GetFavorites ...
func (r *ChiefComplaintType) GetFavorites(p PaginationInput, searchTerm *string, userId int) ([]ChiefComplaintType, int64, error) {
	var result []ChiefComplaintType
	var count int64
	var err error

	var favoriteIds []int
	var entity FavoriteChiefComplaint
	favoriteChiefComplaints, _ := entity.GetByUser(userId)
	for _, e := range favoriteChiefComplaints {
		favoriteIds = append(favoriteIds, e.ChiefComplaintTypeID)
	}

	if len(favoriteIds) > 0 {
		var favorites []ChiefComplaintType

		favoritesDb := DB.Where("id IN ?", favoriteIds)
		if searchTerm != nil {
			favoritesDb.Where("title ILIKE ?", "%"+*searchTerm+"%")
		}
		favoritesDb.Find(&favorites)

		result = append(result, favorites...)

		var nonFavorites []ChiefComplaintType
		nonFavoritesDb := DB.Not(favoriteIds).Scopes(Paginate(&p))
		if searchTerm != nil {
			nonFavoritesDb.Where("title ILIKE ?", "%"+*searchTerm+"%")
		}
		nonFavoritesDb.Find(&nonFavorites)

		result = append(result, nonFavorites...)

		if len(nonFavorites) > 0 {
			count = nonFavorites[0].Count + int64(len(favoriteIds))
		}
	} else {
		return r.GetAll(p, searchTerm)
	}

	return result, count, err
}

// Get ...
func (r *ChiefComplaintType) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *ChiefComplaintType) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *ChiefComplaintType) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *ChiefComplaintType) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
