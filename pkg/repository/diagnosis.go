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

// Diagnosis ...
type Diagnosis struct {
	ID                     int     `gorm:"primaryKey" json:"id"`
	CategoryCode           *string `json:"categoryCode"`
	DiagnosisCode          *string `json:"diagnosisCode"`
	FullCode               *string `json:"fullCode"`
	AbbreviatedDescription *string `json:"abbreviatedDescription"`
	FullDescription        string  `json:"fullDescription"`
	CategoryTitle          *string `json:"categoryTitle"`
	Active                 bool    `json:"active"`
	Document               string  `gorm:"type:tsvector"`
	Count                  int64   `json:"count"`
}

// Save ...
func (r *Diagnosis) Save() error {
	err := DB.Create(&r).Error
	if err != nil {
		return err
	}

	return nil
}

// GetAll ...
func (r *Diagnosis) GetAll(p PaginationInput, searchTerm *string) ([]Diagnosis, int64, error) {
	var result []Diagnosis

	dbOp := DB.Scopes(Paginate(&p))

	if searchTerm != nil && len(*searchTerm) > 0 {
		dbOp.Where("document @@ plainto_tsquery(?)", *searchTerm)
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
func (r *Diagnosis) GetFavorites(p PaginationInput, searchTerm *string, userId int) ([]Diagnosis, int64, error) {
	var result []Diagnosis
	var count int64
	var err error

	var favoriteIds []int
	var entity FavoriteDiagnosis
	favoriteDiagnosis, _ := entity.GetByUser(userId)
	for _, e := range favoriteDiagnosis {
		favoriteIds = append(favoriteIds, e.DiagnosisID)
	}

	if len(favoriteIds) > 0 {
		var favorites []Diagnosis

		favoritesDb := DB.Where("id IN ?", favoriteIds)
		if searchTerm != nil && len(*searchTerm) > 0 {
			favoritesDb.Where("full_description ILIKE ?", "%"+*searchTerm+"%")
		}
		favoritesDb.Find(&favorites)

		result = append(result, favorites...)

		var nonFavorites []Diagnosis
		nonFavoritesDb := DB.Not(favoriteIds).Scopes(Paginate(&p))
		if searchTerm != nil && len(*searchTerm) > 0 {
			nonFavoritesDb.Where("full_description ILIKE ?", "%"+*searchTerm+"%")
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
func (r *Diagnosis) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByTitle ...
func (r *Diagnosis) GetByTitle(title string) error {
	return DB.Where("title = ?", title).Take(&r).Error
}

// Update ...
func (r *Diagnosis) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *Diagnosis) Delete(ID int) error {
	return DB.Where("id = ?", ID).Delete(&r).Error
}
