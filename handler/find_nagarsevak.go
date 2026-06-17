package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"

	"ivr_ataljanseva/asterisk"
	"ivr_ataljanseva/db/repository"
	"ivr_ataljanseva/models"
)

type NagarsevakHandler struct {
	politicalRepo *repository.PoliticalUserRepository
	citizenRepo   *repository.CitizenRepository
}

func NewNagarsevakHandler(
	politicalRepo *repository.PoliticalUserRepository,
	citizenRepo *repository.CitizenRepository,
) *NagarsevakHandler {
	return &NagarsevakHandler{
		politicalRepo: politicalRepo,
		citizenRepo:   citizenRepo,
	}
}

func (h *NagarsevakHandler) FindNagarsevak(c *gin.Context) {
	var req models.NagarsevakLookupRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	nagarsevaks, err := h.politicalRepo.FindNagarsevaks(
		c.Request.Context(),
		req.Pincode,
		req.Ward,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	count := len(nagarsevaks)

	switch {
	case count == 0:
		c.JSON(http.StatusOK, gin.H{
			"status": "not_found",
			"_channel_vars": asterisk.FromNagarsevak(
				"not_found", req.PhoneNumber, req.Language,
				req.Pincode, req.Ward, "", "",
			),
		})

	case count == 1:
		ns := nagarsevaks[0]

		err := h.citizenRepo.UpsertCitizen(
			c.Request.Context(),
			req.PhoneNumber,
			req.Language,
			req.Pincode,
			req.Ward,
			ns.ID,
		)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"status":     "single",
			"auto_saved": true,
			"nagarsevak": map[string]string{
				"id":   ns.ID.String(),
				"name": ns.Name,
			},
			"_channel_vars": asterisk.FromNagarsevak(
				"single", req.PhoneNumber, req.Language,
				req.Pincode, req.Ward, ns.ID.String(), ns.Name,
			),
		})

	case count >= 2 && count <= 5:
		list := make([]models.NagarsevakItem, 0, count)

		for _, ns := range nagarsevaks {
			list = append(list, models.NagarsevakItem{
				ID:   ns.ID.String(),
				Name: ns.Name,
			})
		}

		c.JSON(http.StatusOK, gin.H{
			"status": "choose",
			"list":   list,
			"_channel_vars": asterisk.FromNagarsevak(
				"choose", req.PhoneNumber, req.Language,
				req.Pincode, req.Ward, "", "",
			),
		})

	default:
		c.JSON(http.StatusOK, gin.H{
			"status": "too_many",
			"_channel_vars": asterisk.FromNagarsevak(
				"too_many", req.PhoneNumber, req.Language,
				req.Pincode, req.Ward, "", "",
			),
		})
	}
}

func (h *NagarsevakHandler) CompleteCitizen(c *gin.Context) {
	var req models.RegisterCitizenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	ns, err := h.politicalRepo.FindNagarsevakByID(
		c.Request.Context(),
		req.NagarsevakID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if ns == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "nagarsevak not found"})
		return
	}

	err = h.citizenRepo.UpsertCitizen(
		c.Request.Context(),
		req.PhoneNumber,
		req.Language,
		req.Pincode,
		req.Ward,
		ns.ID,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success":         true,
		"saved":           true,
		"nagarsevak_name": ns.Name,
		"_channel_vars": asterisk.FromComplete(
			req.PhoneNumber, req.Language, req.Pincode, req.Ward, ns.Name,
		),
	})
}
