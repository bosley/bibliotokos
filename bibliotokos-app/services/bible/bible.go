package bible

import (
	"bibliotokos/platform"
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"os"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"

	"github.com/adrg/xdg"
	_ "modernc.org/sqlite"
)

//go:embed bible.db
var embeddedDB []byte

type Book struct {
	Abbr string `json:"abbr"`
	Name string `json:"name"`
	Ord  int    `json:"ord"`
}

type Version struct {
	ID   string `json:"id"`
	Name string `json:"name"`
	Lang string `json:"lang"`
}

type Verse struct {
	Book    string `json:"book"`
	Chapter int    `json:"chapter"`
	Verse   int    `json:"verse"`
	Version string `json:"version"`
	Text    string `json:"text"`
}

type BibleService struct {
	db         *sql.DB
	bookLookup map[string]string
}

var refRE = regexp.MustCompile(`^([1-9]?[A-Za-z][A-Za-z0-9]*)(?:\s+(\d+)(?::(\d+)(?:-(\d+))?)?)?$`)

var extraAliases = map[string]string{
	"psalm": "Psa", "song": "Sol", "songs": "Sol", "canticles": "Sol",
	"canticle": "Sol", "sos": "Sol", "sg": "Sol", "eccl": "Ecc",
	"qoh": "Ecc", "zech": "Zac", "phil": "Phi", "phm": "Plm",
	"phlm": "Plm", "jas": "Jam", "jude": "Jde",
	"1sam": "Sa1", "2sam": "Sa2", "1kgs": "Kg1", "2kgs": "Kg2",
	"1ki": "Kg1", "2ki": "Kg2", "1chr": "Ch1", "2chr": "Ch2",
	"1chron": "Ch1", "2chron": "Ch2", "1cor": "Co1", "2cor": "Co2",
	"1thess": "Th1", "2thess": "Th2", "1tim": "Ti1", "2tim": "Ti2",
	"1pet": "Pe1", "2pet": "Pe2", "1john": "Jo1", "2john": "Jo2",
	"3john": "Jo3", "1jn": "Jo1", "2jn": "Jo2", "3jn": "Jo3",
	"apoc": "Rev", "mk": "Mar", "lk": "Luk", "jn": "Joh", "mt": "Mat",
	"ac": "Act", "gn": "Gen", "dt": "Deu", "deut": "Deu",
	"ex": "Exo", "exod": "Exo", "lv": "Lev", "nm": "Num",
	"jsh": "Jos", "jg": "Jdg", "jdgs": "Jdg", "ru": "Rut",
	"ezk": "Eze", "zeph": "Zep", "obad": "Oba",
	"lamentations": "Lam", "philemon": "Plm",
}

func (b *BibleService) buildBookLookup() error {
	rows, err := b.db.Query("SELECT abbr, name FROM books")
	if err != nil {
		return err
	}
	defer rows.Close()
	lookup := make(map[string]string)
	for rows.Next() {
		var abbr, name string
		if err := rows.Scan(&abbr, &name); err != nil {
			return err
		}
		lookup[strings.ToLower(abbr)] = abbr
		lookup[strings.ToLower(name)] = abbr
		nospace := strings.ReplaceAll(strings.ToLower(name), " ", "")
		lookup[nospace] = abbr
		if len(name) >= 3 {
			lookup[strings.ToLower(name[:3])] = abbr
		}
	}
	for alias, abbr := range extraAliases {
		lookup[alias] = abbr
	}
	b.bookLookup = lookup
	return rows.Err()
}

func (b *BibleService) normalizeBook(raw string) string {
	key := strings.ToLower(strings.TrimSpace(raw))
	if canon, ok := b.bookLookup[key]; ok {
		return canon
	}
	lower := strings.ToLower(raw)
	return strings.ToUpper(lower[:1]) + lower[1:]
}

type ref struct {
	book       string
	chapter    int
	verseStart int
	verseEnd   int
}

func parseRef(token string) (ref, error) {
	token = strings.TrimSpace(token)
	m := refRE.FindStringSubmatch(token)
	if m == nil {
		return ref{}, fmt.Errorf("cannot parse reference: %q", token)
	}
	lower := strings.ToLower(m[1])
	r := ref{book: strings.ToUpper(lower[:1]) + lower[1:]}
	if m[2] != "" {
		r.chapter, _ = strconv.Atoi(m[2])
	}
	if m[3] != "" {
		r.verseStart, _ = strconv.Atoi(m[3])
		r.verseEnd = r.verseStart
	}
	if m[4] != "" {
		r.verseEnd, _ = strconv.Atoi(m[4])
	}
	return r, nil
}

func (b *BibleService) Init() error {
	dbDir := filepath.Join(xdg.DataHome, platform.GetInstallAppName())
	if err := os.MkdirAll(dbDir, 0755); err != nil {
		return fmt.Errorf("create data dir: %w", err)
	}
	dbPath := filepath.Join(dbDir, "bible.db")

	if _, err := os.Stat(dbPath); errors.Is(err, os.ErrNotExist) {
		if err := os.WriteFile(dbPath, embeddedDB, 0644); err != nil {
			return fmt.Errorf("write db: %w", err)
		}
	}

	db, err := sql.Open("sqlite", "file:"+dbPath+"?mode=ro")
	if err != nil {
		return fmt.Errorf("open db: %w", err)
	}
	b.db = db
	return b.buildBookLookup()
}

func (b *BibleService) GetBooks() ([]Book, error) {
	rows, err := b.db.Query("SELECT abbr, name, ord FROM books ORDER BY ord")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var books []Book
	for rows.Next() {
		var bk Book
		if err := rows.Scan(&bk.Abbr, &bk.Name, &bk.Ord); err != nil {
			return nil, err
		}
		books = append(books, bk)
	}
	return books, rows.Err()
}

func (b *BibleService) GetVersions() ([]Version, error) {
	rows, err := b.db.Query("SELECT id, name, lang FROM versions ORDER BY id")
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	var versions []Version
	for rows.Next() {
		var v Version
		if err := rows.Scan(&v.ID, &v.Name, &v.Lang); err != nil {
			return nil, err
		}
		versions = append(versions, v)
	}
	return versions, rows.Err()
}

func placeholders(n int) string {
	if n == 0 {
		return ""
	}
	return strings.TrimRight(strings.Repeat("?,", n), ",")
}

func anySlice[T any](s []T) []any {
	out := make([]any, len(s))
	for i, v := range s {
		out[i] = v
	}
	return out
}

func (b *BibleService) Query(queryStr string, selectedBooks []string, selectedVersions []string) ([]Verse, error) {
	if b.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}

	queryStr = strings.TrimSpace(queryStr)
	if queryStr == "" {
		return nil, fmt.Errorf("empty query")
	}

	if len(selectedVersions) == 0 {
		vrows, err := b.db.Query("SELECT id FROM versions ORDER BY id")
		if err != nil {
			return nil, err
		}
		for vrows.Next() {
			var id string
			if err := vrows.Scan(&id); err != nil {
				vrows.Close()
				return nil, err
			}
			selectedVersions = append(selectedVersions, id)
		}
		vrows.Close()
		if err := vrows.Err(); err != nil {
			return nil, err
		}
	}

	var refStart, refEnd ref
	var hasEnd bool
	if idx := strings.Index(queryStr, " - "); idx != -1 {
		left := queryStr[:idx]
		right := queryStr[idx+3:]
		var err error
		refStart, err = parseRef(left)
		if err != nil {
			return nil, err
		}
		refStart.book = b.normalizeBook(refStart.book)
		refEnd, err = parseRef(right)
		if err != nil {
			return nil, err
		}
		refEnd.book = b.normalizeBook(refEnd.book)
		hasEnd = true
	} else {
		var err error
		refStart, err = parseRef(queryStr)
		if err != nil {
			return nil, err
		}
		refStart.book = b.normalizeBook(refStart.book)
	}

	verPH := placeholders(len(selectedVersions))
	verParams := anySlice(selectedVersions)

	bookFilterSQL := ""
	var bookParams []any
	if len(selectedBooks) > 0 {
		bookFilterSQL = fmt.Sprintf(" AND v.book IN (%s)", placeholders(len(selectedBooks)))
		bookParams = anySlice(selectedBooks)
	}

	var sqlStr string
	var params []any

	if !hasEnd {
		r := refStart
		if r.chapter == 0 {
			sqlStr = fmt.Sprintf(`
				SELECT v.book, v.chapter, v.verse, v.version, v.text
				FROM verses v
				JOIN books b ON b.abbr = v.book
				WHERE v.book = ? AND v.version IN (%s)%s
				ORDER BY v.chapter, v.verse, v.version
			`, verPH, bookFilterSQL)
			params = append(params, r.book)
			params = append(params, verParams...)
			params = append(params, bookParams...)
		} else if r.verseStart == 0 {
			sqlStr = fmt.Sprintf(`
				SELECT v.book, v.chapter, v.verse, v.version, v.text
				FROM verses v
				WHERE v.book = ? AND v.chapter = ? AND v.version IN (%s)%s
				ORDER BY v.verse, v.version
			`, verPH, bookFilterSQL)
			params = append(params, r.book, r.chapter)
			params = append(params, verParams...)
			params = append(params, bookParams...)
		} else {
			vsEnd := r.verseEnd
			if vsEnd == 0 {
				vsEnd = r.verseStart
			}
			sqlStr = fmt.Sprintf(`
				SELECT v.book, v.chapter, v.verse, v.version, v.text
				FROM verses v
				WHERE v.book = ? AND v.chapter = ? AND v.verse BETWEEN ? AND ?
				  AND v.version IN (%s)%s
				ORDER BY v.verse, v.version
			`, verPH, bookFilterSQL)
			params = append(params, r.book, r.chapter, r.verseStart, vsEnd)
			params = append(params, verParams...)
			params = append(params, bookParams...)
		}
	} else {
		ordStart, err := b.bookOrd(refStart.book)
		if err != nil {
			return nil, err
		}
		ordEnd, err := b.bookOrd(refEnd.book)
		if err != nil {
			return nil, err
		}
		chStart := refStart.chapter
		if chStart == 0 {
			chStart = 1
		}
		vsStart := refStart.verseStart
		if vsStart == 0 {
			vsStart = 1
		}
		chEnd := refEnd.chapter
		if chEnd == 0 {
			chEnd = 999999
		}
		vsEnd := refEnd.verseEnd
		if vsEnd == 0 {
			if refEnd.verseStart != 0 {
				vsEnd = refEnd.verseStart
			} else {
				vsEnd = 999999
			}
		}

		sqlStr = fmt.Sprintf(`
			SELECT v.book, v.chapter, v.verse, v.version, v.text
			FROM verses v
			JOIN books b ON b.abbr = v.book
			WHERE v.version IN (%s)
			  AND b.ord BETWEEN ? AND ?
			  AND NOT (b.ord = ? AND (v.chapter < ? OR (v.chapter = ? AND v.verse < ?)))
			  AND NOT (b.ord = ? AND (v.chapter > ? OR (v.chapter = ? AND v.verse > ?)))%s
			ORDER BY b.ord, v.chapter, v.verse, v.version
		`, verPH, bookFilterSQL)
		params = append(params, verParams...)
		params = append(params,
			ordStart, ordEnd,
			ordStart, chStart, chStart, vsStart,
			ordEnd, chEnd, chEnd, vsEnd,
		)
		params = append(params, bookParams...)
	}

	rows, err := b.db.Query(sqlStr, params...)
	if err != nil {
		return nil, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	var verses []Verse
	for rows.Next() {
		var vr Verse
		if err := rows.Scan(&vr.Book, &vr.Chapter, &vr.Verse, &vr.Version, &vr.Text); err != nil {
			return nil, err
		}
		verses = append(verses, vr)
	}
	return verses, rows.Err()
}

func (b *BibleService) bookOrd(abbr string) (int, error) {
	var ord int
	err := b.db.QueryRow("SELECT ord FROM books WHERE abbr = ?", abbr).Scan(&ord)
	if err != nil {
		return 0, fmt.Errorf("unknown book %q", abbr)
	}
	return ord, nil
}
