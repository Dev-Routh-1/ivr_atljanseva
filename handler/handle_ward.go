package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ivr_ataljanseva/asterisk"
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

			c.JSON(http.StatusOK, gin.H{
				"status": "not_found",
				"_channel_vars": asterisk.FromResolveWard(
					"not_found", req.Phone, req.Language, req.Pincode,
					"", "", "",
				),
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

			c.JSON(http.StatusOK, gin.H{
				"status":          "single_ward",
				"ward":            match.Ward,
				"nagarsevak_id":   match.NagarsevakID.String(),
				"nagarsevak_name": match.NagarsevakName,
				"_channel_vars": asterisk.FromResolveWard(
					"single_ward", req.Phone, req.Language, req.Pincode,
					match.Ward, match.NagarsevakID.String(), match.NagarsevakName,
				),
			})

		case len(matches) <= 4:

			var wards []string

			for _, m := range matches {
				wards = append(wards, m.Ward)
			}

			c.JSON(http.StatusOK, gin.H{
				"status": "choose_ward",
				"wards":  wards,
				"_channel_vars": asterisk.FromResolveWard(
					"choose_ward", req.Phone, req.Language, req.Pincode,
					"", "", "",
				),
			})

		default:

			c.JSON(http.StatusOK, gin.H{
				"status": "too_many_matches",
				"_channel_vars": asterisk.FromResolveWard(
					"too_many_matches", req.Phone, req.Language, req.Pincode,
					"", "", "",
				),
			})

	}
}