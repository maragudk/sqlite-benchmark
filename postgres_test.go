//go:build postgres

package sqlite_test

import (
	"database/sql"
	_ "embed"
	"fmt"
	"math/rand"
	"testing"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/maragudk/sqlite-benchmark"
)

func BenchmarkDB_Postgres_ReadPostAndMaybeWriteComment(b *testing.B) {
	db := setupPostgres(b, false)

	b.ResetTimer()

	for _, commentRate := range []float64{0.01, 0.1, 1} {
		b.Run(fmt.Sprintf("comment rate %v", commentRate), func(b *testing.B) {
			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _, err := db.ReadPost(1)
					noErr(b, err)
					if rand.Float64() < commentRate {
						err = db.WriteComment(1, "Love it!", "Great post. :D")
						noErr(b, err)
					}
				}
			})
		})
	}
}

//go:embed postgres.sql
var postgresSchema string

func setupPostgres(tb testing.TB, withMutex bool) *sqlite.DB {
	tb.Helper()

	db, err := sql.Open("pgx", "postgresql://test:123@localhost:5432/benchmark?sslmode=disable")
	noErr(tb, err)

	_, err = db.Exec(postgresSchema)
	noErr(tb, err)

	tb.Cleanup(func() {
		_, err := db.Exec(`drop table comments; drop table posts;`)
		noErr(tb, err)
	})

	newDB := sqlite.NewDB(db, withMutex)

	err = newDB.WritePost("First post!", loremIpsum)
	noErr(tb, err)

	return newDB
}
