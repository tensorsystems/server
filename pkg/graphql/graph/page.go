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

package graph

import (
	"math"

	"github.com/tensoremr/server/pkg/graphql/graph/model"
	"github.com/tensoremr/server/pkg/repository"
)

// ConvertEntityToConnection ...
func GetPageInfo[R any](entities []R, count int64, page repository.PaginationInput) (*model.PageInfo, int) {
	if len(entities) == 0 {
		return &model.PageInfo{}, 0
	}

	pageInfo := model.PageInfo{}
	totalPages := math.Ceil(float64(count) / float64(page.Size))
	pageInfo.TotalPages = int(totalPages)

	return &pageInfo, int(count)
}
