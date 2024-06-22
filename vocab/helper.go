package vocab

import (
	"fmt"
	"log"
	"math"
	"math/rand"
	"os"
	"vocab8/config"
	"vocab8/domain/db"
	"vocab8/domain/transport"
	"vocab8/vocab/repository"

	htgotts "github.com/hegedustibor/htgo-tts"
	"github.com/hegedustibor/htgo-tts/handlers"
	"github.com/hegedustibor/htgo-tts/voices"
)

func RenewThePool() (err error) {
	lowDraw, err := repository.FetchLowDraw()
	if err != nil {
		return
	}

	passedDraw, err := repository.FetchPassedDraw()
	if err != nil {
		return
	}

	highFailedDraw, err := repository.FetchHighFailedDraw()
	if err != nil {
		return
	}

	shuffled_high_failed := make([]db.Word, len(highFailedDraw))
	perm := rand.Perm(len(highFailedDraw))
	for i, v := range perm {
		shuffled_high_failed[v] = db.Word{
			Word:        highFailedDraw[i].Word,
			Description: highFailedDraw[i].Description,
			Id:          highFailedDraw[i].Id,
		}
	}

	high_failed_index := math.Max(0, 10)

	var all_words []db.Word
	all_words = append(all_words, lowDraw...)
	all_words = append(all_words, passedDraw...)
	all_words = append(all_words, shuffled_high_failed[:int(high_failed_index)]...)

	dest := make([]transport.Word, len(all_words))
	perm = rand.Perm(len(all_words))
	for i, v := range perm {
		dest[v] = transport.Word{
			Title:       all_words[i].Word,
			Description: all_words[i].Description,
			Id:          all_words[i].Id,
		}
	}

	final_word_pool := make(map[int]transport.Word)
	for _, el := range dest {
		final_word_pool[el.Id] = el
	}

	config.WordPool = final_word_pool

	return
}

func speechVocab(id int) (speechUrl string, err error) {
	var data db.Word
	data.Id = id
	err = config.DB.Model(&data).WherePK().Select()
	if err != nil {
		return
	}

	err = os.Remove(config.Cfg.SpeechPath + "speech.mp3")
	if err != nil {
		return
	}

	speech := htgotts.Speech{Folder: config.Cfg.SpeechPath, Language: voices.English, Handler: &handlers.Native{}}
	filePath, err := speech.CreateSpeechFile(data.Word, "speech")
	if err != nil {
		return
	}

	speechUrl = config.Cfg.SpeechPath + filePath

	return
}

func calculateCoverage() (coverage string, dailyActivity string, err error) {
	todayActivity, err := repository.CountTodayActivity()
	if err != nil {
		return
	}
	log.Println("today activity:", todayActivity)

	visitedDraw, err := repository.CountVisitedDraw()
	if err != nil {
		return
	}
	log.Println("visited draw:", visitedDraw)

	totalDraw, err := repository.CountTotalDraw()
	if err != nil {
		return
	}
	log.Println("total draw:", totalDraw)

	coverageF := (float64(visitedDraw) / float64(totalDraw)) * 100.0
	dailyActivityF := (float64(todayActivity) / float64(100)) * 100.0

	coverage = fmt.Sprintf("%.3f%%", coverageF)
	dailyActivity = fmt.Sprintf("%.3f%%", dailyActivityF)

	log.Println("coverage", coverage)
	log.Println("dailyActivity", dailyActivity)

	return
}

// deprecated
func Draw_V1() (r transport.Word) {
	pool_length := len(config.WordPool)
	word_index := rand.Intn(pool_length - 1)
	r = config.WordPool[word_index]

	return
}

// using map to inherently return shuffled draw
func Draw_V2() (r transport.Word) {
	r = func() transport.Word {
		for _, el := range config.WordPool {
			return el
		}
		return transport.Word{}
	}()

	return
}
