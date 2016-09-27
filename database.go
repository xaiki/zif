package zif

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

type Database struct {
	path string
	conn *sql.DB
}

func NewDatabase(path string) *Database {
	var db Database
	db.path = path

	return &db
}

func (db *Database) Connect() error {
	var err error

	db.conn, err = sql.Open("sqlite3", db.path)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(sql_create_post_table)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(sql_create_fts_post)
	if err != nil {
		return err
	}

	_, err = db.conn.Exec(sql_create_upload_date_index)
	if err != nil {
		return err
	}

	return nil
}

func (db *Database) InsertPost(post Post) error {
	// TODO: Is preparing all statements before hand worth doing for perf?
	stmt, err := db.conn.Prepare(sql_insert_post)
	if err != nil {
		return err
	}

	_, err = stmt.Exec(post.InfoHash, post.Title, post.Size, post.FileCount, post.Seeders,
		post.Leechers, post.UploadDate, post.Source[:], post.Tags)

	if err != nil {
		return err
	}

	return nil
}

func (db *Database) GenerateFts(since uint64) error {
	stmt, err := db.conn.Prepare(sql_generate_fts)

	if err != nil {
		return err
	}

	stmt.Exec(since)

	return nil
}

func (db *Database) QueryRecent(page int) ([]Post, error) {
	page_size := 25
	posts := make([]Post, 0, page_size)

	rows, err := db.conn.Query(sql_query_recent_post, page_size*page,
		page_size)

	if err != nil {
		return nil, err
	}

	for rows.Next() {
		var post Post

		err := rows.Scan(&post.Id, &post.InfoHash, &post.Title, &post.Size,
			&post.FileCount, &post.Seeders, &post.Leechers, &post.UploadDate,
			&post.Source, &post.Tags)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (db *Database) Search(query string, page int) ([]Post, error) {
	page_size := 25 // TODO: Configure this elsewhere

	posts := make([]Post, 0, page_size)
	rows, err := db.conn.Query(sql_search_post, query, page*page_size,
		page_size)

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var result uint

		err = rows.Scan(&result)

		if err != nil {
			return nil, err
		}

		post, err := db.QueryPostId(result)

		if err != nil {
			return nil, err
		}

		posts = append(posts, post)
	}

	return posts, nil
}

func (db *Database) QueryPostId(id uint) (Post, error) {
	var post Post
	rows, err := db.conn.Query(sql_query_post_id, id)

	if err != nil {
		return post, err
	}

	for rows.Next() {

		err := rows.Scan(&post.Id, &post.InfoHash, &post.Title, &post.Size,
			&post.FileCount, &post.Seeders, &post.Leechers, &post.UploadDate,
			&post.Source, &post.Tags)

		if err != nil {
			return post, err
		}
	}

	return post, nil
}

func (db *Database) QueryPiece(id int) (*Piece, error) {
	page_size := PieceSize // TODO: Configure this elsewhere
	var piece Piece

	rows, err := db.conn.Query(sql_query_paged_post, id*page_size,
		page_size)

	if err != nil {
		return nil, err
	}

	for rows.Next() {

		var post Post

		err := rows.Scan(&post.Id, &post.InfoHash, &post.Title, &post.Size,
			&post.FileCount, &post.Seeders, &post.Leechers, &post.UploadDate,
			&post.Source, &post.Tags)

		if err != nil {
			return nil, err
		}

		piece.Add(post)
	}

	return &piece, nil
}

func (db *Database) PostCount() uint {
	var res uint

	db.conn.QueryRow(sql_count_post).Scan(&res)

	return res
}

func (db *Database) Close() {
	db.conn.Close()
}
