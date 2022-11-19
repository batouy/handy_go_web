package models

import (
	"database/sql"
	"time"
)

// 实际存储 博客 数据
type Blog struct {
	ID        int
	Title     string
	Content   string
	CreatedAt time.Time
	UpdatedAt time.Time
}

// 封装与数据库交互功能的
type BlogModel struct {
	DB *sql.DB
}

func (m *BlogModel) Insert(title, content string) (int, error) {
	stmt := "insert into blogs (title, content, created_at, updated_at) values(?,?,?,?)"
	nowTime := time.Now()
	result, err := m.DB.Exec(stmt, title, content, nowTime, nowTime)
	if err != nil {
		return 0, nil
	}

	id, err := result.LastInsertId()
	if err != nil {
		return 0, nil
	}
	return int(id), nil
}

func (m *BlogModel) Get(id int) (*Blog, error) {
	stmt := "select id, title, content, created_at, updated_at from blogs where id=?"
	row := m.DB.QueryRow(stmt, id)
	blog := &Blog{}

	err := row.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return blog, nil
}

func (m *BlogModel) Latest() ([]*Blog, error) {
	stmt := "select id, title, content, created_at, updated_at from blogs order by updated_at desc limit 5"
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, err
	}

	blogs := []*Blog{}

	for rows.Next() {
		blog := &Blog{}

		err := rows.Scan(&blog.ID, &blog.Title, &blog.Content, &blog.CreatedAt, &blog.UpdatedAt)
		if err != nil {
			return nil, err
		}

		blogs = append(blogs, blog)
	}

	return blogs, nil
}
