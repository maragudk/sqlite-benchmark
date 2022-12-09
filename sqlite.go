package sqlite

import (
	"database/sql"
	"errors"
	"sync"
)

type DB struct {
	DB         *sql.DB
	useMutex   bool
	writeMutex sync.Mutex
}

func NewDB(path string, useMutex bool) (*DB, error) {
	db, err := sql.Open("sqlite3", path+"?_journal=WAL&_timeout=10000&_fk=true")
	if err != nil {
		return nil, err
	}
	newDB := &DB{DB: db, useMutex: useMutex}

	_, err = db.Exec(`
			create table posts (
				id integer primary key,
				title text not null,
				content text not null,
				created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ'))
			)`)
	if err != nil {
		return newDB, err
	}

	_, err = db.Exec(`
			create table comments (
				id integer primary key,
				post_id int not null references posts (id),
				name text not null,
				content text not null,
				created text not null default (strftime('%Y-%m-%dT%H:%M:%fZ'))
			)`)
	if err != nil {
		return newDB, err
	}

	return newDB, nil
}

func (d *DB) WritePost(title, content string) error {
	if d.useMutex {
		d.writeMutex.Lock()
		defer d.writeMutex.Unlock()
	}

	_, err := d.DB.Exec(`insert into posts (title, content) values (?, ?)`, title, content)
	return err
}

type Post struct {
	ID      int
	Title   string
	Content string
}

type Comment struct {
	ID      int
	Name    string
	Content string
}

func (d *DB) ReadPost(id int) (p *Post, cs []*Comment, err error) {
	p = &Post{ID: id}

	row := d.DB.QueryRow(`select title, content from posts where id = ?`, id)
	if err = row.Scan(&p.Title, &p.Content); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return
	}

	var rows *sql.Rows
	rows, err = d.DB.Query(`select id, name, content from comments where post_id = ? order by created`, id)
	if err != nil {
		return
	}

	for rows.Next() {
		c := &Comment{}
		if err = rows.Scan(&c.ID, &c.Name, &c.Content); err != nil {
			return
		}
		cs = append(cs, c)
	}

	err = rows.Err()

	return
}

func (d *DB) WriteComment(postID int, name, content string) error {
	if d.useMutex {
		d.writeMutex.Lock()
		defer d.writeMutex.Unlock()
	}

	_, err := d.DB.Exec(`insert into comments (post_id, name, content) values (?, ?, ?)`, postID, name, content)
	return err
}
