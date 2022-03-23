package mysql

import (
	"database/sql"
	"errors"

	"github.com/daffaz/snippetbox/pkg/models"
)

type SnippetModel struct {
	DB *sql.DB
}

func (m *SnippetModel) Insert(title, content, expires string) (int, error) {
	query := `INSERT INTO snippets (title, content, created, expires)
	VALUES (?, ?, UTC_TIMESTAMP(), DATE_ADD(CURRRENT_TIMESTAMP(), INTERVAL ? DAY))`

	res, err := m.DB.Exec(query, title, content, expires)
	if err != nil {
		return 0, err
	}

	id, err := res.LastInsertId()
	if err != nil {
		return 0, err
	}

	return int(id), nil
}

func (m *SnippetModel) Get(id int) (*models.Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP() AND id = ?`
	s := &models.Snippet{}
	row := m.DB.QueryRow(query, id)
	err := row.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires)
	if err != nil {
		if errors.Is(err, sql.ErrNoRows) {
			return nil, models.ErrNoRecord
		} else {
			return nil, err
		}
	}

	return s, nil
}

func (m *SnippetModel) Latest() ([]*models.Snippet, error) {
	query := `SELECT id, title, content, created, expires FROM snippets
	WHERE expires > CURRENT_TIMESTAMP() ORDER BY created LIMIT 10`

	rows, err := m.DB.Query(query)
	if err != nil {
		return nil, err
	}

	defer rows.Close()
	snippets := []*models.Snippet{}

	for rows.Next() {
		s := &models.Snippet{}

		if err = rows.Scan(&s.ID, &s.Title, &s.Content, &s.Created, &s.Expires); err != nil {
			return nil, err
		}

		snippets = append(snippets, s)
	}

	if err = rows.Err(); err != nil {
		return nil, err
	}

	return snippets, nil
}
