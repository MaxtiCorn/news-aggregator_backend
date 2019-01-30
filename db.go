package main

import (
	"database/sql"

	_ "github.com/mattn/go-sqlite3"
)

const (
	createTableNewsQuery = `
		create table if not exists "news" (
		"id" integer primary key autoincrement,
    	"title" text,
    	"link" varchar(64) unique,
    	"description" text
	);`

	createIndexQuery = `
		create index if not exists "news_title_index"
		on news(title);
	`

	insertNewsQuery = `
		insert or ignore into news(title, link, description) 
		values(?,?,?);
	`

	selectAllNewsQuery = `
		select id, title, link, description from news
		order by id desc
	`

	selectNewsByTitleQuery = `
		select id, title, link, description from news
		where title LIKE '%' || ? || '%'
		order by id desc
	`
)

func newDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	createQueries := []string{createTableNewsQuery, createIndexQuery}

	for _, createQuery := range createQueries {
		statement, err := db.Prepare(createQuery)
		if err != nil {
			return nil, err
		}

		_, err = statement.Exec()
		if err != nil {
			return nil, err
		}

		statement.Close()
	}

	return db, nil
}

func execQueryInTransaction(transaction *sql.Tx, query string, values []interface{}) error {
	statement, err := transaction.Prepare(query)
	if err != nil {
		transactionRollbackErr := transaction.Rollback()
		if transactionRollbackErr != nil {
			panic(transactionRollbackErr)
		}

		return err
	}

	_, err = statement.Exec(values...)
	if err != nil {
		transactionRollbackErr := transaction.Rollback()
		if transactionRollbackErr != nil {
			panic(transactionRollbackErr)
		}

		return err
	}

	return nil
}

func insertNews(transaction *sql.Tx, news *News) (err error) {
	err = execQueryInTransaction(
		transaction,
		insertNewsQuery,
		[]interface{}{news.Title, news.Link, news.Description})

	return
}

func (agr Aggregator) saveNews(news *News) error {
	transaction, err := agr.db.Begin()
	if err != nil {
		return err
	}

	err = insertNews(transaction, news)
	if err != nil {
		return err
	}

	err = transaction.Commit()
	if err != nil {
		panic(err)
	}

	return nil
}

func (agr Aggregator) getAllNews() ([]News, error) {
	statement, err := agr.db.Prepare(selectAllNewsQuery)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query()
	if err != nil {
		return nil, err
	}

	var foundNews []News

	for rows.Next() {
		var id int
		var news News

		if err := rows.Scan(&id, &news.Title, &news.Link, &news.Description); err != nil {
			return nil, err
		}

		foundNews = append(foundNews, news)
	}

	return foundNews, nil
}

func (agr Aggregator) searchNews(query string) ([]News, error) {
	statement, err := agr.db.Prepare(selectNewsByTitleQuery)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(query)
	if err != nil {
		return nil, err
	}

	var foundNews []News

	for rows.Next() {
		var id int
		var news News

		if err := rows.Scan(&id, &news.Title, &news.Link, &news.Description); err != nil {
			return nil, err
		}

		foundNews = append(foundNews, news)
	}

	return foundNews, nil
}
