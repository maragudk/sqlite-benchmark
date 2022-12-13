package sqlite

import (
	"database/sql"
	_ "embed"
	"errors"
	"sync"
)

type DB struct {
	DB       *sql.DB
	useMutex bool
	mutex    sync.RWMutex
}

func NewDB(db *sql.DB, useMutex bool) *DB {
	return &DB{DB: db, useMutex: useMutex}
}

func (d *DB) WritePost(title, content string) error {
	if d.useMutex {
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}

	_, err := d.DB.Exec(`insert into posts (title, content) values ($1, $2)`, title, content)
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
	if d.useMutex {
		d.mutex.RLock()
		defer d.mutex.RUnlock()
	}

	p = &Post{ID: id}

	row := d.DB.QueryRow(`select title, content from posts where id = $1`, id)
	if err = row.Scan(&p.Title, &p.Content); err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			p = nil
		}
		return
	}

	var rows *sql.Rows
	rows, err = d.DB.Query(`select id, name, content from comments where post_id = $1 order by created`, id)
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
		d.mutex.Lock()
		defer d.mutex.Unlock()
	}

	_, err := d.DB.Exec(`insert into comments (post_id, name, content) values ($1, $2, $3)`, postID, name, content)
	return err
}
