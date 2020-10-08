package bot

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestTip_Help(t *testing.T) {
	bot, err := NewSys("./../../data")
	require.NoError(t, err)
	assert.Equal(t, "say! _– набраться мудрости_\n"+
		"tip! _– получить подсказку_\n"+
		"ping _– ответит pong_\n"+
		"пинг _– ответит понг_\n"+
		"кто?, who? _– Разработчики Krista BI_\n"+
		"когда?, when? _– Когда будет версия_\n"+
		"как?, how? _– Демо версия Krista BI_\n"+
		"правила, rules?, правила? _– правила общения в чате_\n",
		bot.Help())
}

func TestTip_Failed(t *testing.T) {
	_, err := NewSys("/tmp/no-such-place")
	require.Error(t, err)
}
