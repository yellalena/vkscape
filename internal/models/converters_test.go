package models

import (
	"reflect"
	"testing"
)

func TestFilterAlbumsByIDs(t *testing.T) {
	albums := []PhotoAlbum{
		{ID: 1, Title: "A"},
		{ID: 2, Title: "B"},
		{ID: 3, Title: "C"},
	}

	tests := []struct {
		name        string
		inputIDs    []string
		wantIDs     []int
		wantInvalid []string
	}{
		{name: "all valid", inputIDs: []string{"1", "3"}, wantIDs: []int{1, 3}, wantInvalid: nil},
		{name: "invalid mixed", inputIDs: []string{"2", "x", "3", "bad"}, wantIDs: []int{2, 3}, wantInvalid: []string{"x", "bad"}},
		{name: "none", inputIDs: []string{}, wantIDs: nil, wantInvalid: nil},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			got, invalid := FilterAlbumsByIDs(tc.inputIDs, albums)
			if len(invalid) != len(tc.wantInvalid) {
				t.Fatalf("expected invalid %v, got %v", tc.wantInvalid, invalid)
			}
			if !reflect.DeepEqual(invalid, tc.wantInvalid) {
				t.Fatalf("expected invalid %v, got %v", tc.wantInvalid, invalid)
			}

			if len(got) != len(tc.wantIDs) {
				t.Fatalf("expected %d albums, got %d", len(tc.wantIDs), len(got))
			}
			if !reflect.DeepEqual(extractIDs(got), tc.wantIDs) {
				t.Fatalf("expected IDs %v, got %v", tc.wantIDs, extractIDs(got))
			}
		})
	}
}

func extractIDs(albums []PhotoAlbum) []int {
	out := make([]int, 0, len(albums))
	for _, a := range albums {
		out = append(out, a.ID)
	}
	return out
}
