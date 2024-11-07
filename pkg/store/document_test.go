package store

import (
	"context"
	"testing"

	"github.com/oklog/ulid/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.etcd.io/bbolt"
)

type testItem string

var bucketName = "items"

func TestDocumentMany(t *testing.T) {
	testCases := []struct {
		desc              string
		existing          []testItem
		count             int
		expected          []testItem
		expectedCursorNil bool
	}{
		{
			desc:              "no items returns empty slice",
			count:             0,
			existing:          []testItem{},
			expected:          []testItem{},
			expectedCursorNil: true,
		},
		{
			desc:              "items list smaller than count",
			count:             5,
			existing:          []testItem{"one", "two", "three"},
			expected:          []testItem{"three", "two", "one"},
			expectedCursorNil: true,
		},
		{
			desc:              "items list larger than count",
			count:             3,
			existing:          []testItem{"one", "two", "three", "four", "five"},
			expected:          []testItem{"five", "four", "three"},
			expectedCursorNil: false,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.desc, func(t *testing.T) {
			withTestDB(t, bucketName, func(db *bbolt.DB) {
				testDocDB := docDB[testItem]{db}
				for _, e := range tc.existing {
					require.NoError(t, testDocDB.Save(context.Background(), bucketName, ulid.Make().String(), e))
				}
				items, cursor, err := testDocDB.Many(context.Background(), bucketName, nil, tc.count)
				if tc.expectedCursorNil {
					assert.Nil(t, cursor)
				} else {
					assert.NotNil(t, cursor)
				}
				assert.NoError(t, err)
				assert.Equal(t, items, tc.expected)
			})
		})
	}

	t.Run("paging with cursor to last item", func(t *testing.T) {
		withTestDB(t, bucketName, func(db *bbolt.DB) {
			testDocDB := docDB[testItem]{db}
			items := []testItem{"one", "two", "three", "four", "five"}
			itemKeys := make(map[testItem]ulid.ULID, len(items))
			for _, i := range items {
				itemKeys[i] = ulid.Make()
				require.NoError(t, testDocDB.Save(context.Background(), bucketName, itemKeys[i].String(), i))
			}
			cursorKey := itemKeys[testItem("three")].String()
			page, cursor, err := testDocDB.Many(context.Background(), bucketName, &cursorKey, 3)
			assert.NoError(t, err)
			assert.Nil(t, cursor)
			assert.Equal(t, []testItem{"three", "two", "one"}, page)
		})
	})

	t.Run("paging with returned cursor", func(t *testing.T) {
		withTestDB(t, bucketName, func(db *bbolt.DB) {
			testDocDB := docDB[testItem]{db}
			items := []testItem{"one", "two", "three", "four", "five"}
			for _, i := range items {
				require.NoError(t, testDocDB.Save(context.Background(), bucketName, ulid.Make().String(), i))
			}
			firstPage, firstCursor, err := testDocDB.Many(context.Background(), bucketName, nil, 3)
			assert.NoError(t, err)
			assert.NotNil(t, firstCursor)
			assert.Equal(t, []testItem{"five", "four", "three"}, firstPage)
			secondPage, secondCursor, err := testDocDB.Many(context.Background(), bucketName, firstCursor, 3)
			assert.NoError(t, err)
			assert.Nil(t, secondCursor)
			assert.Equal(t, []testItem{"two", "one"}, secondPage)
		})
	})
}
