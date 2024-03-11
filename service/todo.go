package service

import (
	"context"
	"database/sql"

	"github.com/TechBowl-japan/go-stations/model"
)

// A TODOService implements CRUD of TODO entities.
type TODOService struct {
	db *sql.DB
}

// NewTODOService returns new TODOService.
func NewTODOService(db *sql.DB) *TODOService {
	return &TODOService{
		db: db,
	}
}

// CreateTODO creates a TODO on DB.
func (s *TODOService) CreateTODO(ctx context.Context, subject, description string) (*model.TODO, error) {
	const (
		insert  = `INSERT INTO todos(subject, description) VALUES(?, ?)`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)

	stmt, err := s.db.PrepareContext(ctx, insert) // SQLステートメントを準備
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, subject, description) // SQLステートメントを実行
	if err != nil {
		return nil, err
	}

	id, err := result.LastInsertId() // resultからIDを取得
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{}
	todo.ID = id                                  // 取得したIDを新しいTODOに格納(sta09でこの部分が引っかかった)
	row := s.db.QueryRowContext(ctx, confirm, id) // 取得したIDを元にTODOを取得
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err // Goの慣習
	}

	return todo, nil // TODOを返す
}

// ReadTODO reads TODOs on DB.
func (s *TODOService) ReadTODO(ctx context.Context, prevID, size int64) ([]*model.TODO, error) {
	const (
		read       = `SELECT id, subject, description, created_at, updated_at FROM todos ORDER BY id DESC LIMIT ?`
		readWithID = `SELECT id, subject, description, created_at, updated_at FROM todos WHERE id < ? ORDER BY id DESC LIMIT ?`
	)
	// PrevID が指定されているかどうかで、定義されている Query を使い分けましょう。
	// 場合分けする必要があるかも。初期値が0だから、指定されていれば0以外の値になっているはず。
	// PrevIDが0かそうでないかで場合分け
	todos := []*model.TODO{} // 空のスライス

	if prevID == 0 {
		rows, err := s.db.QueryContext(ctx, read, size)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() { // 次の行があればtrue, なければfalse
			var todo model.TODO
			if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}
	} else {
		rows, err := s.db.QueryContext(ctx, readWithID, prevID, size)
		if err != nil {
			return nil, err
		}
		defer rows.Close()

		for rows.Next() {
			var todo model.TODO
			if err := rows.Scan(&todo.ID, &todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt); err != nil {
				return nil, err
			}
			todos = append(todos, &todo)
		}
	}

	return todos, nil
}

// UpdateTODO updates the TODO on DB.
func (s *TODOService) UpdateTODO(ctx context.Context, id int64, subject, description string) (*model.TODO, error) {
	const (
		update  = `UPDATE todos SET subject = ?, description = ? WHERE id = ?`
		confirm = `SELECT subject, description, created_at, updated_at FROM todos WHERE id = ?`
	)
	stmt, err := s.db.PrepareContext(ctx, update)
	if err != nil {
		return nil, err
	}
	defer stmt.Close()

	result, err := stmt.ExecContext(ctx, subject, description, id) // updateの順に
	if err != nil {
		return nil, err
	}

	num, err := result.RowsAffected() // RowsAffectedで変更した行数(row)を取得
	if num == 0 {
		return nil, &model.ErrNotFound{}
	}
	if err != nil {
		return nil, err
	}

	todo := &model.TODO{
		ID: id,
	}
	row := s.db.QueryRowContext(ctx, confirm, id)
	err = row.Scan(&todo.Subject, &todo.Description, &todo.CreatedAt, &todo.UpdatedAt)
	if err != nil {
		return nil, err
	}

	return todo, nil
}

// DeleteTODO deletes TODOs on DB by ids.
func (s *TODOService) DeleteTODO(ctx context.Context, ids []int64) error {
	const deleteFmt = `DELETE FROM todos WHERE id IN (?%s)`

	return nil
}
