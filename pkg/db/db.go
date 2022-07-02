package db

import (
	"database/sql"
	"fmt"

	_ "github.com/mattn/go-sqlite3"
	"github.com/pkg/errors"
	log "github.com/sirupsen/logrus"
)

type DB struct {
	path    string
	db      *sql.DB
	queue   chan Annotation
	tblName string
}

func NewDb(dbPath string) (*DB, error) {
	db, err := sql.Open("sqlite3", dbPath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to open database")
	}
	return &DB{path: dbPath, db: db, queue: make(chan Annotation)}, nil
}

func (db *DB) Close() error {
	return db.db.Close()
}

func (db *DB) Commit(annotation Annotation) {
	db.queue <- annotation
}

func (db *DB) Run(title string) error {
	db.tblName = title
	if err := db.createTables(title); err != nil {
		return errors.Wrap(err, "failed to create tables")
	}
	go db.processQueue()
	return nil
}

func (db *DB) createTables(tblName string) error {
	_, err := db.db.Exec(fmt.Sprintf(`
		CREATE TABLE IF NOT EXISTS %s (
			id INTEGER PRIMARY KEY,
			type TEXT,
			text TEXT,
			startx Integer,
			starty Integer,
			endx Integer,
			endy Integer,
			color TEXT,
			author TEXT,
			date TEXT,
			comment TEXT
		);
	`, tblName))
	if err != nil {
		return errors.Wrap(err, "failed to create table")
	}
	return nil
}

func (db *DB) processQueue() {
	for anno := range db.queue {
		if err := db.insertAnnotation(anno); err != nil {
			log.Error(err)
		}
	}
}

func (db *DB) insertAnnotation(anno Annotation) error {
	_, err := db.db.Exec(fmt.Sprintf(`
		INSERT INTO %s (type, text, startx, starty, endx, endy, color, author, date, comment)
		VALUES (?, ?, ?, ?, ?, ?, ?, ?, ?, ?);
	`, db.tblName), anno.Type, anno.Text, anno.StartX, anno.StartY, anno.EndX, anno.EndY,
		anno.Color, anno.Author, anno.Date, anno.Comment)
	if err != nil {
		return errors.Wrap(err, "failed to insert annotation")
	}
	return nil
}
