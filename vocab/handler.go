package vocab

import (
	"encoding/json"
	"errors"
	"log"
	"time"
	"vocab8/config"
	"vocab8/domain/transport"
	"vocab8/vocab/repository"

	"github.com/gin-gonic/gin"
)

func DrawVocab(c *gin.Context) {
	r := Draw_V2()

	speechURL, err := speechVocab(r.Id)
	if err != nil {
		log.Println(err)
		handleResponse(c, 500, r)
		return
	}
	r.Speech = speechURL

	r.Coverage, r.TodayActivity, err = calculateCoverage()
	if err != nil {
		log.Println(err)
		handleResponse(c, 500, r)
		return
	}

	handleResponse(c, 200, r)
}

func AddVocab(c *gin.Context) {
	var r transport.Word
	err := c.Bind(&r)
	if err != nil {
		log.Println(err)
		handleResponse(c, 400, r)
		return
	}

	if len(r.Title) == 0 || len(r.Description) == 0 {
		err = errors.New("empty body")
		log.Println(err)
		handleResponse(c, 400, r)
		return
	}

	if r.Id > 0 {
		err = repository.UpdateWord(r.Id, r.Title, r.Description)
		if err != nil {
			log.Println(err)
			handleResponse(c, 500, r)
			return
		}
	} else {
		err = repository.InsertWord(r.Title, r.Description)
		if err != nil {
			log.Println(err)
			handleResponse(c, 500, r)
			return
		}
	}

	r.Status = true
	handleResponse(c, 200, r)
}

func SubmitFeedback(c *gin.Context) {
	var r transport.Feedback
	err := c.Bind(&r)
	if err != nil {
		log.Println(err)
		handleResponse(c, 400, r)
		return
	}

	if r.Id <= 0 || (!r.Success && !r.Fail) {
		err = errors.New("empty body")
		log.Println(err)
		handleResponse(c, 400, r)
		return
	}

	data, err := repository.FetchDrawById(r.Id)
	if err != nil {
		log.Println(err)
		handleResponse(c, 500, r)
		return
	}

	switch {
	case r.Success:
		data.DrawSuccess++
	case r.Fail:
		data.DrawFail++
	}
	data.DrawNo++
	data.UpdatedAt = time.Now()

	err = repository.UpdateDrawStats(data.Id, data.DrawFail, data.DrawNo, data.DrawSuccess)
	if err != nil {
		log.Println(err)
		handleResponse(c, 500, r)
		return
	}

	delete(config.WordPool, r.Id)

	r.Status = true
	handleResponse(c, 200, r)
}

func LastFeedbackedVocab(c *gin.Context) {
	var r []transport.Word

	words, err := repository.FetchLastNDraw(100)
	if err != nil {
		log.Println(err)
		handleResponse(c, 500, r)
		return
	}

	for _, word := range words {
		r = append(r, transport.Word{
			Title: word.Word,
		})
	}

	handleResponse(c, 200, r)
}

// todo: consider adding error code on failure
func handleResponse(c *gin.Context, statusCode int, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	c.Data(statusCode, "application/json; charset=utf-8", b)
}
