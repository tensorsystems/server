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
	"fmt"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/tensoremr/server/pkg/repository"
)

// RecreateOpthalmologyExam ...
func RecreateOpthalmologyExam(c *gin.Context) {
	patientChartIdQuery := c.Query("patientChartId")

	patientChartId, err := strconv.Atoi(patientChartIdQuery)
	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}

	var entity repository.OpthalmologyExam

	if err := entity.Recreate(patientChartId); err != nil {
		c.String(http.StatusNotFound, fmt.Sprintf("error: %s", err))
	}

	c.JSON(200, entity)
}

// ClearPatientsRecord ...
func ClearPatientsRecord(c *gin.Context) {
	var entity repository.Patient

	err := entity.Clean()

	if err != nil {
		c.String(http.StatusInternalServerError, fmt.Sprintf("error: %s", err))
	}

	c.JSON(200, "Success")
}
