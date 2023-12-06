package index

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestIndex_ToTree(t *testing.T) {
	object := IndexObject{
		Version: 1,
		Entries: map[string]IndexEntry{
			"main.go":          IndexEntry{FullPathName: "main.go"},
			"src/caca.go":      IndexEntry{FullPathName: "src/caca.go"},
			"src/fortine.go":   IndexEntry{FullPathName: "src/fortine.go"},
			"src/casa/algo.go": IndexEntry{FullPathName: "src/casa/algo.go"},
		},
	}

	rootNode := object.ToTree()

	assert.True(t, rootNode.Root)
	assert.Equal(t, len(rootNode.Children), 2)

	childMainGo := rootNode.Children["main.go"]
	assert.False(t, childMainGo.Root)
	assert.Equal(t, childMainGo.Name, "main.go")
	assert.Equal(t, childMainGo.Entry.FullPathName, "main.go")
	assert.Equal(t, len(childMainGo.Children), 0)

	childSrc := rootNode.Children["src"]
	assert.False(t, childSrc.Root)
	assert.Equal(t, childSrc.Name, "src")
	assert.Equal(t, len(childSrc.Children), 3)

	childSrcCaca := childSrc.Children["caca.go"]
	assert.False(t, childSrcCaca.Root)
	assert.Equal(t, childSrcCaca.Name, "caca.go")
	assert.Equal(t, childSrcCaca.Entry.FullPathName, "src/caca.go")
	assert.Equal(t, len(childSrcCaca.Children), 0)

	childSrcFortine := childSrc.Children["fortine.go"]
	assert.False(t, childSrcFortine.Root)
	assert.Equal(t, childSrcFortine.Name, "fortine.go")
	assert.Equal(t, childSrcFortine.Entry.FullPathName, "src/fortine.go")
	assert.Equal(t, len(childSrcFortine.Children), 0)

	childSrcCasa := childSrc.Children["casa"]
	assert.False(t, childSrcCasa.Root)
	assert.Equal(t, childSrcCasa.Name, "casa")
	assert.Equal(t, len(childSrcCasa.Children), 1)

	childSrcCasaAlgo := childSrcCasa.Children["algo.go"]
	assert.False(t, childSrcCasaAlgo.Root)
	assert.Equal(t, childSrcCasaAlgo.Name, "algo.go")
	assert.Equal(t, childSrcCasaAlgo.Entry.FullPathName, "src/casa/algo.go")
	assert.Equal(t, len(childSrcCasaAlgo.Children), 0)
}
