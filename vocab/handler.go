package vocab

import (
	"encoding/json"
	"errors"
	"log"
	"math/rand"
	"time"
	"vocab8/config"
	"vocab8/domain/db"
	"vocab8/domain/transport"

	"github.com/gin-gonic/gin"
	"github.com/go-pg/pg/v10"
)

func RenewThePool() {
	var low_drawed []db.Word
	err := config.DB.Model(&low_drawed).OrderExpr("draw_no ASC").Limit(50).Select()
	if err != nil {
		log.Println(err)
		return
	}

	var high_failed []db.Word
	err = config.DB.Model(&high_failed).
		Where("?<=?", pg.Ident("updated_at"), time.Now().Add(-time.Hour*24)).
		OrderExpr("draw_fail DESC").
		Limit(50).
		Select()
	if err != nil {
		log.Println(err)
		return
	}

	var all_words []db.Word
	all_words = append(all_words, low_drawed...)
	all_words = append(all_words, high_failed...)

	dest := make([]transport.Word, len(all_words))
	perm := rand.Perm(len(all_words))
	for i, v := range perm {
		dest[v] = transport.Word{
			Title:       all_words[i].Word,
			Description: all_words[i].Description,
			Id:          all_words[i].Id,
		}
	}

	config.WordPool = dest
}

func DrawVocab(c *gin.Context) {
	var r transport.Word

	pool_length := len(config.WordPool)
	word_insex := rand.Intn(pool_length - 1)
	r = config.WordPool[word_insex]

	handleResponse(c, r)
}

func AddVocab(c *gin.Context) {
	var r transport.Word
	err := c.Bind(&r)
	if err != nil {
		log.Println(err)
		handleResponse(c, r)
		return
	}

	if len(r.Title) == 0 || len(r.Description) == 0 {
		err = errors.New("empty body")
		log.Println(err)
		handleResponse(c, r)
		return
	}

	data := db.Word{
		Word:        r.Title,
		Description: r.Description,
	}
	_, err = config.DB.Model(&data).Insert()
	if err != nil {
		log.Println(err)
		handleResponse(c, r)
		return
	}

	r.Status = true
	handleResponse(c, r)
}

func SubmitFeedback(c *gin.Context) {
	var r transport.Feedback
	err := c.Bind(&r)
	if err != nil {
		log.Println(err)
		handleResponse(c, r)
		return
	}

	if r.Id <= 0 || (!r.Success && !r.Fail) {
		err = errors.New("empty body")
		log.Println(err)
		handleResponse(c, r)
		return
	}

	var data db.Word
	data.Id = r.Id
	err = config.DB.Model(&data).WherePK().Select()
	if err != nil {
		log.Println(err)
		handleResponse(c, r)
		return
	}

	switch {
	case r.Success:
		data.DrawSuccess++
	case r.Fail:
		data.DrawFail++
	}
	data.DrawNo++

	_, err = config.DB.Model(&data).WherePK().Update()
	if err != nil {
		log.Println(err)
		handleResponse(c, r)
		return
	}

	r.Status = true
	handleResponse(c, r)
}

func handleResponse(c *gin.Context, body interface{}) {
	b, err := json.Marshal(body)
	if err != nil {
		log.Println(err)
	}
	c.Data(200, "application/json; charset=utf-8", b)
}
