package routing

import (
	"testing"

	"gotest.tools/assert"
)

func TestAdjustedRefreshrate(t *testing.T) {

	assert.Equal(t, GetAdjustedRefreshrate(15, 16, 8, 2), 8)
	assert.Equal(t, GetAdjustedRefreshrate(14, 16, 8, 2), 7)
	assert.Equal(t, GetAdjustedRefreshrate(13, 16, 8, 2), 6)
	assert.Equal(t, GetAdjustedRefreshrate(12, 16, 8, 2), 5)
	assert.Equal(t, GetAdjustedRefreshrate(11, 16, 8, 2), 4)
	assert.Equal(t, GetAdjustedRefreshrate(10, 16, 8, 2), 4)
	assert.Equal(t, GetAdjustedRefreshrate(9, 16, 8, 2), 3)
	assert.Equal(t, GetAdjustedRefreshrate(8, 16, 8, 2), 2)
	assert.Equal(t, GetAdjustedRefreshrate(7, 16, 8, 2), 2)
	assert.Equal(t, GetAdjustedRefreshrate(6, 16, 8, 2), 2)
	assert.Equal(t, GetAdjustedRefreshrate(5, 16, 8, 2), 1)
	assert.Equal(t, GetAdjustedRefreshrate(4, 16, 8, 2), 1)
	assert.Equal(t, GetAdjustedRefreshrate(3, 16, 8, 2), 1)

	assert.Equal(t, GetAdjustedRefreshrate(15, 16, 8, 3), 7)
	assert.Equal(t, GetAdjustedRefreshrate(14, 16, 8, 3), 6)
	assert.Equal(t, GetAdjustedRefreshrate(13, 16, 8, 3), 5)
	assert.Equal(t, GetAdjustedRefreshrate(12, 16, 8, 3), 4)
	assert.Equal(t, GetAdjustedRefreshrate(11, 16, 8, 3), 3)
	assert.Equal(t, GetAdjustedRefreshrate(10, 16, 8, 3), 2)
	assert.Equal(t, GetAdjustedRefreshrate(9, 16, 8, 3), 2)
	assert.Equal(t, GetAdjustedRefreshrate(8, 16, 8, 3), 1)
	assert.Equal(t, GetAdjustedRefreshrate(7, 16, 8, 3), 1)
	assert.Equal(t, GetAdjustedRefreshrate(3, 16, 8, 3), 1)
	assert.Equal(t, GetAdjustedRefreshrate(15, 16, 2, 3), 2)
	assert.Equal(t, GetAdjustedRefreshrate(14, 16, 2, 3), 2)
	assert.Equal(t, GetAdjustedRefreshrate(13, 16, 2, 3), 2)
}
