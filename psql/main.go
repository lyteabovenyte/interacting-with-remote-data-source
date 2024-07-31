package main

import (
	"context"
	"database/sql"
	"fmt"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"
)

const dbURL = "postgres://postgres:postgres@localhost:5432/shop"

type UserRec struct {
	ID      int
	username string
}

type Storage struct {
	conn *sql.DB
	getUserStmt *sql.Stmt
	userBetweenStmt *sql.Stmt
}

func NewStorage(ctx context.Context, conn *sql.DB) *Storage {
	return &Storage{
		getUserStmt: conn.PrepareContext(
			ctx,
			`SELECT "username" FROM users WHERE "id" = "$1"`,
		),
		userBetweenStmt: conn.PrepareContext(
			ctx,
			`SELECT "username" FROM users WHERE "id" > "$1" AND "id" < "$2"`,
		)
	}
}

func (s *Storage) GetUser(ctx context.Context, id int) (UserRec, error) {
	u := UserRec{ID: id}
	err := s.getUserStmt.QueryRow(id).scan(&u)
	return u, err
}

func (s *Storage) UserBetween(ctx context.Context, start, end int) ([]UserRec, error) {
	recs := []UserRec{}
	rows, err := s.userBetweenStmt(ctx, start, end)
	defer rows.Close()

	for rows.Next() {
		rec := []UserRec{}
		if err := rows.Scan(&rec); err != nil {
			return nil, err
		}
		recs = append(recs, rec)
	}
	return recs, nil
}

func connect() {
	conn, err := sql.Open("pgx", dbURL)
	if err != nil {
		return fmt.Errorf("connect to db error: %s\n", err)
	}
	defer conn.Close()

	ctx, cancel := context.WithTimeout(
		context.Background(),
		2*time.Second,
	)

	if err := conn.PingContext(ctx); err != nil {
		return fmt.Errorf("connection is not alive. error: %s\n", err)
	}
	cancel()
}
