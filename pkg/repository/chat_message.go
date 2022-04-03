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

// ChatMessage ...
type ChatMessage struct {
	gorm.Model
	ID     int    `gorm:"primaryKey"`
	Body   string `json:"body"`
	ChatID int    `json:"chatId"`
	UserID int    `json:"userId"`
}

// Save ...
func (r *ChatMessage) Save() error {
	return DB.Create(&r).Error
}

// Get ...
func (r *ChatMessage) Get(ID int) error {
	return DB.Where("id = ?", ID).Take(&r).Error
}

// GetByChatID ...
func (r *ChatMessage) GetByChatID(ID int) ([]*ChatMessage, error) {
	var messages []*ChatMessage
	err := DB.Where("chat_id = ?", ID).Order("created_at asc").Find(&messages).Error
	return messages, err
}

// Update ...
func (r *ChatMessage) Update() error {
	return DB.Updates(&r).Error
}
