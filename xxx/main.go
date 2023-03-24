package main

import (
	"context"
	"database/sql"
	_ "embed"
	"fmt"

	sq "github.com/Masterminds/squirrel"
	_ "github.com/go-sql-driver/mysql"
	"github.com/tikinang/scuffolding/g"
)

//go:embed updateArtTitle.sql
var updateArtTitleSql string

func main() {
	ctx := context.Background()

	db, err := sql.Open("mysql", "root:rootik@tcp(localhost:3306)/icbaat")
	if err != nil {
		panic(err)
	}
	defer db.Close()

	tx, err := db.BeginTx(ctx, nil)
	if err != nil {
		panic(err)
	}
	defer tx.Rollback()

	options := []g.Option[Cnt]{
		WithQueryVisible(),
	}

	cnt := g.ApplyOptions(options...)
	fmt.Println(cnt.QueryVisible, cnt.Whop)

	cnt = g.ApplyOptionsOnDefault(
		Cnt{
			Whop: "xxx",
		},
		options...,
	)
	fmt.Println(cnt.QueryVisible, cnt.Whop)

	b := sq.
		Select("id", "title", "description", "MATCH(description) AGAINST ('sunset' WITH QUERY EXPANSION) AS relevance").
		From("art").
		Where(
			sq.Expr("MATCH(description) AGAINST ('sunset' WITH QUERY EXPANSION)"),
		).
		OrderBy("relevance DESC")

	docs, err := g.FindManyDocs[ArtSearch](ctx, tx, b)
	if err != nil {
		panic(err)
	}
	for _, doc := range docs {
		fmt.Println(doc.Id, doc.Description, doc.Relevance)
	}

	b = b.Limit(1)
	docNull, err := g.FoundToNull(g.FindOneDoc[ArtSearch](ctx, tx, b))
	if err != nil {
		panic(err)
	}
	if doc, filled := docNull.Ask(); filled {
		fmt.Println(doc.Id, doc.Description, doc.Relevance)
	}

	rowsAffected, err := g.Exec(ctx, tx, sq.Expr(updateArtTitleSql))
	if err != nil {
		panic(err)
	}
	fmt.Println("update rows affected:", rowsAffected)

	bi := sq.
		Insert("art").
		Columns("customer_id", "title", "content", "description").
		Values(
			sq.Expr("(SELECT id FROM customer WHERE name = 'Matthew')"),
			"plop.jpeg",
			[]byte{'x', 'y', 'z'},
			"Plop plop plo-pity plop.",
		)

	fmt.Println(bi.MustSql())
	id, err := g.InsertValuesReturningId(ctx, tx, bi)
	if err != nil {
		panic(err)
	}
	fmt.Println("insert id:", id)

	newArts := []ArtInsert{
		{
			CustomerId:  "737a2de5-c5ac-11ed-bf5f-0242ac110002",
			Title:       "plop.jpeg",
			Content:     []byte{'x', 'y', 'z'},
			Description: "Plop plop plo-pity plop.",
		},
		{
			CustomerId:  "737b7bd2-c5ac-11ed-bf5f-0242ac110002",
			Title:       "plop.jpeg",
			Content:     []byte{'x', 'y', 'z', '1'},
			Description: "Plop plop plo-pity plop.",
		},
	}

	id, err = g.InsertOneDocReturningId(ctx, tx, newArts[0])
	if err != nil {
		panic(err)
	}
	fmt.Println("insert id:", id)

	err = g.InsertOneDoc(ctx, tx, newArts[1])
	if err != nil {
		panic(err)
	}

	err = g.InsertManyDocs(ctx, tx, newArts...)
	if err != nil {
		panic(err)
	}

	if err := tx.Commit(); err != nil {
		panic(err)
	}
}

type Cnt struct {
	QueryVisible bool
	Whop         string
}

func WithQueryVisible() g.Option[Cnt] {
	return func(cnt *Cnt) {
		cnt.QueryVisible = true
	}
}

func WithWhop(whop string) g.Option[Cnt] {
	return func(cnt *Cnt) {
		cnt.Whop = whop
	}
}

type ArtSearch struct {
	Id          string
	Title       string
	Description string
	Relevance   float64
}

func (r *ArtSearch) ScanPointers() []any {
	return []any{
		&r.Id,
		&r.Title,
		&r.Description,
		&r.Relevance,
	}
}

type ArtInsert struct {
	CustomerId  string
	Title       string
	Content     []byte // FIXME(mpavlicek): Some kind of io.Reader.
	Description string
}

func (r ArtInsert) Table() string {
	return "art"
}

func (r ArtInsert) Columns() []string {
	return []string{
		"customer_id",
		"title",
		"content",
		"description",
	}
}

func (r ArtInsert) Values() []any {
	return []any{
		r.CustomerId,
		r.Title,
		r.Content,
		r.Description,
	}
}
