package bible

import (
	"database/sql"
	_ "embed"
	"errors"
	"fmt"
	"math/rand"
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
	"matt": "Mat",
	"ac":   "Act", "gn": "Gen", "dt": "Deu", "deut": "Deu",
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
	type bookRow struct{ abbr, name string }
	var all []bookRow
	for rows.Next() {
		var br bookRow
		if err := rows.Scan(&br.abbr, &br.name); err != nil {
			return err
		}
		all = append(all, br)
	}
	if err := rows.Err(); err != nil {
		return err
	}
	lookup := make(map[string]string)
	for _, br := range all {
		lookup[strings.ToLower(br.name)] = br.abbr
		nospace := strings.ReplaceAll(strings.ToLower(br.name), " ", "")
		lookup[nospace] = br.abbr
		if len(br.name) >= 3 {
			key := strings.ToLower(br.name[:3])
			if _, ok := lookup[key]; !ok {
				lookup[key] = br.abbr
			}
		}
	}
	for alias, abbr := range extraAliases {
		lookup[alias] = abbr
	}
	for _, br := range all {
		lookup[strings.ToLower(br.abbr)] = br.abbr
	}
	b.bookLookup = lookup
	return nil
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

func (b *BibleService) Init(targetName string) error {
	dbDir := filepath.Join(xdg.DataHome, targetName)
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

type VersePage struct {
	Total  int     `json:"total"`
	Offset int     `json:"offset"`
	Verses []Verse `json:"verses"`
}

func (b *BibleService) buildRangeSQL(queryStr string, versionID string) (string, []any, error) {
	queryStr = strings.TrimSpace(queryStr)
	if queryStr == "" {
		return "", nil, fmt.Errorf("empty query")
	}

	var refStart, refEnd ref
	var hasEnd bool
	if idx := strings.Index(queryStr, " - "); idx != -1 {
		var err error
		refStart, err = parseRef(queryStr[:idx])
		if err != nil {
			return "", nil, err
		}
		refStart.book = b.normalizeBook(refStart.book)
		refEnd, err = parseRef(queryStr[idx+3:])
		if err != nil {
			return "", nil, err
		}
		refEnd.book = b.normalizeBook(refEnd.book)
		hasEnd = true
	} else {
		var err error
		refStart, err = parseRef(queryStr)
		if err != nil {
			return "", nil, err
		}
		refStart.book = b.normalizeBook(refStart.book)
	}

	base := `
		FROM verses v
		JOIN books b ON b.abbr = v.book
		WHERE v.version = ?`

	if !hasEnd {
		r := refStart
		if r.chapter == 0 {
			return base + " AND v.book = ?", []any{versionID, r.book}, nil
		}
		if r.verseStart == 0 {
			return base + " AND v.book = ? AND v.chapter = ?", []any{versionID, r.book, r.chapter}, nil
		}
		vsEnd := r.verseEnd
		if vsEnd == 0 {
			vsEnd = r.verseStart
		}
		return base + " AND v.book = ? AND v.chapter = ? AND v.verse BETWEEN ? AND ?",
			[]any{versionID, r.book, r.chapter, r.verseStart, vsEnd}, nil
	}

	ordStart, err := b.bookOrd(refStart.book)
	if err != nil {
		return "", nil, err
	}
	ordEnd, err := b.bookOrd(refEnd.book)
	if err != nil {
		return "", nil, err
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

	sqlStr := base + `
		  AND b.ord BETWEEN ? AND ?
		  AND NOT (b.ord = ? AND (v.chapter < ? OR (v.chapter = ? AND v.verse < ?)))
		  AND NOT (b.ord = ? AND (v.chapter > ? OR (v.chapter = ? AND v.verse > ?)))`
	params := []any{
		versionID,
		ordStart, ordEnd,
		ordStart, chStart, chStart, vsStart,
		ordEnd, chEnd, chEnd, vsEnd,
	}
	return sqlStr, params, nil
}

func (b *BibleService) QueryPage(queryStr string, versionID string, offset int, limit int) (VersePage, error) {
	if b.db == nil {
		return VersePage{}, fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(versionID) == "" {
		return VersePage{}, fmt.Errorf("empty version")
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	core, params, err := b.buildRangeSQL(queryStr, versionID)
	if err != nil {
		return VersePage{}, err
	}

	var total int
	if err := b.db.QueryRow("SELECT COUNT(*) "+core, params...).Scan(&total); err != nil {
		return VersePage{}, fmt.Errorf("count error: %w", err)
	}

	selectParams := append(append([]any{}, params...), limit, offset)
	rows, err := b.db.Query(
		"SELECT v.book, v.chapter, v.verse, v.version, v.text "+core+
			" ORDER BY b.ord, v.chapter, v.verse LIMIT ? OFFSET ?",
		selectParams...)
	if err != nil {
		return VersePage{}, fmt.Errorf("query error: %w", err)
	}
	defer rows.Close()

	verses := []Verse{}
	for rows.Next() {
		var vr Verse
		if err := rows.Scan(&vr.Book, &vr.Chapter, &vr.Verse, &vr.Version, &vr.Text); err != nil {
			return VersePage{}, err
		}
		verses = append(verses, vr)
	}
	if err := rows.Err(); err != nil {
		return VersePage{}, err
	}
	return VersePage{Total: total, Offset: offset, Verses: verses}, nil
}

type Collection int

const (
	OldTestament Collection = iota
	NewTestament
	Apocrypha
)

func (b *BibleService) HasCollection(versionID string, col Collection) (bool, error) {
	var minOrd, maxOrd int
	switch col {
	case OldTestament:
		minOrd, maxOrd = 1, 39
	case NewTestament:
		minOrd, maxOrd = 40, 66
	case Apocrypha:
		minOrd, maxOrd = 100, 999
	}
	var exists bool
	err := b.db.QueryRow(`
		SELECT EXISTS(
			SELECT 1 FROM verses v
			JOIN books b ON b.abbr = v.book
			WHERE v.version = ? AND b.ord BETWEEN ? AND ?
		)`,
		versionID, minOrd, maxOrd,
	).Scan(&exists)
	if err != nil {
		return false, err
	}
	return exists, nil
}

func (b *BibleService) Search(versionID string, book string, phrase string, offset int, limit int) (VersePage, error) {
	if b.db == nil {
		return VersePage{}, fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(versionID) == "" {
		return VersePage{}, fmt.Errorf("empty version")
	}
	if strings.TrimSpace(phrase) == "" {
		return VersePage{}, fmt.Errorf("empty search phrase")
	}
	if offset < 0 {
		offset = 0
	}
	if limit <= 0 {
		limit = 50
	}
	if limit > 200 {
		limit = 200
	}

	pattern := "%" + strings.ReplaceAll(phrase, "%", "\\%") + "%"

	var (
		countSQL  string
		selectSQL string
		params    []any
	)

	if book == "" {
		countSQL = `
			SELECT COUNT(*) FROM verses v
			JOIN books b ON b.abbr = v.book
			WHERE v.version = ? AND v.text LIKE ? ESCAPE '\'`
		selectSQL = `
			SELECT v.book, v.chapter, v.verse, v.version, v.text
			FROM verses v
			JOIN books b ON b.abbr = v.book
			WHERE v.version = ? AND v.text LIKE ? ESCAPE '\'
			ORDER BY b.ord, v.chapter, v.verse
			LIMIT ? OFFSET ?`
		params = []any{versionID, pattern}
	} else {
		normalizedBook := b.normalizeBook(book)
		countSQL = `
			SELECT COUNT(*) FROM verses v
			JOIN books b ON b.abbr = v.book
			WHERE v.version = ? AND v.book = ? AND v.text LIKE ? ESCAPE '\'`
		selectSQL = `
			SELECT v.book, v.chapter, v.verse, v.version, v.text
			FROM verses v
			JOIN books b ON b.abbr = v.book
			WHERE v.version = ? AND v.book = ? AND v.text LIKE ? ESCAPE '\'
			ORDER BY b.ord, v.chapter, v.verse
			LIMIT ? OFFSET ?`
		params = []any{versionID, normalizedBook, pattern}
	}

	var total int
	if err := b.db.QueryRow(countSQL, params...).Scan(&total); err != nil {
		return VersePage{}, fmt.Errorf("search count error: %w", err)
	}

	rows, err := b.db.Query(selectSQL, append(append([]any{}, params...), limit, offset)...)
	if err != nil {
		return VersePage{}, fmt.Errorf("search error: %w", err)
	}
	defer rows.Close()

	verses := []Verse{}
	for rows.Next() {
		var vr Verse
		if err := rows.Scan(&vr.Book, &vr.Chapter, &vr.Verse, &vr.Version, &vr.Text); err != nil {
			return VersePage{}, err
		}
		verses = append(verses, vr)
	}
	if err := rows.Err(); err != nil {
		return VersePage{}, err
	}
	return VersePage{Total: total, Offset: offset, Verses: verses}, nil
}

func (b *BibleService) GetChapterCount(versionID string, book string) (int, error) {
	if b.db == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	normalizedBook := b.normalizeBook(book)
	var count int
	err := b.db.QueryRow(`
		SELECT COUNT(DISTINCT chapter) FROM verses
		WHERE version = ? AND book = ?`,
		versionID, normalizedBook,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("chapter count error: %w", err)
	}
	return count, nil
}

func (b *BibleService) GetVerseCount(versionID string, book string, chapter int) (int, error) {
	if b.db == nil {
		return 0, fmt.Errorf("database not initialized")
	}
	normalizedBook := b.normalizeBook(book)
	var count int
	err := b.db.QueryRow(`
		SELECT COUNT(*) FROM verses
		WHERE version = ? AND book = ? AND chapter = ?`,
		versionID, normalizedBook, chapter,
	).Scan(&count)
	if err != nil {
		return 0, fmt.Errorf("verse count error: %w", err)
	}
	return count, nil
}

func (b *BibleService) RandomVerse(versionID string, book string) (Verse, error) {
	if b.db == nil {
		return Verse{}, fmt.Errorf("database not initialized")
	}
	if strings.TrimSpace(versionID) == "" {
		return Verse{}, fmt.Errorf("empty version")
	}

	var countSQL, selectSQL string
	var params []any

	if book == "" {
		countSQL = `SELECT COUNT(*) FROM verses WHERE version = ?`
		params = []any{versionID}
	} else {
		normalizedBook := b.normalizeBook(book)
		countSQL = `SELECT COUNT(*) FROM verses WHERE version = ? AND book = ?`
		params = []any{versionID, normalizedBook}
	}

	var total int
	if err := b.db.QueryRow(countSQL, params...).Scan(&total); err != nil {
		return Verse{}, fmt.Errorf("random verse count error: %w", err)
	}
	if total == 0 {
		return Verse{}, fmt.Errorf("no verses found")
	}

	offset := rand.Intn(total)

	if book == "" {
		selectSQL = `
			SELECT book, chapter, verse, version, text FROM verses
			WHERE version = ?
			LIMIT 1 OFFSET ?`
		params = []any{versionID, offset}
	} else {
		normalizedBook := b.normalizeBook(book)
		selectSQL = `
			SELECT book, chapter, verse, version, text FROM verses
			WHERE version = ? AND book = ?
			LIMIT 1 OFFSET ?`
		params = []any{versionID, normalizedBook, offset}
	}

	var vr Verse
	err := b.db.QueryRow(selectSQL, params...).Scan(&vr.Book, &vr.Chapter, &vr.Verse, &vr.Version, &vr.Text)
	if err != nil {
		return Verse{}, fmt.Errorf("random verse error: %w", err)
	}
	return vr, nil
}

type MultiVersionVerse struct {
	Book    string            `json:"book"`
	Chapter int               `json:"chapter"`
	Verse   int               `json:"verse"`
	Texts   map[string]string `json:"texts"`
}

func (b *BibleService) QueryMultiVersion(queryStr string, versionIDs []string) ([]MultiVersionVerse, error) {
	if b.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	if len(versionIDs) == 0 {
		return nil, fmt.Errorf("no versions specified")
	}

	first, err := b.QueryPage(queryStr, versionIDs[0], 0, 1000)
	if err != nil {
		return nil, err
	}

	type loc struct {
		book    string
		chapter int
		verse   int
	}
	locMap := make(map[loc]*MultiVersionVerse)
	order := make([]loc, 0, len(first.Verses))

	for _, v := range first.Verses {
		l := loc{v.Book, v.Chapter, v.Verse}
		mv := &MultiVersionVerse{
			Book:    v.Book,
			Chapter: v.Chapter,
			Verse:   v.Verse,
			Texts:   make(map[string]string),
		}
		mv.Texts[versionIDs[0]] = v.Text
		locMap[l] = mv
		order = append(order, l)
	}

	for _, vid := range versionIDs[1:] {
		page, err := b.QueryPage(queryStr, vid, 0, 1000)
		if err != nil {
			continue
		}
		for _, v := range page.Verses {
			l := loc{v.Book, v.Chapter, v.Verse}
			if mv, ok := locMap[l]; ok {
				mv.Texts[vid] = v.Text
			}
		}
	}

	result := make([]MultiVersionVerse, len(order))
	for i, l := range order {
		result[i] = *locMap[l]
	}
	return result, nil
}

func (b *BibleService) GetBooksByCollection(col Collection) ([]Book, error) {
	if b.db == nil {
		return nil, fmt.Errorf("database not initialized")
	}
	var minOrd, maxOrd int
	switch col {
	case OldTestament:
		minOrd, maxOrd = 1, 39
	case NewTestament:
		minOrd, maxOrd = 40, 66
	case Apocrypha:
		minOrd, maxOrd = 100, 999
	}
	rows, err := b.db.Query(`
		SELECT abbr, name, ord FROM books
		WHERE ord BETWEEN ? AND ?
		ORDER BY ord`,
		minOrd, maxOrd,
	)
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

func (b *BibleService) bookOrd(abbr string) (int, error) {
	var ord int
	err := b.db.QueryRow("SELECT ord FROM books WHERE abbr = ?", abbr).Scan(&ord)
	if err != nil {
		return 0, fmt.Errorf("unknown book %q", abbr)
	}
	return ord, nil
}

type PassageRange struct {
	Book     string `json:"book"`
	Display  string `json:"display"`
	StartPos int    `json:"startPos"`
	EndPos   int    `json:"endPos"`
}

func encodePos(ord, chapter, verse int) int {
	return ord*1_000_000 + chapter*1_000 + verse
}

func (b *BibleService) bookInfo(abbr string) (string, int, error) {
	var name string
	var ord int
	err := b.db.QueryRow("SELECT name, ord FROM books WHERE abbr = ?", abbr).Scan(&name, &ord)
	if err != nil {
		return "", 0, fmt.Errorf("unknown book %q", abbr)
	}
	return name, ord, nil
}

func renderRef(bookName string, r ref) string {
	if r.chapter == 0 {
		return bookName
	}
	if r.verseStart == 0 {
		return fmt.Sprintf("%s %d", bookName, r.chapter)
	}
	if r.verseEnd != 0 && r.verseEnd != r.verseStart {
		return fmt.Sprintf("%s %d:%d-%d", bookName, r.chapter, r.verseStart, r.verseEnd)
	}
	return fmt.Sprintf("%s %d:%d", bookName, r.chapter, r.verseStart)
}

func (b *BibleService) ResolveRange(refStr string) (PassageRange, error) {
	if b.db == nil {
		return PassageRange{}, fmt.Errorf("database not initialized")
	}
	refStr = strings.TrimSpace(refStr)
	if refStr == "" {
		return PassageRange{}, fmt.Errorf("empty reference")
	}

	if idx := strings.Index(refStr, " - "); idx != -1 {
		left, err := parseRef(refStr[:idx])
		if err != nil {
			return PassageRange{}, err
		}
		right, err := parseRef(refStr[idx+3:])
		if err != nil {
			return PassageRange{}, err
		}
		left.book = b.normalizeBook(left.book)
		right.book = b.normalizeBook(right.book)

		leftName, leftOrd, err := b.bookInfo(left.book)
		if err != nil {
			return PassageRange{}, err
		}
		rightName, rightOrd, err := b.bookInfo(right.book)
		if err != nil {
			return PassageRange{}, err
		}

		chStart, vsStart := left.chapter, left.verseStart
		if chStart == 0 {
			chStart = 1
		}
		if vsStart == 0 {
			vsStart = 1
		}
		chEnd := right.chapter
		if chEnd == 0 {
			chEnd = 999
		}
		vsEnd := right.verseEnd
		if vsEnd == 0 {
			if right.verseStart != 0 {
				vsEnd = right.verseStart
			} else {
				vsEnd = 999
			}
		}

		pr := PassageRange{
			Book:     left.book,
			Display:  renderRef(leftName, left) + " - " + renderRef(rightName, right),
			StartPos: encodePos(leftOrd, chStart, vsStart),
			EndPos:   encodePos(rightOrd, chEnd, vsEnd),
		}
		if pr.StartPos > pr.EndPos {
			return PassageRange{}, fmt.Errorf("reference range is reversed: %q", refStr)
		}
		return pr, nil
	}

	r, err := parseRef(refStr)
	if err != nil {
		return PassageRange{}, err
	}
	r.book = b.normalizeBook(r.book)
	name, ord, err := b.bookInfo(r.book)
	if err != nil {
		return PassageRange{}, err
	}

	chStart, chEnd := r.chapter, r.chapter
	vsStart, vsEnd := r.verseStart, r.verseEnd
	if r.chapter == 0 {
		chStart, chEnd = 1, 999
		vsStart, vsEnd = 1, 999
	} else if r.verseStart == 0 {
		vsStart, vsEnd = 1, 999
	} else if vsEnd == 0 {
		vsEnd = vsStart
	}
	if vsEnd < vsStart {
		return PassageRange{}, fmt.Errorf("verse range is reversed: %q", refStr)
	}

	return PassageRange{
		Book:     r.book,
		Display:  renderRef(name, r),
		StartPos: encodePos(ord, chStart, vsStart),
		EndPos:   encodePos(ord, chEnd, vsEnd),
	}, nil
}
