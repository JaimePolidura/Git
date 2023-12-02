package ignore

import (
	"bytes"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestGitIgnore_Deserialize(t *testing.T) {
	gitIgnoreBytes := []byte("#Mi comentario\n" +
		".idea\n" +
		"node_modules\n\n" +
		"#Mi otro comentario\n" +
		" \n" +
		"   .ignorar  \n")

	gitIgnore, err := Deserialize(bytes.NewReader(gitIgnoreBytes))

	assert.Nil(t, err)
	assert.Equal(t, len(gitIgnore.ignoredRules), 3)

	gitIgnore.IsIgnored(".idea")

	ignored, err := gitIgnore.IsIgnored(".idea")
	assert.Nil(t, err)
	assert.True(t, ignored)
	
	ignored, err = gitIgnore.IsIgnored("node_modules")
	assert.Nil(t, err)
	assert.True(t, ignored)

	ignored, err = gitIgnore.IsIgnored(".ignorar")
	assert.Nil(t, err)
	assert.True(t, ignored)

	ignored, err = gitIgnore.IsIgnored(".a")
	assert.Nil(t, err)
	assert.False(t, ignored)
}
