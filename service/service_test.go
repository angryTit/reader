package service

import (
	"bytes"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/angryTit/reader/types"
)

var (
	contentStr = []string{
		"1, 127.0.0.1, 17:51:59\n",
		"2, 127.0.0.1, 17:52:59\n",
		"1, 127.0.0.2, 17:53:59\n",
		"2, 127.0.0.2, 17:54:59\n",
		"2, 127.0.0.3, 17:55:59\n",
		"3, 127.0.0.3, 17:55:59\n",
		"3, 127.0.0.1, 17:56:59\n",
		"4, 127.0.0.1, 17:57:59\n",
	}
)

func TestIsSame(t *testing.T) {
	content := bytes.NewReader([]byte(strings.Join(contentStr, "")))

	t.Run("successfully execute", func(t *testing.T) {
		storage := types.NewStorage()
		_, err := FillStorage(content, 0, storage)
		require.NoError(t, err)

		require.True(t, IsSame("1", "2", storage))
		require.True(t, IsSame("2", "1", storage))
		require.True(t, IsSame("2", "3", storage))
		require.True(t, IsSame("3", "2", storage))
		require.True(t, IsSame("1", "1", storage))

		require.False(t, IsSame("1", "4", storage))
		require.False(t, IsSame("1", "3", storage))
		require.False(t, IsSame("3", "1", storage))
	})
}

func TestFillStorage(t *testing.T) {
	duplicates := make([]string, 0)
	copy(duplicates, contentStr)
	duplicates = append(duplicates, contentStr...)
	content := bytes.NewReader([]byte(strings.Join(duplicates, "")))

	t.Run("successfully fill storage", func(t *testing.T) {
		storage := types.NewStorage()
		currentPosition, err := FillStorage(content, 0, storage)
		require.NoError(t, err)
		require.Greater(t, *currentPosition, int64(0))

		target := *storage.Get("1").GetSlice()
		require.Equal(t, 2, len(target))
		require.True(t, contains(target, "127.0.0.1"))
		require.True(t, contains(target, "127.0.0.2"))

		target = *storage.Get("2").GetSlice()
		require.Equal(t, 3, len(target))
		require.True(t, contains(target, "127.0.0.1"))
		require.True(t, contains(target, "127.0.0.2"))
		require.True(t, contains(target, "127.0.0.3"))

		target = *storage.Get("3").GetSlice()
		require.Equal(t, 2, len(target))
		require.True(t, contains(target, "127.0.0.1"))
		require.True(t, contains(target, "127.0.0.3"))

		target = *storage.Get("4").GetSlice()
		require.Equal(t, 1, len(target))
		require.True(t, contains(target, "127.0.0.1"))
	})
}

func TestParse(t *testing.T) {
	t.Run("successfully parse", func(t *testing.T) {
		target := parse(contentStr)
		require.Equal(t, 2, len(target["1"]))
		require.True(t, contains(target["1"], "127.0.0.1"))
		require.True(t, contains(target["1"], "127.0.0.2"))

		require.Equal(t, 3, len(target["2"]))
		require.True(t, contains(target["2"], "127.0.0.1"))
		require.True(t, contains(target["2"], "127.0.0.2"))
		require.True(t, contains(target["2"], "127.0.0.3"))

		require.Equal(t, 2, len(target["3"]))
		require.True(t, contains(target["3"], "127.0.0.1"))
		require.True(t, contains(target["3"], "127.0.0.3"))

		require.Equal(t, 1, len(target["4"]))
		require.True(t, contains(target["4"], "127.0.0.1"))
	})
}

func contains(source []string, target string) bool {
	for _, each := range source {
		if each == target {
			return true
		}
	}
	return false
}

func TestReadFrom(t *testing.T) {
	content := bytes.NewReader([]byte(strings.Join(contentStr, "")))

	t.Run("successfully read whole content", func(t *testing.T) {
		result, position, err := readFrom(content, 0)

		require.NoError(t, err)
		require.Greater(t, *position, int64(0))
		require.Equal(t, contentStr, result)
		require.Equal(t, true, true)
	})

	t.Run("successfully read from position", func(t *testing.T) {
		result, position, err := readFrom(content, 23)
		require.NoError(t, err)
		require.Greater(t, *position, int64(23))
		require.Equal(t, contentStr[1:], result)

		result, position, err = readFrom(content, 46)
		require.NoError(t, err)
		require.Greater(t, *position, int64(23))
		require.Equal(t, contentStr[2:], result)

		result, position, err = readFrom(content, 69)
		require.NoError(t, err)
		require.Greater(t, *position, int64(23))
		require.Equal(t, contentStr[3:], result)
	})

}
