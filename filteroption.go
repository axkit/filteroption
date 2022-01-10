// The package filteroption provides helpers to simplify catching & normalization
// HTTP GET url attributes like "/customers&sortBy=name&pageSize=20"
//
// Supported URL parameters:
//  pageNumber: 0...
//  pageSize: 0...
//  sortBy: "name" or "-name"
//  show: "" - means not deleted rows only, "deleted" - deleted only, "all" - all.
//  Lang: "", "en", "fr"...
//  Download: true or false
package filteroption

import "sort"

type ShowDeletedRule string

const (
	HideDeleted     ShowDeletedRule = ""
	ShowDeletedOnly ShowDeletedRule = "deleted"
	ShowAll         ShowDeletedRule = "all"
)

// DefaultPageSize holds maximum amount of rows to be returned if PageSize not specified.
var DefaultPageSize = 10

// FilterOption holds input
type FilterOption struct {
	// PageNumber start from 0.
	PageNumber int `schema:"pageNumber" param:"pageNumber"`
	// PageSize holds
	PageSize int             `schema:"pageSize" param:"pageSize"`
	SortBy   string          `schema:"sortBy" param:"sortBy"`
	Show     ShowDeletedRule `schema:"show" param:"show"`
	Download bool            `schema:"download" param:"download"`
	Lang     string          `schema:"lang" param:"lang"`
	sortBy   string
	desc     bool
	total    int
	ids      []int
}

// Add adds object id to the pre-resultset.
func (fo *FilterOption) Add(id int) {
	fo.ids = append(fo.ids, id)
	fo.total++
}

// IDs returns stored pre-resultset.
func (fo *FilterOption) IDs() []int {
	return fo.ids
}

// ApplyDefaults assign defaults values to FilterOption. The method is supposed
// to be called after parsing and applying data from URL to FilterOption{}.
//
// Function call modifies PageNumber and PageSize if required.
func (fo *FilterOption) ApplyDefaults() *FilterOption {

	if fo.Download {
		fo.PageNumber = 0
		fo.PageSize = 0
	} else {
		if fo.PageSize == 0 {
			fo.PageSize = DefaultPageSize
		}
	}

	fo.sortBy = fo.SortBy

	if len(fo.SortBy) > 0 && []rune(fo.SortBy)[0] == '-' {
		fo.desc = true
		fo.sortBy = string([]rune(fo.SortBy)[1:])
		fo.SortBy = fo.sortBy
	}
	return fo
}

// SetLang assigns language to the Lang attribute.
func (fo *FilterOption) SetLang(lang string) *FilterOption {
	fo.Lang = lang
	return fo
}

// Len returns length of pre-resultset.
func (fo *FilterOption) Len() int {
	return len(fo.ids)
}

// PageRange returns position [from:to] in pre-resultset in accordance
// to FilterOption parameters. Retrived values can be used to cut
// result slice.
func (fo *FilterOption) PageRange() (from, to int) {

	return fo.sliceRange(len(fo.ids))
}

func (fo *FilterOption) CalcRange(length int) (from, to int) {
	return fo.sliceRange(length)
}

func (fo *FilterOption) sliceRange(full int) (from, to int) {

	// if all rows without pagination.
	if fo.PageSize <= 0 {
		return 0, full
	}

	from = fo.PageNumber * fo.PageSize

	if from > full {
		from = full - from
	}
	if from < 0 {
		from = 0
	}

	to = from + fo.PageSize
	if to > full {
		to = full
	}
	return from, to
}

// SortSlice sort slice with reversing if SortDesc() == true.
func (fo *FilterOption) SortSlice(s interface{}, f func(int, int) bool) {

	if fo.desc {
		sort.Slice(s, func(i, j int) bool {
			return !f(i, j)
		})
		return
	}

	sort.Slice(s, func(i, j int) bool {
		return f(i, j)
	})
}

// IsIgnoredRow returns true if row deleted flag (i.e. DeletedAt.Valid()) does not
// corresponds with Show rule.
func (fo *FilterOption) IsIgnoredRow(rowDeletedFlag bool) bool {
	if rowDeletedFlag && fo.Show == HideDeleted {
		return true
	}

	if !rowDeletedFlag && fo.Show == ShowDeletedOnly {
		return true
	}

	return false
}

// IsSortRequired returns true if URL had specified "sortBy=".
func (fo *FilterOption) IsSortRequired() bool {
	return fo.sortBy != ""
}

// SortByAttr returns name to be used for resultset sorting.
func (fo *FilterOption) SortAttr() string {
	return fo.sortBy
}

// IsSortDesc return true if sortBy value had "minus" like "sortBy=-name"
func (fo *FilterOption) IsSortDesc() bool {
	return fo.desc
}
