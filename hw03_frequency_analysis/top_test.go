package hw03frequencyanalysis

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = true

var text = `Как видите, он  спускается  по  лестнице  вслед  за  своим
	другом   Кристофером   Робином,   головой   вниз,  пересчитывая
	ступеньки собственным затылком:  бум-бум-бум.  Другого  способа
	сходить  с  лестницы  он  пока  не  знает.  Иногда ему, правда,
		кажется, что можно бы найти какой-то другой способ, если бы  он
	только   мог   на  минутку  перестать  бумкать  и  как  следует
	сосредоточиться. Но увы - сосредоточиться-то ему и некогда.
		Как бы то ни было, вот он уже спустился  и  готов  с  вами
	познакомиться.
	- Винни-Пух. Очень приятно!
		Вас,  вероятно,  удивляет, почему его так странно зовут, а
	если вы знаете английский, то вы удивитесь еще больше.
		Это необыкновенное имя подарил ему Кристофер  Робин.  Надо
	вам  сказать,  что  когда-то Кристофер Робин был знаком с одним
	лебедем на пруду, которого он звал Пухом. Для лебедя  это  было
	очень   подходящее  имя,  потому  что  если  ты  зовешь  лебедя
	громко: "Пу-ух! Пу-ух!"- а он  не  откликается,  то  ты  всегда
	можешь  сделать вид, что ты просто понарошку стрелял; а если ты
	звал его тихо, то все подумают, что ты  просто  подул  себе  на
	нос.  Лебедь  потом  куда-то делся, а имя осталось, и Кристофер
	Робин решил отдать его своему медвежонку, чтобы оно не  пропало
	зря.
		А  Винни - так звали самую лучшую, самую добрую медведицу
	в  зоологическом  саду,  которую  очень-очень  любил  Кристофер
	Робин.  А  она  очень-очень  любила  его. Ее ли назвали Винни в
	честь Пуха, или Пуха назвали в ее честь - теперь уже никто  не
	знает,  даже папа Кристофера Робина. Когда-то он знал, а теперь
	забыл.
		Словом, теперь мишку зовут Винни-Пух, и вы знаете почему.
		Иногда Винни-Пух любит вечерком во что-нибудь поиграть,  а
	иногда,  особенно  когда  папа  дома,  он больше любит тихонько
	посидеть у огня и послушать какую-нибудь интересную сказку.
		В этот вечер...`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("positive test", func(t *testing.T) {
		if taskWithAsteriskIsCompleted {
			expected := []string{
				"а",         // 8
				"он",        // 8
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"в",         // 4
				"его",       // 4
				"если",      // 4
				"кристофер", // 4
				"не",        // 4
			}
			require.Equal(t, expected, Top10(text))
		} else {
			expected := []string{
				"он",        // 8
				"а",         // 6
				"и",         // 6
				"ты",        // 5
				"что",       // 5
				"-",         // 4
				"Кристофер", // 4
				"если",      // 4
				"не",        // 4
				"то",        // 4
			}
			require.Equal(t, expected, Top10(text))
		}
	})
}

func TestAllowedWords(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: "Нога", expected: []string{"нога"}},
		{input: "Нога, нога! нога 'нога'", expected: []string{"нога"}},
		{input: "Нога - нога", expected: []string{"нога"}},
		{input: "-", expected: []string{}},
		{input: "Какой-то какойто", expected: []string{"какой-то", "какойто"}},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			require.Equal(t, tc.expected, Top10(tc.input))
		})
	}
}

func TestEnglishWords(t *testing.T) {
	tests := []struct {
		input    string
		expected []string
	}{
		{input: "Leg", expected: []string{"leg"}},
		{input: "Leg, Leg! leg 'leg'", expected: []string{"leg"}},
		{input: "Leg - leg", expected: []string{"leg"}},
		{
			input:    "One, two, two, three, three, three, four, four, four, four",
			expected: []string{"four", "three", "two", "one"},
		},
	}
	for _, tc := range tests {
		tc := tc
		t.Run(tc.input, func(t *testing.T) {
			require.Equal(t, tc.expected, Top10(tc.input))
		})
	}
}

func TestWordsOrderWithTheSameFrequency(t *testing.T) {
	var sb strings.Builder

	postfixes := []string{"15", "14", "13", "12", "11", "10", "09", "08", "07", "06", "05", "04", "03", "02", "01"}

	for i := 0; i < 15; i++ {
		sb.WriteString(strings.Repeat(fmt.Sprintf("слово-%s ", postfixes[i]), 133))
	}
	for i := 0; i < 10; i++ {
		sb.WriteString(strings.Repeat(fmt.Sprintf("другое-слово-%s ", postfixes[i]), 99))
	}

	text := sb.String()

	expected := []string{
		"слово-01",
		"слово-02",
		"слово-03",
		"слово-04",
		"слово-05",
		"слово-06",
		"слово-07",
		"слово-08",
		"слово-09",
		"слово-10",
	}
	require.Equal(t, expected, Top10(text))
}
