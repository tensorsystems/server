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

package controller

import (
	"github.com/Tensor-Systems/tensoremr-server/pkg/repository"
	"github.com/gin-gonic/gin"
)

// GetUserTypes ...
func GetUserTypes(c *gin.Context) {
	var entity repository.UserType

	paginationInput := repository.PaginationInput{
		Page: 0,
		Size: 20,
	}

	userTypes, _, err := entity.GetAll(paginationInput)

	if err != nil {
		c.JSON(400, gin.H{
			"msg": err,
		})
		c.Abort()

		return
	}

	c.JSON(200, userTypes)
}
