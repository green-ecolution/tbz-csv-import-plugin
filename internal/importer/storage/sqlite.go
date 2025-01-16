package storage

import (
	"context"
	"embed"

	"github.com/green-ecolution/tbz-csv-import-plugin/internal/entities"
	"github.com/jmoiron/sqlx"
	"github.com/pressly/goose/v3"
)

type ImportRepository interface {
	GetAllTrees(ctx context.Context) ([]entities.Tree, error)
	DeleteTreesByID(ctx context.Context, treeID []entities.TreeID) error
	CreateTrees(ctx context.Context, trees []*entities.Tree) error
}

type ImportRepositoryDB struct {
	db *sqlx.DB
}

type ImportRepositoryTx struct {
	db *sqlx.Tx
}

func NewImportRepositoryDB(db *sqlx.DB) *ImportRepositoryDB {
	return &ImportRepositoryDB{
		db,
	}
}

func NewImportRepositoryTx(db *sqlx.Tx) *ImportRepositoryTx {
	return &ImportRepositoryTx{
		db,
	}
}

//go:embed migrations/*.sql
var migrations embed.FS

func (r *ImportRepositoryDB) Setup() error {
	sqlDB := r.db.DB
	goose.SetBaseFS(migrations)
	if err := goose.SetDialect("sqlite3"); err != nil {
		return err
	}

	if err := goose.Up(sqlDB, "migrations"); err != nil {
		return err
	}

	return nil
}

func (r *ImportRepositoryDB) WithTx(ctx context.Context, fn func(context.Context, *ImportRepositoryTx) error) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	if err := fn(ctx, NewImportRepositoryTx(tx)); err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	if err := tx.Commit(); err != nil {
		return err
	}

	return nil
}

const (
	getAllQuery = "SELECT * FROM trees"
	deleteQuery = "DELETE FROM trees WHERE tree_id IN (?)"
	createQuery = "INSERT INTO trees (tree_number, species, area, planting_year, street, latitude, longitude) VALUES (:tree_number, :species, :area, :planting_year, :street, :latitude, :longitude)"
	updateQuery = "UPDATE trees SET tree_number = :tree_number, species = :species, area = :area, planting_year = :planting_year, street = :street, latitude = :latitude, longitude = :longitude WHERE tree_id = :tree_id"
)

func (r *ImportRepositoryDB) GetAllTrees(ctx context.Context) ([]entities.Tree, error) {
	var trees []entities.Tree
	err := r.db.SelectContext(ctx, &trees, getAllQuery)
	return trees, err
}

func (r *ImportRepositoryTx) GetAllTrees(ctx context.Context) ([]entities.Tree, error) {
	var trees []entities.Tree
	err := r.db.SelectContext(ctx, &trees, getAllQuery)
	if err != nil {
		return nil, err
	}

	return trees, nil
}

func (r *ImportRepositoryDB) DeleteTreesByID(ctx context.Context, treeID []entities.TreeID) error {
	_, err := r.db.ExecContext(ctx, deleteQuery, treeID)
	return err
}

func (r *ImportRepositoryTx) DeleteTreesByID(ctx context.Context, treeID []entities.TreeID) error {
	_, err := r.db.ExecContext(ctx, deleteQuery, treeID)
	return err
}

func (r *ImportRepositoryDB) CreateTrees(ctx context.Context, trees []*entities.Tree) error {
	_, err := r.db.NamedExecContext(ctx, createQuery, trees)
	return err
}

func (r *ImportRepositoryTx) CreateTrees(ctx context.Context, trees []*entities.Tree) error {
	_, err := r.db.NamedExecContext(ctx, createQuery, trees)
	return err
}

func (r *ImportRepositoryDB) UpdateTrees(ctx context.Context, trees []*entities.Tree) error {
	_, err := r.db.NamedExecContext(ctx, updateQuery, trees)
	return err
}

func (r *ImportRepositoryTx) UpdateTrees(ctx context.Context, trees []*entities.Tree) error {
	_, err := r.db.NamedExecContext(ctx, updateQuery, trees)
	return err
}

func (r *ImportRepositoryDB) AddImport(ctx context.Context, i entities.Import, treeIDs []entities.TreeID) error {
	tx, err := r.db.BeginTxx(ctx, nil)
	if err != nil {
		return err
	}

	importRow, err := tx.ExecContext(ctx, "INSERT INTO imports (created_at, user_id, raw_csv) VALUES (datetime('now'), ?, ?)", i.UserID, i.RawCSV)
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	importID, err := importRow.LastInsertId()
	if err != nil {
		if err := tx.Rollback(); err != nil {
			return err
		}
		return err
	}

	for _, treeID := range treeIDs {
		if _, err := tx.ExecContext(ctx, "INSERT INTO import_tree (import_id, tree_id) VALUES (?, ?)", importID, treeID); err != nil {
			if err := tx.Rollback(); err != nil {
				return err
			}
			return err
		}
	}

	return err
}
