//go:build postgres

package sqlite_test

import (
	"database/sql"
	_ "embed"
	"fmt"
	"math/rand"
	"os/exec"
	"testing"
	"time"

	_ "github.com/jackc/pgx/v4/stdlib"

	"github.com/maragudk/sqlite-benchmark"
)

func BenchmarkDB_Postgres_ReadPostAndMaybeWriteComment(b *testing.B) {
	for _, commentRate := range []float64{0.01, 0.1, 1} {
		b.Run(fmt.Sprintf("comment rate %v", commentRate), func(b *testing.B) {
			db := setupPostgres(b, false)

			b.ResetTimer()

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

	cmd := exec.Command("docker", "compose", "up", "-d")
	if output, err := cmd.CombinedOutput(); err != nil {
		tb.Fatal(string(output))
	}

	time.Sleep(3 * time.Second)

	_, err = db.Exec(postgresSchema)
	noErr(tb, err)

	tb.Cleanup(func() {
		cmd := exec.Command("docker", "compose", "down")
		if output, err := cmd.CombinedOutput(); err != nil {
			tb.Fatal(string(output))
		}
	})

	newDB := sqlite.NewDB(db, withMutex)

	err = newDB.WritePost("First post!", loremIpsum)
	noErr(tb, err)

	return newDB
}
