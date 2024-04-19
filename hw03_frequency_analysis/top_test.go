package hw03frequencyanalysis

import (
	"testing"

	"github.com/stretchr/testify/require"
)

// Change to true if needed.
var taskWithAsteriskIsCompleted = false

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

var text2 = `Lorem ipsum dolor sit amet, consectetur adipiscing elit.
Aenean vulputate ipsum nec nunc consequat tincidunt. Donec eget volutpat
lorem, id blandit lorem. Pellentesque dapibus mattis metus, at semper 
diam mattis in. Sed tincidunt justo nisi, eget faucibus dolor tristique
at. Vivamus in viverra risus. Mauris ut turpis sagittis massa pretium 
efficitur eget non urna. Duis tristique, erat id accumsan blandit, purus 
diam tristique leo, eget accumsan eros tortor eget diam. Etiam ultrices 
placerat risus, ut fermentum est lobortis eget. Aliquam semper, enim sit 
amet dictum rutrum, lacus tellus pharetra lacus, id vehicula diam lacus 
vitae turpis. Nulla laoreet ultricies vulputate. Pellentesque posuere 
nulla sed cursus iaculis. Pellentesque tempus porttitor massa a fermentum.
Integer convallis consectetur convallis. Maecenas egestas ex tincidunt 
erat vehicula rutrum. Morbi ipsum libero, consequat sed lectus sed, tincidunt
bibendum augue.
Duis consequat feugiat orci lobortis eleifend. Curabitur sollicitudin rutrum 
gravida. Nam efficitur posuere justo sed molestie. Nunc euismod, dui a congue
posuere, mauris dui egestas mi, eu iaculis risus libero et mi. Aliquam ut 
convallis justo. Vestibulum vel eros sed nisi porttitor malesuada. Praesent 
sollicitudin quis nunc at auctor. Cras lectus magna, vehicula a pulvinar vel,
egestas et metus. Morbi congue tristique leo nec imperdiet. Mauris nisi magna,
malesuada eget tellus sed, dignissim aliquam justo. Nulla leo nisl, porta 
in neque a, laoreet sollicitudin nisl. Maecenas nec cursus nisl, sed 
consequat mauris. Etiam neque eros, tincidunt vehicula auctor quis, tincidunt
ut neque. Nam lobortis eros nec velit ultricies, vel dignissim magna faucibus.
Nulla sodales quis ex vel dictum.
Vestibulum ante ipsum primis in faucibus orci luctus et ultrices posuere cubilia curae;
Praesent ac mauris massa. Nam nibh ligula, pulvinar et elit sit amet, interdum efficitur
felis. Fusce non diam in turpis efficitur sagittis at sit amet justo. Etiam turpis massa,
lobortis vel lacus in, vestibulum semper velit. Integer commodo dui elit, sed tincidunt 
diam consectetur quis. In eu dui in lectus mollis tempor sit amet in felis. Interdum et
malesuada fames ac ante ipsum primis in faucibus.
Sed quis pharetra magna. Nam efficitur rutrum augue, eget tincidunt elit venenatis in. Donec 
vel tincidunt ante, ac dignissim lorem. Sed vitae dui convallis, lobortis tellus vitae, 
posuere massa. Pellentesque suscipit enim justo, et sodales orci eleifend non. Curabitur nec 
libero sollicitudin, vulputate diam eget, feugiat felis. Donec eu felis laoreet, tincidunt arcu
eu, pulvinar orci. Maecenas erat metus, dignissim in ultricies non, aliquam id neque. Proin 
quis elementum sapien. Praesent scelerisque cursus sem, nec ullamcorper magna eleifend in.
Donec auctor pulvinar tempus. Sed turpis lacus, mollis sed odio ut, faucibus semper lectus.
Quisque vel tempor odio, eget tristique eros. Donec venenatis ut metus vel tincidunt. Curabitur
molestie tristique finibus.
Vestibulum magna leo, mollis quis tortor nec, rhoncus hendrerit dui. Orci varius natoque penatibus
et magnis dis parturient montes, nascetur ridiculous mus. Nullam vel interdum nulla, et posuere orci.
Curabitur accumsan elit posuere dui facilisis, eu tristique turpis aliquam. Suspendisse sed mi urna.
Nam aliquam sapien mi, vitae laoreet felis facilisis ut. Vestibulum accumsan ipsum non urna elementum
vestibulum. In metus est, accumsan nec leo id, commodo laoreet libero. Cras sodales, tellus nec 
ultrices vulputate, erat felis congue lorem, eget suscipit magna nisl et tellus. Nunc egestas turpis
sit amet porta feugiat. Etiam lobortis diam ac metus laoreet mollis. Mauris sit amet faucibus sapien,
at facilisis enim. Pellentesque erat dui, luctus sed turpis eget, finibus condimentum est.`

func TestTop10(t *testing.T) {
	t.Run("no words in empty string", func(t *testing.T) {
		require.Len(t, Top10(""), 0)
	})

	t.Run("case sensitivity", func(t *testing.T) {
		require.Len(t, Top10("Нога нога"), 2)
	})

	t.Run("punctuation sensitivity", func(t *testing.T) {
		require.Len(t, Top10("нога, нога - ---"), 4)
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

			expected = []string{
				"eget",      // 9
				"et",        // 9
				"sed",       // 9
				"tincidunt", // 9
				"in",        // 8
				"nec",       // 8
				"vel",       // 8
				"diam",      // 7
				"sit",       // 7
				"turpis",    // 7
			}
			require.Equal(t, expected, Top10(text2))
		}
	})
}
