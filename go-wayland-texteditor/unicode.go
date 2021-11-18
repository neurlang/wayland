package main

var extendedaDescriptor = "" +
	`Ā	Ġ	Ŀ	ş	ā	ġ	ŀ	Š
	Ă	Ģ	Ł	š	ă	ģ	ł	Ţ
	Ą	Ĥ	Ń	ţ	ą	ĥ	ń	Ť
	Ć	Ħ	Ņ	ť	ć	ħ	ņ	Ŧ
	Ĉ	Ĩ	Ň	ŧ	ĉ	ĩ	ň	Ũ
	Ċ	Ī	ŉ	ũ	ċ	ī	Ŋ	Ū
	Č	Ĭ	ŋ	ū	č	ĭ	Ō	Ŭ
	Ď	Į	ō	ŭ	ď	į	Ŏ	Ů
	Đ	İ	ŏ	ů	đ	ı	Ő	Ű
	Ē	Ĳ	ő	ű	ē	ĳ	Œ	Ų
	Ĕ	Ĵ	œ	ų	ĕ	ĵ	Ŕ	Ŵ
	Ė	Ķ	ŕ	ŵ	ė	ķ	Ŗ	Ŷ
	Ę	ſ	ŗ	ŷ	ę	ĸ	Ř	Ÿ
	Ě	Ĺ	ř	Ź	ě	ĺ	Ś	ź
	Ĝ	Ļ	ś	Ż	ĝ	ļ	Ŝ	ż
	Ğ	Ľ	ŝ	Ž	ğ	ľ	Ş	ž`

var supplementDescriptor = "" +
	`¡	±	À	Ï	à	ï
	¢	²	Á	Ñ	á	ñ
	£	³	Â	Ò	â	ò
	¤	´	Ã	Ó	ã	ó
	¥	µ	Ä	Ô	ä	ô
	¦	¶	Å	Õ	å	õ
	§	·	Æ	Ö	æ	ö
	¨	¸	Ç	Ø	ç	ø
	©	¹	Ð	Ù	ð	ù
	ª	º	È	Ú	è	ú
	«	»	É	Û	é	û
	¬	¼	Ê	Ü	ê	ü
	½	¾	Ë	Ý	ë	ý
	­®	¿	Ì	Ÿ	ì	ÿ
	¯	÷	Í	Þ	í	þ
	°	×	Î	ẞ	î	ß`

var extendedbDescriptor = "" +
	`ǎ	Ǎ	ǳ	ǲ	ȟ	Ȟ	ǒ	Ǒ	ȑ	Ȑ	ȳ	Ȳ
	ǻ	Ǻ	Ǳ	ǆ	ƕ	Ƕ	ȫ	Ȫ	ȓ	Ȓ	ɏ	Ɏ
	ǟ	Ǟ	ǅ	Ǆ	ǐ	Ǐ	ȭ	Ȭ	ɍ	Ɍ	ƴ	Ƴ
	ȧ	Ȧ	Ɖ	Ɗ	ȉ	Ȉ	ȯ	Ȯ	ș	Ș	ȝ	Ȝ
	ǡ	Ǡ	ȩ	Ȩ	ȋ	Ȋ	ȱ	Ȱ	Ʀ	ȿ	ƶ	Ƶ
	ȁ	Ȁ	ȅ	Ȅ	Ɨ	Ɩ	ǿ	Ǿ	ț	Ț	ȥ	Ȥ
	ȃ	Ȃ	ȇ	Ȇ	ǰ	ȷ	ǫ	Ǫ	ƾ	Ⱦ	ǯ	Ǯ
	ǽ	Ǽ	ɇ	Ɇ	ɉ	Ɉ	ǭ	Ǭ	Ƭ	Ʈ	ƹ	Ƹ
	ǣ	Ǣ	ǝ	Ǝ	ǩ	Ǩ	ȍ	Ȍ	ǔ	Ǔ	ƿ	Ƿ
	Ⱥ	ƀ	Ə	Ɛ	ƙ	Ƙ	ȏ	Ȏ	ǘ	Ǘ	ǜ	Ǜ
	Ƀ	Ɓ	ƒ	Ƒ	ǉ	ǈ	Ɲ	ȵ	ǚ	Ǚ	ƽ	Ƽ
	ƃ	Ƃ	ǵ	Ǵ	ƚ	Ƚ	Ǉ	Ɔ	ǖ	Ǖ	ƅ	Ƅ
	ȼ	Ȼ	ǧ	Ǧ	ȴ	ƛ	ơ	Ơ	ȕ	Ȕ	ƨ	Ƨ
	ƈ	Ƈ	ǥ	Ǥ	ǹ	Ǹ	ȣ	Ȣ	ȗ	Ȗ	ȡ	ƫ
	ȸ	ȹ	Ɠ	Ɣ	ǌ	ǋ	ƥ	Ƥ	ɋ	Ɋ	ư	Ư
	ƌ	Ƌ	ƣ	Ƣ	ƞ	Ƞ	Ǌ	Ɵ	Ʃ	ƪ	ƭ	ȶ`

var ipaDescriptor = "" +
	`ɐ	ɑ	ɒ	ɓ	ɔ	ɕ	ɖ	ɗ	ɘ	ə	ɚ	ɛ
	ɜ	ɝ	ɞ	ɟ	ɠ	ɡ	ɢ	ɣ	ɤ	ɥ	ɦ	ɧ
	ɨ	ɩ	ɪ	ɫ	ɬ	ɭ	ɮ	ɯ	ɰ	ɱ	ɲ	ɳ
	ɴ	ɵ	ɶ	ɷ	ɸ	ɹ	ɺ	ɻ	ɼ	ɽ	ɾ	ɿ
	ʀ	ʁ	ʂ	ʃ	ʄ	ʅ	ʆ	ʇ	ʈ	ʉ	ʊ	ʋ
	ʌ	ʍ	ʎ	ʏ	ʐ	ʑ	ʒ	ʓ	ʔ	ʕ	ʖ	ʗ
	ʘ	ʙ	ʚ	ʛ	ʜ	ʝ	ʞ	ʟ	ʠ	ʡ	ʢ	ʣ
	ʤ	ʥ	ʦ	ʧ	ʨ	ʩ	ʪ	ʫ	ʬ	ʭ	ʮ	ʯ`

var greekDescriptor = "" +
	`Α	α	Ά	ά	Β	β	ϐ
	Γ	γ	Δ	δ	Ε	ε	ϵ
	Έ	έ	Ϝ	ϝ	Ͷ	ͷ	϶
	Ϛ	ϛ	Ζ	ζ	Ͱ	ͱ	ϑ
	Η	η	Ή	ή	Θ	θ	ϴ
	Ι	ι	Ί	ί	Ϊ	ϊ	ΐ
	Ϳ	ϳ	Κ	κ	ϰ	Ϗ	ϗ
	Λ	λ	Μ	μ	Ν	ν	
	Ξ	ξ	Ο	ο	Ό	ό	
	Π	π	ϖ	Ϻ	ϻ	Ϟ	ϟ
	Ϙ	ϙ	Ρ	ρ	ϱ	ϼ	
	Σ	ς	σ	Ϲ	ϲ	Ͼ	ͼ
	Ͻ	ͻ	Ͽ	ͽ	Τ	τ	
	Υ	υ	ϒ	Ύ	ύ	ϓ	
	Ϋ	ϋ	ϔ	ΰ	Φ	φ	ϕ
	Χ	χ	Ψ	ψ	Ω	ω	
	Ώ	ώ	Ϡ	ϡ	Ͳ	ͳ	
	Ϸ	ϸ	Ϣ	ϣ	Ϥ	ϥ	
	Ϧ	ϧ	Ϩ	ϩ	Ϫ	ϫ	
	Ϭ	ϭ	Ϯ	ϯ			`

var cyrillicDescriptor = "" +
	`А	а	Ӑ	ӑ	Ӓ	ӓ	Ә	ә	Ӛ	ӛ
	Ӕ	ӕ	Б	б	В	в	Г	г	Ѓ	ѓ
	Ґ	ґ	Ғ	ғ	Ӻ	ӻ	Ҕ	ҕ	Ӷ	ӷ
	Д	д	Ђ	ђ	Ҙ	ҙ	Е	е	Ѐ	ѐ
	Ӗ	ӗ	Ё	ё	Є	є	Ж	ж	Ӂ	ӂ
	Ӝ	ӝ	Җ	җ	З	з	Ӟ	ӟ	Ѕ	ѕ
	Ӡ	ӡ	И	и	Ѝ	ѝ	Ӥ	ӥ	Ӣ	ӣ
	Ҋ	ҋ	І	і	Ї	ї	Й	й	Ј	ј
	К	к	Ќ	ќ	Қ	қ	Ӄ	ӄ	Ҡ	ҡ
	Ҟ	ҟ	Ҝ	ҝ	Л	л	Ӆ	ӆ	Љ	љ
	М	м	Ӎ	ӎ	Н	н	Ӊ	ӊ	Ң	ң
	Ӈ	ӈ	Ҥ	ҥ	Њ	њ	О	о	Ӧ	ӧ
	Ө	ө	Ӫ	ӫ	П	п	Ҧ	ҧ	Ҁ	ҁ
	Р	р	Ҏ	ҏ	С	с	Ҫ	ҫ	Т	т
	Ҭ	ҭ	Ћ	ћ	У	у	Ў	ў	Ӱ	ӱ
	Ӳ	ӳ	Ӯ	ӯ	Ү	ү	Ұ	ұ	Ѹ	ѹ
	Ф	ф	Х	х	Ӽ	ӽ	Ӿ	ӿ	Ҳ	ҳ
	Һ	һ	Ѡ	ѡ	Ѿ	ѿ	Ѽ	ѽ	Ѻ	ѻ
	Ц	ц	Ҵ	ҵ	Ч	ч	Ӵ	ӵ	Ҷ	ҷ
	Ӌ	ӌ	Ҹ	ҹ	Ҽ	ҽ	Ҿ	ҿ	Џ	џ
	Ш	ш	Щ	щ	Ъ	ъ	Ы	ы	Ӹ	ӹ
	Ь	ь	Ҍ	ҍ	Ѣ	ѣ	Э	э	Ӭ	ӭ
	Ю	ю	Я	я	Ѥ	ѥ	Ѧ	ѧ	Ѫ	ѫ
	Ѩ	ѩ	Ѭ	ѭ	Ѯ	ѯ	Ѱ	ѱ	Ѳ	ѳ
	Ѵ	ѵ	Ѷ	ѷ	Ҩ	ҩ	Ӏ	ӏ	҂	`

var vietnameseDescriptor = "" +
	`ả	Ả	ạ	Ạ	ằ	Ằ	ẳ	Ẳ	ẵ	Ẵ
	ắ	Ắ	ặ	Ặ	ầ	Ầ	ẩ	Ẩ	ẫ	Ẫ
	ấ	Ấ	ậ	Ậ	ẻ	Ẻ	ẽ	Ẽ	ẹ	Ẹ
	ề	Ề	ể	Ể	ễ	Ễ	ế	Ế	ệ	Ệ
	ỉ	Ỉ	ị	Ị	ỏ	Ỏ	ọ	Ọ	ồ	Ồ
	ổ	Ổ	ỗ	Ỗ	ố	Ố	ộ	Ộ	ờ	Ờ
	ở	Ở	ỡ	Ỡ	ớ	Ớ	ợ	Ợ	ủ	Ủ
	ụ	Ụ	ừ	Ừ	ử	Ử	ữ	Ữ	ứ	Ứ
	ự	Ự	ỳ	Ỳ	ỷ	Ỷ	ỹ	Ỹ	ỵ	Ỵ`

var hangul0Descriptor = "" +
	`가	개	갸	걔	거	게	겨	계	고	과	괘	괴	교	구	궈	궤	귀	규	그	긔	기
	까	깨	꺄	꺠	꺼	께	껴	꼐	꼬	꽈	꽤	꾀	꾜	꾸	꿔	꿰	뀌	뀨	끄	끠	끼
	나	내	냐	냬	너	네	녀	녜	노	놔	놰	뇌	뇨	누	눠	눼	뉘	뉴	느	늬	니
	다	대	댜	댸	더	데	뎌	뎨	도	돠	돼	되	됴	두	둬	뒈	뒤	듀	드	듸	디
	따	때	땨	떄	떠	떼	뗘	뗴	또	똬	뙈	뙤	뚀	뚜	뚸	뛔	뛰	뜌	뜨	띄	띠
	라	래	랴	럐	러	레	려	례	로	롸	뢔	뢰	료	루	뤄	뤠	뤼	류	르	릐	리
	마	매	먀	먜	머	메	며	몌	모	뫄	뫠	뫼	묘	무	뭐	뭬	뮈	뮤	므	믜	미
	바	배	뱌	뱨	버	베	벼	볘	보	봐	봬	뵈	뵤	부	붜	붸	뷔	뷰	브	븨	비
	빠	빼	뺘	뺴	뻐	뻬	뼈	뼤	뽀	뽜	뽸	뾔	뾰	뿌	뿨	쀄	쀠	쀼	쁘	쁴	삐
	사	새	샤	섀	서	세	셔	셰	소	솨	쇄	쇠	쇼	수	숴	쉐	쉬	슈	스	싀	시
	싸	쌔	쌰	썌	써	쎄	쎠	쎼	쏘	쏴	쐐	쐬	쑈	쑤	쒀	쒜	쒸	쓔	쓰	씌	씨
	아	애	야	얘	어	에	여	예	오	와	왜	외	요	우	워	웨	위	유	으	의	이
	자	재	쟈	쟤	저	제	져	졔	조	좌	좨	죄	죠	주	줘	줴	쥐	쥬	즈	즤	지
	짜	째	쨔	쨰	쩌	쩨	쪄	쪠	쪼	쫘	쫴	쬐	쬬	쭈	쭤	쮀	쮜	쮸	쯔	쯰	찌
	차	채	챠	챼	처	체	쳐	쳬	초	촤	쵀	최	쵸	추	춰	췌	취	츄	츠	츼	치
	카	캐	캬	컈	커	케	켜	켸	코	콰	쾌	쾨	쿄	쿠	쿼	퀘	퀴	큐	크	킈	키
	타	태	탸	턔	터	테	텨	톄	토	톼	퇘	퇴	툐	투	퉈	퉤	튀	튜	트	틔	티
	파	패	퍄	퍠	퍼	페	펴	폐	포	퐈	퐤	푀	표	푸	풔	풰	퓌	퓨	프	픠	피
	하	해	햐	햬	허	헤	혀	혜	호	화	홰	회	효	후	훠	훼	휘	휴	흐	희	히`

var hangul9Descriptor = "" +
	`ᆻ	ᆨ	ᆺ	ᆫ	ᆼ	ᇀ	ᆷ	ᆸ	ᆮ	ᆯ	ᆾ
	ᆻ9	ᆨ9	ᆺ9	ᆫ9	ᆼ9	ᇀ9	ᆷ9	ᆸ9	ᆮ9	ᆯ9	ᆾ9
	ᆿ	ᇁ	ᇂ	ᆽ	ퟹ	ᆩ	ퟦ	ퟍ	ᆪ	ᆬ	ᆭ
	ᆿ9	ᇁ9	ᇂ9	ᆽ9	ퟹ9	ᆩ9	ퟦ9	ퟍ9	ᆪ9	ᆬ9	ᆭ9
	ᆰ	ᆱ	ᆲ	ᆳ	ᆴ	ᆵ	ᆶ	ᆹ	ᇫ	ᇰ	ᇹ
	ᆰ9	ᆱ9	ᆲ9	ᆳ9	ᆴ9	ᆵ9	ᆶ9	ᆹ9	ᇫ9	ᇰ9	ᇹ9`

var combiningDescriptor = "" +
	`̀	́	̂	̃	̄	̅	̆	̇	̈
	̉	̊		̋	̌	̍	̎	̏	̐
	̑	̒	̓	̔		̕	̖		̗
	̘	̙	̚	̛	̜	̝	̞	̟	̠
	̡	̢		̣	̤	̥	̦	̧	̨
	̩	̪	̫	̬		̭	̮		̯
	̰	̱	̲	̳	̴	̵	̶	̷	̸
	̹	̺		̻		̼	̽	̾	̿
	̀	́	͂	̓	̈́	ͅ	͆		͇
	͈	͉	͊	͋	͌	͍	͎	͏	͐
	͑	͒		͓		͔	͕	͖	͗
	͘	͙	͚	͛	͜		͝	͞	͟
	͠	͡	͢	ͣ	ͤ	ͥ	ͦ	ͧ	ͨ
	ͩ	ͪ		ͫ	ͬ	ͭ	ͮ	ͯ	`

var armenianDescriptor = "" +
	`	Ա	Բ	Գ	Դ	Ե	Զ	Է	Ը	Թ	Ժ	Ի	Լ	Խ	Ծ	Կ
	Հ	Ձ	Ղ	Ճ	Մ	Յ	Ն	Շ	Ո	Չ	Պ	Ջ	Ռ	Ս	Վ	Տ
	Ր	Ց	Ւ	Փ	Ք	Օ	Ֆ			ՙ	՚	՛	՜	՝	՞	՟
		ա	բ	գ	դ	ե	զ	է	ը	թ	ժ	ի	լ	խ	ծ	կ
	հ	ձ	ղ	ճ	մ	յ	ն	շ	ո	չ	պ	ջ	ռ	ս	վ	տ
	ր	ց	ւ	փ	ք	օ	ֆ	և		։	֊			֍	֎	֏`

var UnicodeFont Font

func init() {
	(&UnicodeFont).Load("ascii.png", asciiDescriptor, "")
	(&UnicodeFont).Load("extendeda.png", extendedaDescriptor, "")
	(&UnicodeFont).Load("extendedb.png", extendedbDescriptor, "")
	(&UnicodeFont).Load("supplement.png", supplementDescriptor, "")
	(&UnicodeFont).Load("ipa.png", ipaDescriptor, "")
	(&UnicodeFont).Load("greek.png", greekDescriptor, "")
	(&UnicodeFont).Load("cyrillic.png", cyrillicDescriptor, "")
	(&UnicodeFont).Load("vietnamese.png", vietnameseDescriptor, "")
	(&UnicodeFont).Load("hangul0.png", hangul0Descriptor, "")
	(&UnicodeFont).Load("hangul1.png", hangul0Descriptor, "1")
	(&UnicodeFont).Load("hangul9.png", hangul9Descriptor, "")
	(&UnicodeFont).Multiply(hangul0Descriptor, "x", "1", hangul9Descriptor)
	Each(hangul0Descriptor, func(v string) error {
		const buf = "	\u11a8\u11a9\u11aa\u11ab\u11ac\u11ad\u11ae\u11af\u11b0\u11b1\u11b2" +
			"\u11b3\u11b4\u11b5\u11b6\u11b7\u11b8\u11b9\u11ba\u11bb\u11bc" +
			"\u11bd\u11be\u11bf\u11c0\u11c1\u11c2"
		for i := 1; i < 28; i++ {

			var target = string([]rune(v)[0] + rune(i))
			var bottom = string([]rune(buf)[i])

			//println(v, "|",  bottom + "x" + v)
			(&UnicodeFont).Alias(target, bottom+"x"+v)
		}

		return nil
	})
	err := (&UnicodeFont).Load("combining.png", combiningDescriptor, "")
	if err != nil {
		println(err.Error())
	}
	(&UnicodeFont).Multiply(combiningDescriptor, "", "", cyrillicDescriptor)
	(&UnicodeFont).Load("armenian.png", armenianDescriptor, "")

	(&UnicodeFont).Alias("\t", " ")
	(&UnicodeFont).Alias("", " ")
}
