package discovery

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestDecodeTxtRecord(t *testing.T) {

	source := `id=87cf98a003f1f1dbd2efe6d19055a617|ve=04|md=Chromecast|ic=/setup/icon.png|fn=Chromecast PO|ca=5|st=0|bs=FA8FCA7EE8A9|rs=`

	result := decodeTxtRecord(source)

	assert.Equal(t, result["id"], "87cf98a003f1f1dbd2efe6d19055a617")

}
