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
var UnicodeFont Font

func init() {
	(&UnicodeFont).Load("ascii.png", asciiDescriptor)
	(&UnicodeFont).Load("extendeda.png", extendedaDescriptor)
	(&UnicodeFont).Load("extendedb.png", extendedbDescriptor)
	(&UnicodeFont).Load("supplement.png", supplementDescriptor)
	(&UnicodeFont).Load("ipa.png", ipaDescriptor)
	(&UnicodeFont).Load("greek.png", greekDescriptor)
	(&UnicodeFont).Load("cyrillic.png", cyrillicDescriptor)
	(&UnicodeFont).Load("vietnamese.png", vietnameseDescriptor)
	(&UnicodeFont).Alias("\t", " ")
	(&UnicodeFont).Alias("", " ")
}
