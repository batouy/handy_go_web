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
	return 0, nil
}

func (m *BlogModel) Get() (*Blog, error) {
	return nil, nil
}

func (m *BlogModel) Latest() ([]*Blog, error) {
	stmt := "select id, title, content, created_at, updated_at from blogs order by updated_at desc limit 5"
	rows, err := m.DB.Query(stmt)

	if err != nil {
		return nil, nil
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
