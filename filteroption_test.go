package filteroption

import (
	"strings"
	"testing"
)

func TestFilterOptions(t *testing.T) {

	ts := []FilterOption{
		{},
		{PageNumber: 1},
		{SortBy: "name"},
		{SortBy: "-name"},
	}

	for i := range ts {
		pageSize := ts[i].PageSize
		desc := strings.HasPrefix(ts[i].SortBy, "-")

		ts[i].ApplyDefaults()
		if pageSize == 0 && ts[i].PageSize != DefaultPageSize {
			t.Errorf("case #1 failed, expected %d, got %d", DefaultPageSize, ts[i].PageSize)
		}

		if ts[i].SortBy != "" {
			if ts[i].IsSortDesc() != desc {
				t.Errorf("case #2, %d failed, expected %t, got %t", i, desc, ts[i].IsSortDesc())
			}
		}
	}
}
