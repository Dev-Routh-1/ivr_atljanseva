package handler

import (
	"net/http"
	"log"

	"github.com/gin-gonic/gin"
	"ivr_ataljanseva/db/repository"
	"ivr_ataljanseva/models"

	
)

type CitizenHandler struct {
	repo *repository.CitizenRepository
}

func NewCitizenHandler(
	repo *repository.CitizenRepository,
) *CitizenHandler {
	return &CitizenHandler{
		repo: repo,
	}
}

func (h *CitizenHandler) GetCitizen(c *gin.Context) {
	phone := c.Param("phone")

	citizen, err := h.repo.FindByPhone(
		c.Request.Context(),
		phone,
	)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	if citizen == nil {
		c.JSON(http.StatusOK, models.CitizenLookupResponse{
			Found: false,
		})
		return
	}

	c.JSON(http.StatusOK, models.CitizenLookupResponse{
		Found:      true,
		CitizenIVR: *citizen,
	})
}


func (h *CitizenHandler) RegisterCitizen(c *gin.Context) {

	var req models.RegisterCitizenRequest

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"error": err.Error(),
		})
		return
	}

	err := h.repo.Create(
		c.Request.Context(),
		&req,
	)


	if err != nil {
		log.Printf("register citizen failed: %v", err)

		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"message": "citizen registered successfully",
	})
}