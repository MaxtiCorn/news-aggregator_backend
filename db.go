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
		"link" varchar(64),
		"source" varchar (64),
    	"description" text
	);`

	createTitleIndexQuery = `
		create index if not exists "news_title_index"
		on news(title);
	`

	createSourceIndexQuery = `
		create index if not exists "news_source_index"
		on news(source);
	`

	insertNewsQuery = `
		insert or ignore into news(title, link, source, description) 
		values(?,?,?,?);
	`

	selectAllNewsQuery = `
		select id, title, link, source, description from news
		order by id desc
	`

	selectNewsQuery = `
		select id, title, link, source, description from news
		order by id desc
		limit ? offset ?
	`

	selectNewsByTitleQuery = `
		select id, title, link, source, description from news
		where title LIKE '%' || ? || '%'
		order by id desc
	`

	selectNewsBySourceQuery = `
		select id, title, link, source, description from news
		where source = ?
		order by id desc
	`
)

func newDatabase(path string) (*sql.DB, error) {
	db, err := sql.Open("sqlite3", path)
	if err != nil {
		return nil, err
	}

	createQueries := []string{createTableNewsQuery, createTitleIndexQuery, createSourceIndexQuery}

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
		[]interface{}{news.Title, news.Link, news.Source, news.Description})
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

func (agr Aggregator) getNews(count, offset string) ([]News, error) {
	statement, err := agr.db.Prepare(selectNewsQuery)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(count, offset)
	if err != nil {
		return nil, err
	}

	var foundNews []News = []News{}

	for rows.Next() {
		var id int
		var news News

		if err := rows.Scan(&id, &news.Title, &news.Link, &news.Source, &news.Description); err != nil {
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

	var foundNews []News = []News{}

	for rows.Next() {
		var id int
		var news News

		if err := rows.Scan(&id, &news.Title, &news.Link, &news.Source, &news.Description); err != nil {
			return nil, err
		}

		foundNews = append(foundNews, news)
	}

	return foundNews, nil
}

func (agr Aggregator) getNewsBySource(source string) ([]News, error) {
	statement, err := agr.db.Prepare(selectNewsBySourceQuery)
	if err != nil {
		return nil, err
	}

	rows, err := statement.Query(source)
	if err != nil {
		return nil, err
	}

	var foundNews []News = []News{}

	for rows.Next() {
		var id int
		var news News

		if err := rows.Scan(&id, &news.Title, &news.Link, &news.Source, &news.Description); err != nil {
			return nil, err
		}

		foundNews = append(foundNews, news)
	}

	return foundNews, nil
}
