package handler

import (
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"

	"ivr_ataljanseva/db/repository"
	"ivr_ataljanseva/models"
	"ivr_ataljanseva/plivo"
)

type PlivoHandler struct {
	citizenRepo   *repository.CitizenRepository
	politicalRepo *repository.PoliticalUserRepository
	baseURL       string
	audioBaseURL  string
	maxRetries    int
}

func NewPlivoHandler(
	citizenRepo *repository.CitizenRepository,
	politicalRepo *repository.PoliticalUserRepository,
	baseURL, audioBaseURL string,
) *PlivoHandler {
	return &PlivoHandler{
		citizenRepo:   citizenRepo,
		politicalRepo: politicalRepo,
		baseURL:       baseURL,
		audioBaseURL:  audioBaseURL,
		maxRetries:    3,
	}
}

// POST /ivr/plivo/incoming
func (h *PlivoHandler) Incoming(c *gin.Context) {
	phone := c.PostForm("From")
	if phone == "" {
		c.String(http.StatusBadRequest, plivo.Response(
			plivo.Speak("Invalid request. No caller ID.", "english"),
			plivo.Hangup(),
		))
		return
	}

	citizen, err := h.citizenRepo.FindByPhone(c.Request.Context(), phone)
	if err != nil {
		log.Printf("citizen lookup error: %v", err)
		c.String(http.StatusOK, plivo.Response(
			plivo.Speak("System error. Please try again later.", "english"),
			plivo.Hangup(),
		))
		return
	}

	if citizen != nil {
		lang := citizen.Language
		nsName := ""
		if citizen.NagarsevakID != uuid.Nil {
			ns, err := h.politicalRepo.FindNagarsevakByID(c.Request.Context(), citizen.NagarsevakID)
			if err == nil && ns != nil {
				nsName = ns.Name
			}
		}
		h.returnMainMenu(c, phone, lang, citizen.Pincode, citizen.Ward, citizen.NagarsevakID.String(), nsName)
		return
	}

	audio := plivo.AudioURL(h.audioBaseURL, "english", "welcome.mp3")
	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.GetDigits(h.baseURL+"/ivr/plivo/language?phone="+phone, 1, 10),
	))
}

// POST /ivr/plivo/language
func (h *PlivoHandler) LanguageSelect(c *gin.Context) {
	phone := c.Query("phone")
	digits := c.PostForm("Digits")

	lang := resolveLanguage(digits)
	audio := plivo.AudioURL(h.audioBaseURL, lang, "ward-input.mp3")
	action := h.baseURL + "/ivr/plivo/ward-input?phone=" + phone + "&language=" + lang

	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.GetDigits(action, 20, 15),
	))
}

// POST /ivr/plivo/ward-input
func (h *PlivoHandler) WardInput(c *gin.Context) {
	phone := c.Query("phone")
	lang := c.Query("language")
	digits := c.PostForm("Digits")
	retryStr := c.Query("retry")

	if digits == "" {
		h.wardInputRetry(c, phone, lang, retryStr)
		return
	}

	retry, _ := strconv.Atoi(retryStr)

	pincode, wardInput := splitPincodeWard(digits)
	if pincode == "" || wardInput == "" {
		h.wardInputRetry(c, phone, lang, strconv.Itoa(retry+1))
		return
	}

	matches, err := h.politicalRepo.FindMatchingWards(c.Request.Context(), pincode, wardInput)
	if err != nil {
		log.Printf("ward resolve error: %v", err)
		c.String(http.StatusOK, plivo.Response(
			plivo.Speak("System error. Please try again later.", lang),
			plivo.Hangup(),
		))
		return
	}

	switch {
	case len(matches) == 0:
		h.wardInputRetry(c, phone, lang, strconv.Itoa(retry+1))

	case len(matches) == 1:
		m := matches[0]
		err := h.citizenRepo.UpsertCitizen(c.Request.Context(), phone, lang, pincode, m.Ward, m.NagarsevakID)
		if err != nil {
			log.Printf("auto-save error: %v", err)
		}
		h.returnNagarsevakSingle(c, phone, lang, pincode, m.Ward, m.NagarsevakID.String(), m.NagarsevakName)

	case len(matches) <= 4:
		h.returnWardMenu(c, phone, lang, pincode, matches)

	default:
		h.returnWhatsAppPrompt(c, lang)
	}
}

// POST /ivr/plivo/ward-select
func (h *PlivoHandler) WardSelect(c *gin.Context) {
	phone := c.Query("phone")
	lang := c.Query("language")
	pincode := c.Query("pincode")
	wardsRaw := c.Query("wards")
	digits := c.PostForm("Digits")

	idx, _ := strconv.Atoi(digits)
	idx-- // 1-indexed to 0-indexed

	wards := strings.Split(wardsRaw, ",")
	if idx < 0 || idx >= len(wards) {
		action := h.baseURL + "/ivr/plivo/ward-input?phone=" + phone + "&language=" + lang + "&retry="
		c.String(http.StatusOK, plivo.Response(
			plivo.Speak("Invalid selection. Please try again.", lang),
			plivo.GetDigits(action, 20, 15),
		))
		return
	}

	selectedWard := strings.TrimSpace(wards[idx])

	nagarsevaks, err := h.politicalRepo.FindNagarsevaks(c.Request.Context(), pincode, selectedWard)
	if err != nil {
		log.Printf("nagarsevak lookup error: %v", err)
		c.String(http.StatusOK, plivo.Response(
			plivo.Speak("System error. Please try again later.", lang),
			plivo.Hangup(),
		))
		return
	}

	switch {
	case len(nagarsevaks) == 0:
		h.returnWhatsAppPrompt(c, lang)

	case len(nagarsevaks) == 1:
		ns := nagarsevaks[0]
		err := h.citizenRepo.UpsertCitizen(c.Request.Context(), phone, lang, pincode, selectedWard, ns.ID)
		if err != nil {
			log.Printf("auto-save error: %v", err)
		}
		h.returnNagarsevakSingle(c, phone, lang, pincode, selectedWard, ns.ID.String(), ns.Name)

	case len(nagarsevaks) <= 5:
		h.returnNagarsevakMenu(c, phone, lang, pincode, selectedWard, nagarsevaks)

	default:
		h.returnWhatsAppPrompt(c, lang)
	}
}

// POST /ivr/plivo/nagarsevak-select
func (h *PlivoHandler) NagarsevakSelect(c *gin.Context) {
	phone := c.Query("phone")
	lang := c.Query("language")
	pincode := c.Query("pincode")
	ward := c.Query("ward")
	idsRaw := c.Query("ids")
	digits := c.PostForm("Digits")

	idx, _ := strconv.Atoi(digits)
	idx--

	ids := strings.Split(idsRaw, ",")
	if idx < 0 || idx >= len(ids) {
		action := h.baseURL + "/ivr/plivo/ward-select?phone=" + phone + "&language=" + lang + "&pincode=" + pincode + "&wards=" + ward
		c.String(http.StatusOK, plivo.Response(
			plivo.Speak("Invalid selection. Please try again.", lang),
			plivo.GetDigits(action, 1, 10),
		))
		return
	}

	nsID := strings.TrimSpace(ids[idx])
	parsedUUID, err := uuid.Parse(nsID)
	if err != nil {
		c.String(http.StatusOK, plivo.Response(
			plivo.Speak("System error. Please try again later.", lang),
			plivo.Hangup(),
		))
		return
	}

	ns, err := h.politicalRepo.FindNagarsevakByID(c.Request.Context(), parsedUUID)
	if err != nil || ns == nil {
		c.String(http.StatusOK, plivo.Response(
			plivo.Speak("System error. Please try again later.", lang),
			plivo.Hangup(),
		))
		return
	}

	err = h.citizenRepo.UpsertCitizen(c.Request.Context(), phone, lang, pincode, ward, ns.ID)
	if err != nil {
		log.Printf("save error: %v", err)
	}

	h.returnRegisteredConfirmation(c, lang, ns.Name)
}

// POST /ivr/plivo/main-menu
func (h *PlivoHandler) MainMenu(c *gin.Context) {
	phone := c.Query("phone")
	lang := c.Query("language")
	pincode := c.Query("pincode")
	ward := c.Query("ward")
	nsID := c.Query("nagarsevak_id")
	nsName := c.Query("nagarsevak_name")
	digits := c.PostForm("Digits")

	switch digits {
	case "1":
		h.returnSOS(c, lang)
	case "2":
		h.returnComplaint(c, lang)
	case "3":
		h.corporatorConnect(c, phone, lang, pincode, ward, nsID, nsName)
	default:
		h.returnMainMenu(c, phone, lang, pincode, ward, nsID, nsName)
	}
}

// -------------------- internal response builders --------------------

func (h *PlivoHandler) returnMainMenu(c *gin.Context, phone, lang, pincode, ward, nsID, nsName string) {
	if nsName == "" {
		nsName = "your nagarsevak"
	}

	action := h.baseURL + "/ivr/plivo/main-menu?phone=" + phone + "&language=" + lang +
		"&pincode=" + pincode + "&ward=" + ward +
		"&nagarsevak_id=" + nsID + "&nagarsevak_name=" + nsName

	audio := plivo.AudioURL(h.audioBaseURL, lang, "main-menu.mp3")

	c.String(http.StatusOK, plivo.Response(
		plivo.Speak("Welcome back! Your nagarsevak "+nsName+" is connected.", lang),
		plivo.Wait(1),
		plivo.Play(audio),
		plivo.GetDigits(action, 1, 10),
	))
}

func (h *PlivoHandler) returnSOS(c *gin.Context, lang string) {
	audio := plivo.AudioURL(h.audioBaseURL, lang, "sos.mp3")
	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.Hangup(),
	))
}

func (h *PlivoHandler) returnComplaint(c *gin.Context, lang string) {
	audio := plivo.AudioURL(h.audioBaseURL, lang, "complaint.mp3")
	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.Hangup(),
	))
}

func (h *PlivoHandler) corporatorConnect(c *gin.Context, phone, lang, pincode, ward, nsID, nsName string) {
	if nsName == "" {
		nsName = "your nagarsevak"
	}

	audio := plivo.AudioURL(h.audioBaseURL, lang, "corporator-connect.mp3")

	c.String(http.StatusOK, plivo.Response(
		plivo.Speak("Connecting you to "+nsName+".", lang),
		plivo.Wait(1),
		plivo.Play(audio),
		plivo.Hangup(),
	))
}

func (h *PlivoHandler) wardInputRetry(c *gin.Context, phone, lang, retryStr string) {
	retry, _ := strconv.Atoi(retryStr)
	if retry >= h.maxRetries {
		h.returnWhatsAppPrompt(c, lang)
		return
	}

	audio := plivo.AudioURL(h.audioBaseURL, lang, "no-match.mp3")
	action := h.baseURL + "/ivr/plivo/ward-input?phone=" + phone + "&language=" + lang + "&retry=" + strconv.Itoa(retry)

	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.GetDigits(action, 20, 15),
	))
}

func (h *PlivoHandler) returnWardMenu(c *gin.Context, phone, lang, pincode string, matches []models.WardMatch) {
	var wards []string
	var ttsParts []string
	for i, m := range matches {
		wards = append(wards, m.Ward)
		ttsParts = append(ttsParts, "Press "+strconv.Itoa(i+1)+" for "+m.Ward)
	}

	audio := plivo.AudioURL(h.audioBaseURL, lang, "ward-menu.mp3")
	action := h.baseURL + "/ivr/plivo/ward-select?phone=" + phone + "&language=" + lang +
		"&pincode=" + pincode + "&wards=" + strings.Join(wards, ",")

	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.Speak(strings.Join(ttsParts, ". ")+".", lang),
		plivo.GetDigits(action, 1, 10),
	))
}

func (h *PlivoHandler) returnNagarsevakSingle(c *gin.Context, phone, lang, pincode, ward, nsID, nsName string) {
	if nsName == "" {
		nsName = "your nagarsevak"
	}

	audio := plivo.AudioURL(h.audioBaseURL, lang, "goodbye.mp3")

	c.String(http.StatusOK, plivo.Response(
		plivo.Speak("Your ward is "+ward+". Your nagarsevak is "+nsName+".", lang),
		plivo.Wait(1),
		plivo.Play(audio),
		plivo.Hangup(),
	))
}

func (h *PlivoHandler) returnNagarsevakMenu(c *gin.Context, phone, lang, pincode, ward string, nagarsevaks []models.NagarsevakRecord) {
	var ids []string
	var ttsParts []string
	for i, ns := range nagarsevaks {
		ids = append(ids, ns.ID.String())
		ttsParts = append(ttsParts, "Press "+strconv.Itoa(i+1)+" for "+ns.Name)
	}

	audio := plivo.AudioURL(h.audioBaseURL, lang, "nagarsevak-menu.mp3")
	action := h.baseURL + "/ivr/plivo/nagarsevak-select?phone=" + phone + "&language=" + lang +
		"&pincode=" + pincode + "&ward=" + ward + "&ids=" + strings.Join(ids, ",")

	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.Speak(strings.Join(ttsParts, ". ")+".", lang),
		plivo.GetDigits(action, 1, 10),
	))
}

func (h *PlivoHandler) returnRegisteredConfirmation(c *gin.Context, lang, nsName string) {
	if nsName == "" {
		nsName = "your nagarsevak"
	}
	audio := plivo.AudioURL(h.audioBaseURL, lang, "goodbye.mp3")
	c.String(http.StatusOK, plivo.Response(
		plivo.Speak("You are now registered. Your nagarsevak is "+nsName+".", lang),
		plivo.Wait(1),
		plivo.Play(audio),
		plivo.Hangup(),
	))
}

func (h *PlivoHandler) returnWhatsAppPrompt(c *gin.Context, lang string) {
	audio := plivo.AudioURL(h.audioBaseURL, lang, "whatsapp.mp3")
	c.String(http.StatusOK, plivo.Response(
		plivo.Play(audio),
		plivo.Hangup(),
	))
}

// -------------------- helpers --------------------

func resolveLanguage(digits string) string {
	switch digits {
	case "1":
		return "english"
	case "2":
		return "hindi"
	case "3":
		return "marathi"
	default:
		return "english"
	}
}

func splitPincodeWard(s string) (string, string) {
	s = strings.TrimSpace(s)
	if idx := strings.Index(s, "#"); idx > 0 {
		return strings.TrimSpace(s[:idx]), strings.TrimSpace(s[idx+1:])
	}
	// fallback: first 6 digits as pincode, rest as ward
	if len(s) >= 6 {
		return s[:6], strings.TrimSpace(s[6:])
	}
	return s, ""
}
