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
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
)

// GetPatientQueues ...
func GetPatientQueues(c *gin.Context) {
	var repo repository.PatientQueue

	patientQueues, err := repo.GetAll()

	if err != nil {
		c.JSON(400, gin.H{
			"msg": err,
		})
		c.Abort()

		return
	}

	var result []*model.PatientQueueWithAppointment
	var appointmentRepo repository.Appointment

	for _, patientQueue := range patientQueues {
		var ids []int

		if err := json.Unmarshal([]byte(patientQueue.Queue.String()), &ids); err != nil {
			c.JSON(400, gin.H{
				"msg": err,
			})
			c.Abort()

			return
		}

		if len(ids) == 0 {
			continue
		}

		page := repository.PaginationInput{Page: 0, Size: 1000}

		appointments, _, _ := appointmentRepo.GetByIds(ids, page)

		var orderedAppointments []*repository.Appointment

		for _, id := range ids {
			for _, appointment := range appointments {
				if appointment.ID == id {
					a := appointment
					orderedAppointments = append(orderedAppointments, &a)
				}
			}
		}

		result = append(result, &model.PatientQueueWithAppointment{
			ID:        int(patientQueue.ID),
			QueueName: patientQueue.QueueName,
			QueueType: patientQueue.QueueType,
			Queue:     orderedAppointments,
		})
	}

	c.JSON(200, result)
}
