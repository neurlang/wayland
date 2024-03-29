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

var chinese1Descriptor = "" +
	`一	七	丈	三	上	下	不	且	世	丟	並	中	串	主	乃	久	之	乏	乖	乘	乙	九	也	乾	亂	了	予	事	二	于	云	互	五	井	些	亞	亟	交	亦	享	京	亮	人	什	介	仍
	仔	他	付	仙	代	令	以	仰	仲	件	任	份	仿	企	伊	休	伙	伯	估	伴	伸	似	但	佈	位	低	住	佔	何	佛	作	你	佣	佩	佳	併	使	來	例	供	依	侮	侵	侷	便	係
	促	俄	俗	保	俟	信	修	俱	倉	個	倍	們	倒	倘	候	借	倡	倫	值	假	偉	偏	做	停	健	偵	偶	偷	偽	傑	傘	備	傢	催	傳	債	傷	傻	傾	僅	像	僑	價	儀	億	儒
	儘	優	儲	允	元	兄	充	兇	先	光	克	免	兒	入	內	全	兩	八	公	六	共	兵	其	具	典	兼	冊	再	冒	冠	冰	冷	准	凝	凡	凸	出	刀	分	切	刊	刑	划	列	初	判
	別	利	刪	到	制	刺	刻	則	前	剖	剛	剝	剩	剪	副	割	創	劃	劇	劉	劍	力	功	加	助	努	勇	勉	勒	動	勘	務	勝	勞	募	勢	勤	勸	包	匈	化	北	匯	匹	區	十
	千	卅	升	午	半	卓	協	南	博	占	卡	印	危	即	卷	卻	厚	原	厲	去	參	又	及	友	反	叔	取	受	叛	叢	口	古	句	另	只	叫	召	可	台	史	右	司	吃	各	合	吋
	同	名	吐	向	君	否	吧	含	听	吳	吵	吸	吹	吾	呀	呆	呈	告	呎	呢	周	味	呵	呼	命	和	咖	咦	咬	咱	品	哇	哈	哎	員	哥	哦	哩	哪	哭	哲	哼	唇	唉	售	唯
	唱	唸	商	啊	問	啟	啤	啦	喂	善	喉	喊	喔	喘	喚	喜	喝	喪	單	喲	嗎	嗜	嗩	嗯	嘆	嘉	嘔	嘗	嘛	嘴	嘿	噁	噢	器	噪	噴	噸	嚇	嚮	嚴	嚼	囉	四	回	因	困
	固	圈	國	圍	園	圓	圖	團	土	在	地	均	坊	坐	坦	坪	垂	垃	型	埃	埋	城	執	培	基	堂	堅	堆	堪	報	場	塊	塑	塔	塗	塞	填	境	墓	增	墨	墮	墾	壓	壞	壟
	士	壯	壺	壽	夏	夕	外	多	夜	夠	夢	夥	大	天	太	夫	央	失	夸	夾	奇	奉	契	奔	套	奠	奧	奪	奮	女	奴	奶	她	好	如	妓	妙	妥	妨	妳	妹	妻	姊	始	姐	姑
	姓	委	姿	威	娃	娛	娶	婆	婚	婦	媒	媳	媽	嫁	嫌	嬉	嬰	子	孔	字	存	孝	孟	季	孤	孩	孫	學	它	宇	守	安	宋	完	宗	官	定	宛	宜	客	宣	室	宮	害	家	容
	宿	寂	寄	密	富	寒	察	寢	實	寧	審	寫	寬	寵	寶	寸	封	射	將	專	尊	尋	對	導	小	少	尖	尚	尤	就	尷	尺	尼	尾	尿	局	屁	居	屆	屋	屍	屏	展	屢	層	屬
	山	岩	島	崇	崗	崩	嶄	州	巡	巢	工	左	巧	巨	差	已	巴	巷	市	布	希	帛	帝	師	席	帳	帶	常	帽	幅	幕	幢	幫	干	平	年	幸	幹	幻	幼	幽	幾	床	序	底	店
	度	座	庫	庭	康	廁	廚	廟	廠	廢	廣	廳	延	建	廿	弄	弊	式	引	弟	弦	弱	張	強	彈	彌	彎	彗	彙	形	彩	彭	彰	影	彷	彼	往	征	待	很	律	後	徐	得	徘	從
	復	循	微	徵	德	徹	心	必	忍	志	忘	忙	忠	快	念	忽	怎	怕	思	急	性	怨	怪	恐	恢	恨	恭	息	恰	悄	悠	患	您	悲	悶	情	惟	惡	想	惹	愈	愉	意	愛	感	慈
	態	慎	慕	慘	慚	慢	慣	慶	慾	憂	憑	憤	憧	憲	懂	應	懲	懷	戀	成	我	戒	或	截	戰	戲	戴	戶	房	所	扇	手	才	打	托	扛	扭	扮	扶	批	找	承	技	抄	抉	把
	抑	抓	投	抗	折	抨	抬	抱	抵	抹	抽	拆	拉	拋	拍	拒	拓	拔	拖	拚	招	拜	拯	拷	拼	拿	持	指	按	挑	挖	挨	挫	振	挺	挽	捉	捏	捐	捕	捧	捨	捲	捷	掀	掃
	授	掉	掌	排	掙	掛	採	探	接	控	推	掩	措	描	提	插	揚	換	握	揭	揮	援	損	搖	搜	搞	搬	搭	搶	摔	摘	摧	摩	摸	撐	撒	撞	撤	撥	播	撰	撲	撿	擁	擅	擊
	擋	操	擔	據	擠	擦	擬	擴	擷	擺	攜	攝	攤	支	收	攸	改	攻	放	政	故	效	敏	救	敗	敘	教	敢	散	敲	整	敵	敷	數	文	斜	斤	新	斷	方	於	施	旁	旅	旋	族
	旗	既	日	早	旺	昂	昆	明	易	星	是	時	晉	晒	晚	普	景	晶	智	暗	暢	暨	暫	暴	曉	曝	曰	曲	更	書	曹	曾	替	最	會	月	有	朋	服	望	朝	期	木	未	末	本
	朱	朵	李	材	村	杜	束	杭	杯	東	板	枕	林	枚	果	枝	架	柏	某	柔	查	校	株	核	根	格	栽	桃	案	桌	桶	梅	條	棄	棒	棟	森	棲	棵	椅	植	楊	業	極	概	榕
	榮	構	槍	樁	樂	樓	標	模	樣	樹	橋	機	檔	檢	檳	櫥	櫻	權	欠	次	欣	欲	欸	欺	款	歌	歐	歡	止	正	此	步	武	歧	歲	歷	歸	歹	死	殖	殘	段	殷	殺	殼	殿
	毀	毆	母	每	毒	比	毛	毫	民	氣	氧	水	永	氾	求	汗	汙	江	污	決	汽	沈	沉	沒	沙	沮	河	油	治	沼	沾	沿	況	泉	法	泡	波	泥	注	泰	洋	洗	洛	洞	洪	洶
	活	洽	派	流	浪	浮	浴	海	消	涉	液	涵	涼	淘	淚	淡	淨	淪	深	淵	混	淹	淺	添	清	減	渡	測	港	渴	游	湖	湯	源	準	溝	溪	溫	滅	滋	滑	滴	滿	漁	漂	漏
	演	漠	漢	漫	漲	漸	潛	潮	潰	澄	澎	澳	激	濃	濕	濫	瀏	瀑	瀕	瀰	灌	火	灰	災	炎	炒	炫	炸	為	烈	烏	烤	無	焦	然	煉	煙	煤	照	煩	煮	熊	熟	熱	燃	燈
	燒	燕	營	燦	爆	爛	爬	爭	父	爸	爹	爺	爽	牆	片	版	牌	牙	牛	牟	牠	物	特	牽	犀	犧	犯	狀	狗	狩	狹	狼	猛	猜	猴	猶	獅	獎	獨	獲	獵	率	玉	王	玩	玫
	玻	珊	珍	班	現	球	理	琴	瑞	環	瓦	瓶	瓷	甄	甘	甚	甜	生	產	用	甫	田	由	甲	申	男	界	畔	留	畝	畢	略	番	畫	異	當	疏	疑	疫	疲	疼	疾	病	症	痕	痛
	痣	瘋	瘦	療	癌	登	發	白	百	的	皆	皇	皮	盆	盈	盒	盛	盞	盟	盡	監	盤	目	盯	盲	直	相	盼	省	眉	看	真	眷	眼	眾	睜	睡	督	瞧	瞪	瞭	矚	矛	矣	知	矩
	短	矮	石	砂	砍	研	破	砸	硬	碎	碗	碩	碰	確	碼	磚	磨	礦	示	社	祇	祈	祕	祖	祝	神	票	祭	禁	福	禪	禮	私	秉	科	秒	秘	租	秦	秩	移	稀	稅	程	稍	種
	稱	稽	稿	積	穩	究	空	穿	突	窄	窗	窮	立	站	竟	童	端	競	竹	笑	笙	笛	符	笨	第	筆	等	答	策	筷	算	管	箭	箱	節	範	篇	篩	簡	簽	籃	籌	籠	籤	米	粒
	粗	精	糊	糖	糟	糧	系	糾	紀	約	紅	納	紐	紓	純	紙	級	紛	素	紡	索	紫	紮	累	細	紳	終	組	結	絕	給	統	絲	綁	經	綜	綠	維	綱	網	綽	綿	緊	線	締	緣
	編	緩	練	縣	縮	縱	總	績	繁	繞	繩	繪	繫	繳	繼	續	纖	缺	罐	罕	罪	置	罰	罵	罷	罹	羅	羊	美	群	羨	義	羽	翅	習	翻	老	考	者	而	耐	耕	耗	耳	耶	聆
	聊	聖	聘	聚	聞	聯	聰	聲	職	聽	肉	肌	肚	肝	股	肢	肥	肩	肯	育	肺	胃	背	胖	胡	胸	能	脂	脆	脈	脖	脫	脾	腎	腐	腦	腫	腰	腳	腿	膨	臉	臨	自	臭	至
	致	臺	舅	與	興	舉	舊	舌	舒	舞	航	般	船	艘	良	艱	色	芝	芭	花	苗	若	苦	英	茁	茫	茶	草	荒	荷	莊	莫	菁	菜	菩	華	菲	萎	萬	落	葉	著	葡	董	蒐	蒙
	蒸	蓋	蓬	蓮	蔓	蔡	蔣	蔬	蕭	薄	薪	藉	藍	藏	藝	藥	蘇	蘋	蘭	虎	虐	虔	處	虛	號	蚊	蛇	蛋	蝴	蝸	螃	融	螞	螢	蟋	蟑	蟲	蠟	蠻	血	行	衍	街	衛	衝	衡
	衣	表	衰	被	裁	裏	補	裝	裡	製	複	褲	西	要	覆	見	規	覓	視	親	覺	觀	角	解	觸	言	訂	計	訊	討	訓	記	訝	訪	設	許	訴	診	註	評	詞	詢	試	詩	詮	話
	該	詳	誇	認	誕	誘	語	誠	誤	說	誰	課	調	談	請	諒	論	諮	諷	諸	諾	謀	謂	謊	謎	謙	講	謝	謹	證	識	譜	警	譬	議	譴	護	讀	變	讓	讚	豈	豐	象	豪	豬
	貓	貝	負	財	貢	貧	貨	販	貪	貫	責	貴	買	貸	費	貼	貿	資	賊	賓	賜	賞	賠	賣	賦	質	賭	賺	購	賽	贈	贊	贏	赤	走	赴	起	趁	超	越	趕	趙	趟	趣	趨	足
	跌	跑	距	跟	跡	跨	路	跳	踏	踢	踩	踴	蹤	蹲	蹺	身	躲	躺	車	軌	軍	軟	軸	較	載	輔	輕	輛	輝	輩	輪	輸	輻	輿	轉	轎	轟	辛	辣	辦	辨	辭	辯	農	迄	迅
	迎	近	返	迫	迴	迷	追	退	送	逃	透	逐	途	這	通	逛	速	造	逢	連	逮	週	進	逼	逾	遂	遇	遊	運	遍	過	道	達	違	遙	遠	遣	適	遭	遲	遴	遵	遷	選	遺	避
	邀	邁	還	邊	邏	那	郊	郝	部	郭	郵	都	鄉	鄭	鄰	配	酒	酗	酸	醒	醫	釋	里	重	野	量	釐	金	針	釣	鈔	鈣	銀	銅	銜	銷	鋪	鋼	錄	錢	錦	錯	錶	鍵	鎮	鏡
	鐘	鐵	鑑	鑰	鑽	長	門	閃	閉	開	閏	閒	間	閩	閱	闖	關	闢	防	阻	阿	附	陌	降	限	院	陣	除	陪	陰	陳	陶	陷	陸	陽	隊	階	隔	隕	障	隧	隨	隱	隻	雄	集
	雇	雌	雕	雖	雙	雜	雞	離	難	雨	雪	雲	零	雷	電	需	震	霧	露	靈	青	靜	非	靠	面	革	鞋	鞏	韓	音	響	頁	頂	項	順	須	預	頑	頒	頓	頗	領	頭	頻	顆	題
	額	顏	願	顛	類	顧	顯	風	颱	飄	飛	食	飯	飲	飼	飽	飾	餅	養	餐	餓	餘	館	餵	首	香	馬	駐	駕	駱	騎	騙	騷	驅	驕	驗	驚	骨	髒	體	高	鬆	鬥	鬧	鬼	魅
	鮑	鮮	鳥	鴨	鴿	鷹	鹽	鹿	麥	麵	麻	黃	黑	默	點	黨	鼓	鼻	齊	齣	龍	龐																								`

var devanagari1Descriptor = "" +
	`ऍ	ए	ऎ	ऐ	ऑ	ऒ	ओ	औ	क	ः
	ऄ	अ	आ	इ	ई	उ	ऊ	ऋ	ऌ	ं
	ख	ग	घ	ङ	च	छ	ज	झ	ञ	ट
	ठ	ड	ढ	ण	त	थ	द	ध	न	ऩ
	प	फ	ब	भ	म	य	र	ऱ	ल	ळ
	क़	ख़	ग़	ज़	ड	ढ	फ़	य़	ॠ	ॡ
	ऴ	व	ॸ	ष	स	ह	_?_	श	ो	ौ
	०	१	२	३	४	५	६	७	८	९
	ा	ऻ	ि	ी	े	ॅ	ॆ	ै	ॉ	ॊ
	़	ॎ	ऀ	ँ	ऺ	ऽ	ु	ू	ृ	ॄ
	ॏ	ॕ	ॖ	ॗ	ॢ	ॣ	ꣿ	॥	।	ॱ`

var devanagari2Descriptor = "" +
	`ऍ्	ए्	ऎ्	ऐ्	ऑ्	ऒ्	ओ्	औ्	क्	्
	ऄ्	अ्	आ्	इ	ई	उ्	ऊ्	ऋ्	ऌ्	ॐ
	ख्	ग्	घ्	ङ्	च्	छ्	ज्	झ्	ञ्	ट्
	ठ्	ड्	ढ्	ण्	त्	थ्	द्	ध्	न्	ऩ्
	प्	फ्	ब्	भ्	म्	य्	र्	ऱ्	ल्	ळ्
	क़्	ख़्	ग़्	ज़्	ड़्	ढ़्	फ़्	य़्	ॠ्	ॡ्
	ऴ्	व्	_?_	ष्	स्	ह्	_?_	श्`

var devanagari3Descriptor = "" +
	`ॳ	ॹ	ॺ	ॻ	ॼ	ॽ	ॾ	ॿ	ꣴ	ꣳ
	ꣵ	ꣶ	ꣷ	꣸	꣹	꣺	ꣻ	ꣽ	꣼	`

var spacingmodDescriptor = "" +
	`ʰ	ˀ	ː	ˠ	˰
	ʱ	ˁ	ˑ	ˡ	˱
	ʲ	˂	˒	ˢ	˲
	ʳ	˃	˓	ˣ	˳
	ʴ	˄	˔	ˤ	˴
	ʵ	˅	˕	˥	˵
	ʶ	ˆ	˖	˦	˶
	ʷ	ˇ	˗	˧	˷
	ʸ	ˈ	˘	˨	˸
	ʹ	ˉ	˙	˩	˹
	ʺ	ˊ	˚	˪	˺
	ʻ	ˋ	˛	˫	˻
	ʼ	ˌ	˜	ˬ	˼
	ʽ	ˍ	˝	˭	˽
	ʾ	ˎ	˞	ˮ	˾
	ʿ	ˏ	˟	˯	˿`

var UnicodeFont Font

func init() {
	(&UnicodeFont).Load("ascii.png", asciiDescriptor, "")
	(&UnicodeFont).Load("extendeda.png", extendedaDescriptor, "")
	(&UnicodeFont).Load("extendedb.png", extendedbDescriptor, "")
	(&UnicodeFont).Load("supplement.png", supplementDescriptor, "")
	(&UnicodeFont).Load("spacingmod.png", spacingmodDescriptor, "")
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

	(&UnicodeFont).Load("chinese1.jpg", chinese1Descriptor, "")

	(&UnicodeFont).Load("devanagari1.png", devanagari1Descriptor, "")
	(&UnicodeFont).Load("devanagari2.png", devanagari2Descriptor, "")
	(&UnicodeFont).Load("devanagari3.png", devanagari3Descriptor, "")
	(&UnicodeFont).Combine("ः", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ं", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ा", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ऻ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ि", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ी", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("े", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॅ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॆ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ै", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॉ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॊ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ो", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ौ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("़", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॎ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ऀ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ँ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ऺ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ु", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ू", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ृ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॄ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॏ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॕ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॖ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॗ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॢ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ॣ", devanagari1Descriptor, "")

	(&UnicodeFont).Combine("ों", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ें", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("़ा", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ो़", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("़ि", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("ूँ", devanagari1Descriptor, "")
	(&UnicodeFont).Combine("़ो", devanagari1Descriptor, "")

	(&UnicodeFont).Combine("ꣿ", "ए", "")
	(&UnicodeFont).Alias("ꣾ", "एꣿ")
	(&UnicodeFont).Alias("क़्", "क़्")
	(&UnicodeFont).Alias("ख़्", "ख़्")
	(&UnicodeFont).Alias("ग़्", "ग़्")
	(&UnicodeFont).Alias("ज़्", "ज़्")
	(&UnicodeFont).Alias("ड़्", "ड़्")
	(&UnicodeFont).Alias("ढ़्", "ढ़्")
	(&UnicodeFont).Alias("फ़्", "फ़्")
	(&UnicodeFont).Alias("य़्", "य़्")
	(&UnicodeFont).Alias("ड़", "ड़")
	(&UnicodeFont).Alias("ढ़", "ढ़")
	(&UnicodeFont).Alias("ॴ", "आऺ")
	(&UnicodeFont).Alias("ॶ", "अॖ")
	(&UnicodeFont).Alias("ॷ", "अॗ")
	(&UnicodeFont).Alias("ॵ", "अॏ")
	(&UnicodeFont).Alias("ॲ", "अॅ")
	(&UnicodeFont).Alias("ꣲ", "ँ")
	(&UnicodeFont).Alias("॰", "°")

	(&UnicodeFont).Alias("\t", " ")
	(&UnicodeFont).Alias("", " ")
}
