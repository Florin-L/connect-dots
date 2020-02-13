package game

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

var blobJson = []byte(`
{
"size": 5,
"difficulty": 0,
"dots": [
    {
        "x": 0,
        "y": 0,
        "color": "red"
    },
    {
        "x": 2,
        "y": 0,
        "color": "yellow"
    },
    {
        "x": 3,
        "y": 0,
        "color": "green"
    },
    {
        "x": 4,
        "y": 0,
		"color": "orange"
    },
    {
        "x": 4,
        "y": 2,
        "color": "blue"
    },
    {
        "x": 1,
        "y": 3,
        "color": "yellow"
    },
    {
        "x": 2,
        "y": 3,
        "color": "green"
    },
    {
        "x": 4,
        "y": 3,
        "color": "orange"
    },
    {
        "x": 2,
        "y": 4,
        "color": "red"
    },
    {
        "x": 4,
        "y": 4,
        "color": "blue"
    }
]
}
`)

func TestLoad(t *testing.T) {
	l, err := Load(blobJson)
	assert.Nil(t, err)
	assert.Equal(t, l.Size, int32(5))
}

func TestLoadWrongSizeValue(t *testing.T) {
	var json = []byte(`
	{
	"size": 0,
	"difficulty": 1,
	"dots": [
		{
		"x": 1,
		"y": 2,
		"color": "red"
		}
	]
	}
`)

	l, err := Load(json)
	assert.Nil(t, l)
	assert.NotNil(t, err)
}

func TestLoadNoDots(t *testing.T) {
	var json = []byte(`
	{
	"size": 5,
	"difficulty": 1,
	"dots": []
	}
`)

	l, err := Load(json)
	assert.Nil(t, l)
	assert.NotNil(t, err)
}

func TestLoadWrongDifficultyValue(t *testing.T) {
	var json = []byte(`
	{
	"size": 5,
	"difficulty": 4,
	"dots": [
		{
		"x": 1,
		"y": 2,
		"color": "red"
		}
	]
	}
`)

	l, err := Load(json)
	assert.Nil(t, l)
	assert.NotNil(t, err)
}
