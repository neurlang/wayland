package main

var extendedaDescriptor = `Ā	Ġ	Ŀ	ş	ā	ġ	ŀ	Š
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

var supplementDescriptor = `¡	±	À	Ï	à	ï
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

var extendedbDescriptor = `ǎ	Ǎ	ǳ	ǲ	ȟ	Ȟ	ǒ	Ǒ	ȑ	Ȑ	ȳ	Ȳ
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

var UnicodeFont Font

func init() {
	(&UnicodeFont).Load("ascii.png", asciiDescriptor)
	(&UnicodeFont).Load("extendeda.png", extendedaDescriptor)
	(&UnicodeFont).Load("extendedb.png", extendedbDescriptor)
	(&UnicodeFont).Load("supplement.png", supplementDescriptor)
}
