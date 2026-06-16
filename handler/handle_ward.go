package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"ivr_ataljanseva/db/repository"
	"ivr_ataljanseva/models"
)

type WardHandler struct {
	politicalRepo *repository.PoliticalUserRepository
	citizenRepo   *repository.CitizenRepository
}
func NewWardHandler(
	politicalRepo *repository.PoliticalUserRepository,
	citizenRepo *repository.CitizenRepository,
) *WardHandler {
	return &WardHandler{
		politicalRepo: politicalRepo,
		citizenRepo:   citizenRepo,
	}
}

func (h *WardHandler) ResolveWard(c *gin.Context) {

	var req models.ResolveWardRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	matches, err := h.politicalRepo.FindMatchingWards(
		c.Request.Context(),
		req.Pincode,
		req.WardInput,
	)

	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	switch {

		case len(matches) == 0:

			c.JSON(http.StatusOK, models.ResolveWardResponse{
				Status: "not_found",
			})

		case len(matches) == 1:

			match := matches[0]

			err := h.citizenRepo.UpsertCitizen(
				c.Request.Context(),
				req.Phone,
				req.Language,
				req.Pincode,
				match.Ward,
				match.NagarsevakID,
			)

			if err != nil {
				c.JSON(http.StatusInternalServerError, gin.H{
					"error": err.Error(),
				})
				return
			}

			c.JSON(http.StatusOK, models.ResolveWardResponse{
				Status:         "single_ward",
				Ward:           match.Ward,
				NagarsevakID:   match.NagarsevakID,
				NagarsevakName: match.NagarsevakName,
			})

		case len(matches) <= 4:

			var wards []string

			for _, m := range matches {
				wards = append(wards, m.Ward)
			}

			c.JSON(http.StatusOK, models.ResolveWardResponse{
				Status: "choose_ward",
				Wards:  wards,
			})

		default:

			c.JSON(http.StatusOK, models.ResolveWardResponse{
				Status: "too_many_matches",
			})

	}
}