package db

import "time"

type Word struct {
	tableName   struct{}  `pg:"word"`
	Id          int       `pg:"id"`
	Word        string    `pg:"word"`
	Description string    `pg:"description"`
	DrawNo      int       `pg:"draw_no,use_zero"`
	DrawSuccess int       `pg:"draw_success,use_zero"`
	DrawFail    int       `pg:"draw_fail,use_zero"`
	CreatedAt   time.Time `pg:"created_at"`
	UpdatedAt   time.Time `pg:"updated_at"`
}
type WordMigrations struct {
	tableName struct{}  `pg:"word_migrations"`
	Id        int       `pg:"id"`
	No        int       `pg:"no"`
	UpdatedAt time.Time `pg:"updated_at"`
}
