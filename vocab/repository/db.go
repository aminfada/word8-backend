package repository

import (
	"log"
	"time"
	"vocab8/config"
	"vocab8/domain/db"

	"github.com/go-pg/pg/v10"
)

func UpdateWord(id int, word, description string) (err error) {
	data := db.Word{
		Id:          id,
		Word:        word,
		Description: description,
	}
	_, err = config.DB.Model(&data).
		WherePK().
		Column("description").
		Column("word").
		Update()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func UpdateDrawStats(id int, drawFail, drawNo, drawSuccess int) (err error) {
	data := db.Word{
		Id:          id,
		DrawFail:    drawFail,
		DrawNo:      drawNo,
		DrawSuccess: drawSuccess,
	}
	_, err = config.DB.Model(&data).
		WherePK().
		Column("draw_fail").
		Column("draw_success").
		Column("draw_no").
		Update()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func InsertWord(word, description string) (err error) {
	data := db.Word{
		Word:        word,
		Description: description,
	}
	_, err = config.DB.Model(&data).Insert()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func FetchDrawById(id int) (draw db.Word, err error) {
	var data db.Word
	data.Id = id

	err = config.DB.Model(&data).WherePK().Select()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func FetchLastNDraw(n int) (lastDraws []db.Word, err error) {
	err = config.DB.Model(&lastDraws).
		Where("?>?", pg.Ident("draw_no"), 0).
		OrderExpr("updated_at DESC").
		Limit(n).
		Select()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func FetchLowDraw() (lowDraw []db.Word, err error) {
	err = config.DB.Model(&lowDraw).
		Where("?<?", pg.Ident("draw_success"), 7).
		OrderExpr("draw_no ASC").
		Limit(20).
		Select()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func FetchPassedDraw() (passedDraw []db.Word, err error) {
	err = config.DB.Model(&passedDraw).
		Where("?>=?", pg.Ident("draw_success"), 7).
		OrderExpr("draw_no ASC").
		Limit(5).
		Select()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func FetchHighFailedDraw() (highFailedDraw []db.Word, err error) {
	err = config.DB.Model(&highFailedDraw).
		Where("?<=?", pg.Ident("updated_at"), time.Now().Add(-time.Hour*24)).
		Where("?>?", pg.Ident("draw_fail"), 0).
		Where("?<?", pg.Ident("draw_success"), 7).
		Select()
	if err != nil {
		log.Println(err)
		return
	}

	return
}

func CountTodayActivity() (todayActivity int, err error) {
	todayActivity, err = config.DB.Model(&db.Word{}).
		Where("?>?", pg.Ident("draw_no"), 0).
		Where("date(updated_at)=?", time.Now().Format("2006-01-02")).
		Count()
	if err != nil {
		return
	}

	return
}

func CountVisitedDraw() (visitedDraw int, err error) {
	visitedDraw, err = config.DB.Model(&db.Word{}).
		Where("?>?", pg.Ident("draw_no"), 0).
		Count()
	if err != nil {
		return
	}

	return
}

func CountTotalDraw() (totalDraw int, err error) {
	totalDraw, err = config.DB.Model(&db.Word{}).
		Count()
	if err != nil {
		return
	}

	return
}
