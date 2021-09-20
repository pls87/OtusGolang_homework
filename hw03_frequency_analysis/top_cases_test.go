package hw03frequencyanalysis_test

var tests = []struct {
	input    string
	name     string
	expected []string
}{
	{
		name:     "Empty string",
		input:    "",
		expected: []string{},
	}, {
		name: "Винни-Пух", input: `Как видите, он  спускается  по  лестнице  вслед  за  своим
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
		В этот вечер...`,
		expected: []string{
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
		},
	}, {
		name: "Собака...",
		input: `Соб@к@ бывает кусачей
	Только от жизни собачей.
	Только от жизни, от жизни собачей,
	Соб@к@ бывает кусачей.

	Соб@к@ хватает зубами за пятку,
	Соб@к@ съедает гражданку лошадку,
	И с ней гражданина кота,
	Когда проживает собака не в будке,
	Когда у нее завывает в желудке.
	И каждому ясно, что эта собака
К	руглая сирота.`,
		expected: []string{
			"соб@к@",  // 4
			"жизни",   // 3
			"от",      // 3
			"бывает",  // 2
			"в",       // 2
			"и",       // 2
			"когда",   // 2
			"кусачей", // 2
			"собака",  // 2
			"собачей", // 2
		},
	}, {
		name: "Плов",
		input: `Ингридиенты:
    * 1 кг длиннозёрного риса
    * 1 кг баранины
    * 1 кг моркови
    * 300 мл растительного масла
    * 4 небольшие луковицы
    * 2 небольших сухих острых перчика
    * чеснок
    * 1 ст. л. сушеного барбариса
    * 1 ст. л. зиры
    * 1 ч. л. семян кориандра
    * соль
	Приготовление:
	* Шаг 1: Рис для узбекского плова в казане промыть холодной водой, меняя ее несколько раз. 
	Последняя вода после промывки должна остаться совершенно прозрачной.
	* Шаг 2: Баранину для плова вымыть и нарезать кубиками. 3 луковицы и всю морковь очистить. 
	Лук порезать тонкими полукольцами, морковь – длинными брусками толщиной 1 см.
	* шаг 3: Казан или толстостенную кастрюлю разогреть, влить масло и прокалить его до появления светлого дымка. 
	Добавить оставшуюся луковицу и хорошо обжарить ее до темно-золотистого цвета. Вытащить луковицу из кастрюли.
	* Шаг 4: Подготовить зирвак (основу узбекского плова). 
	Положить в казан нарезанный лук и, помешивая, обжарить до темно-золотистого цвета в течение 5-7 мин. 
	Следите, чтобы он не пригорел. Добавить к луку нарезанную баранину. 
	Помешивая кулинарной лопаткой, жарить ингредиенты зирвака до коричневой румяной корочки. 
	Это может занять около 10–15 минут.
	Выложить в казан к мясу с луком морковь. Жарить, не перемешивая, 3 минуты. 
	Затем все содержимое казана перемешать и готовить 10 минут, слегка помешивая лопаткой.
	Продолжение следует...
`,
		expected: []string{
			"1",       // 8
			"и",       // 6
			"в",       // 4
			"до",      // 4
			"шаг",     // 4
			"3",       // 3
			"казан",   // 3
			"кг",      // 3
			"л",       // 3
			"морковь", // 2
		},
	}, {
		name: "Garry Potter",
		input: `THE BOY WHO LIVED
	Mr. and Mrs. Dursley,	 of number four, Privet Drive, 
	were proud to say that they were perfectly normal, 
	thank you very much. They were the last people you’d 
	expect to be involved in anything strange or 
	mysterious, because they just didn’t hold with such 
	nonsense. 

	Mr. Dursley was the director of a firm called 
	Grunnings, which made drills. He was a big, beefy 
	man with hardly any neck, although he did have a 
	very large mustache. Mrs. Dursley was thin and 
	blonde and had nearly twice the usual amount of 
	neck, which came in 	very useful as she spent so 
	much of her time craning over garden fences, spying 
	on the neighbors. The Dursley s had a small son 
	called Dudley and in their opinion there was no finer 
	boy anywhere. 

	The Dursleys had everything they wanted, but they 
	also had a secret, and their greatest fear was that 
	somebody would discover it. They didn’t think they 
	could bear it if anyone found out about the Potters. 
	Mrs. Potter was Mrs. Dursley’s sister, but they hadn’t
	met for several years;	 in fact, Mrs. Dursley pretended 
	she didn’t have a sister, because her sister and her 
	good-for-nothing husband were as unDursleyish as it 
	was possible to be. The Dursleys shuddered to think 
	what the neighbors would say if the Potters arrived in 
	the street. The Dursleys knew that the Potters had a 
	small son, too, but they had never even seen him. 

	This boy was another good reason for keeping the 
	Potters away; they didn’t want Dudley mixing with a 
	child like that. 
	
	When Mr. and Mrs. Dursley woke up on the dull, gray 
	Tuesday our story starts, there was nothing about the 
	cloudy sky outside to suggest that strange and 
	mysterious things would soon be happening all over 
	the country. Mr. Dursley hummed as he picked out 
	his most boring tie for work, and Mrs. Dursley 
	gossiped away happily as she wrestled a screaming 
	Dudley into his high chair. 
	
	None of them noticed a large, tawny owl flutter past 
	the window. `,
		expected: []string{
			"the",     // 18
			"a",       // 10
			"they",    // 10
			"and",     // 9
			"was",     // 8
			"dursley", // 8
			"mrs",     // 7
			"had",     // 6
			"as",      // 5
			"in",      // 5
		},
	}, {
		name: "One syntetic case with less than 10 words",
		input: `	 (asd)! (rtsfg?! rtsfg asd. a a a!?, c-c-c, ccc, ccc 	`,
		expected: []string{
			"a",     // 3
			"asd",   // 2
			"ccc",   // 2
			"rtsfg", // 2
			"c-c-c", // 1
		},
	},
}
