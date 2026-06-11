package datastore

import (
	"database/sql"
	"fmt"

	_ "modernc.org/sqlite"
)

type DB struct {
	writer *sql.DB
	reader *sql.DB
}

func Open(location string) (*DB, error) {
	writerDSN := fmt.Sprintf(
		"file:%s?_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)&_pragma=foreign_keys(ON)",
		location,
	)
	readerDSN := fmt.Sprintf(
		"file:%s?mode=ro&_pragma=journal_mode(WAL)&_pragma=busy_timeout(5000)",
		location,
	)

	writer, err := sql.Open("sqlite", writerDSN)
	if err != nil {
		return nil, err
	}
	writer.SetMaxOpenConns(1)

	reader, err := sql.Open("sqlite", readerDSN)
	if err != nil {
		writer.Close()
		return nil, err
	}
	reader.SetMaxOpenConns(8)

	return &DB{writer: writer, reader: reader}, nil
}

func (d *DB) Writer() *sql.DB { return d.writer }
func (d *DB) Reader() *sql.DB { return d.reader }

func (d *DB) Close() error {
	werr := d.writer.Close()
	rerr := d.reader.Close()
	if werr != nil {
		return werr
	}
	return rerr
}
