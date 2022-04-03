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

// ChatMute ...
type ChatMute struct {
	gorm.Model
	ID     int `gorm:"primaryKey"`
	UserID int `json:"userId"`
	ChatID int `json:"chatId"`
}

// Save ...
func (r *ChatMute) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *ChatMute) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// Update ...
func (r *ChatMute) Update() error {
	return DB.Updates(&r).Error
}

// Delete ...
func (r *ChatMute) Delete(userID int, chatID int) error {
	return DB.Where("user_id = ? AND chat_id = ?", userID, chatID).Delete(&r).Error
}
