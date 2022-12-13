package sqlite_test

import (
	"database/sql"
	_ "embed"
	"fmt"
	"math/rand"
	"path"
	"strconv"
	"testing"

	_ "github.com/mattn/go-sqlite3"

	"github.com/maragudk/sqlite-benchmark"
)

func TestPragma(t *testing.T) {
	t.Run("sets up a new DB", func(t *testing.T) {
		db := setupSQLite(t, false)

		for _, pragma := range []string{"synchronous", "journal_mode", "busy_timeout", "auto_vacuum", "foreign_keys"} {
			t.Log("PRAGMA", pragma, getPragma(db, pragma))
		}
	})
}

func getPragma(db *sqlite.DB, name string) string {
	var s string
	if err := db.DB.QueryRow(`PRAGMA ` + name).Scan(&s); err != nil {
		panic(err)
	}
	return s
}

func BenchmarkSelect1(b *testing.B) {
	db := setupSQLite(b, false)

	b.ResetTimer()

	b.RunParallel(func(pb *testing.PB) {
		for pb.Next() {
			_, err := db.DB.Exec(`select 1`)
			noErr(b, err)
		}
	})
}

func BenchmarkDB_ReadPost(b *testing.B) {
	for _, withMutex := range []bool{false, true} {
		b.Run("mutex "+strconv.FormatBool(withMutex), func(b *testing.B) {
			db := setupSQLite(b, withMutex)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					_, _, err := db.ReadPost(1)
					noErr(b, err)
				}
			})
		})
	}
}

func BenchmarkDB_WritePost(b *testing.B) {
	for _, withMutex := range []bool{false, true} {
		b.Run("mutex "+strconv.FormatBool(withMutex), func(b *testing.B) {
			db := setupSQLite(b, withMutex)

			b.ResetTimer()

			b.RunParallel(func(pb *testing.PB) {
				for pb.Next() {
					err := db.WritePost("Lorem Ipsum, don't you think?", loremIpsum)
					noErr(b, err)
				}
			})
		})
	}
}

func BenchmarkDB_ReadPostAndMaybeWriteComment(b *testing.B) {
	for _, withMutex := range []bool{false, true} {
		b.Run("mutex "+strconv.FormatBool(withMutex), func(b *testing.B) {
			db := setupSQLite(b, withMutex)

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
		})
	}
}

//go:embed sqlite.sql
var sqliteSchema string

func setupSQLite(tb testing.TB, withMutex bool) *sqlite.DB {
	tb.Helper()

	db, err := sql.Open("sqlite3", path.Join(tb.TempDir(), "benchmark.db")+"?_journal=WAL&_timeout=10000&_fk=true")
	noErr(tb, err)

	_, err = db.Exec(sqliteSchema)
	noErr(tb, err)

	newDB := sqlite.NewDB(db, withMutex)

	err = newDB.WritePost("First post!", loremIpsum)
	noErr(tb, err)

	return newDB
}

func noErr(tb testing.TB, err error) {
	tb.Helper()

	if err != nil {
		tb.Fatal("Error is not nil:", err)
	}
}

const loremIpsum = `Lorem ipsum dolor sit amet, consectetur adipiscing elit. Duis eget sapien accumsan, commodo ligula iaculis, pretium ligula. In ac lobortis nulla. Donec lobortis metus sed mauris iaculis euismod. Ut vehicula velit vitae dolor maximus euismod. Nulla vel risus eros. Vivamus porttitor odio eleifend, imperdiet tellus sed, feugiat augue. Donec faucibus eget nunc facilisis gravida. Fusce posuere ac lacus eu rutrum. Vivamus nec nibh sed nisl maximus varius. Pellentesque in placerat eros. Vivamus efficitur in dolor nec eleifend. Proin quis nibh quis enim rutrum posuere. Aliquam odio metus, scelerisque quis massa imperdiet, tincidunt placerat orci. Maecenas posuere, ex vitae porttitor tristique, risus erat bibendum augue, nec tempus eros dolor vitae ipsum.

Praesent euismod dui eu tortor aliquet cursus. Suspendisse non augue a odio placerat tincidunt in commodo mi. Nulla facilisi. Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae; Aliquam risus ipsum, vestibulum at metus id, sodales gravida felis. Integer justo tellus, blandit id orci id, fringilla feugiat sapien. Cras velit neque, pretium quis pharetra vitae, faucibus vel quam. Sed ut ipsum ac nunc blandit condimentum. Curabitur vel accumsan lacus.

Nam felis libero, gravida euismod ultricies sit amet, sodales non mauris. Aliquam rhoncus sed ex id convallis. Vivamus ac hendrerit dui. Maecenas urna nisi, imperdiet eu nunc id, sollicitudin sodales dolor. Proin et metus interdum, suscipit turpis eu, imperdiet libero. Mauris a mattis purus. Ut arcu dolor, pulvinar quis suscipit quis, tempor sed lacus. Donec vitae ipsum sed felis euismod finibus et sit amet diam. Donec rhoncus, nunc ultricies varius faucibus, urna lectus auctor arcu, viverra bibendum mauris lorem et diam. Sed erat massa, vulputate quis fermentum vitae, laoreet ut dui. Fusce lacinia accumsan accumsan. Nulla facilisi. Vestibulum at turpis sed leo dapibus bibendum nec sed ex.

Sed tempor porta orci ac luctus. Integer feugiat vel arcu vel suscipit. Mauris blandit, justo at suscipit tincidunt, odio sem maximus ligula, sit amet placerat risus felis quis turpis. Maecenas eget metus urna. Donec semper eros vel felis dapibus vestibulum. Class aptent taciti sociosqu ad litora torquent per conubia nostra, per inceptos himenaeos. Curabitur aliquet dapibus risus, ac hendrerit lorem lobortis nec. Sed scelerisque blandit elit, ut pharetra magna venenatis nec. Ut venenatis lobortis felis, at tincidunt nulla congue vitae. Sed viverra, ipsum a posuere finibus, ipsum est condimentum odio, sed dictum nisi velit eget libero. Pellentesque velit urna, eleifend et fringilla sed, bibendum ut sapien. Aenean eu velit et ante rutrum interdum in eu mi. Aliquam erat volutpat.

Suspendisse tempor vestibulum ante, mollis accumsan libero viverra ultricies. Fusce sagittis velit vel urna bibendum, at placerat lorem hendrerit. Nulla facilisi. Aenean at felis nisl. Quisque malesuada ultrices est eu lacinia. Curabitur volutpat sem non nulla ultrices accumsan. Etiam interdum elementum ante vel facilisis. Donec vel nisi quis ante ullamcorper consectetur. Vestibulum elit diam, interdum eget eleifend imperdiet, sollicitudin quis mi. Maecenas nec quam id orci tincidunt venenatis id id purus. Vivamus tristique libero quis purus vulputate rutrum. Nunc sed eleifend orci, at mattis nibh. Quisque rhoncus pharetra velit vel congue.`
