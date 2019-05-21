package main

import (
	"fmt"
	"io/ioutil"
	"math"
	"strings"
	"unicode"
)

func main() {
	// f, err := os.Open(filepath.Join("input", "careless.kfn-instru.ogg"))
	// if err != nil {
	// 	log.Fatal(err)
	// }

	// streamer, format, err := vorbis.Decode(f)
	// if err != nil {
	// 	log.Fatal(err)
	// }
	// defer streamer.Close()

	// speaker.Init(format.SampleRate, format.SampleRate.N(time.Second/10))

	// done := make(chan bool)
	// speaker.Play(beep.Seq(streamer, beep.Callback(func() {
	// 	done <- true
	// })))

	printLyrics()

	// <-done
}

func printLyrics() {
	formtLyricsAsHTML()
	// ticker := time.NewTicker(10 * time.Millisecond)
	// secondsElapsed := 0
	// printNextLines(lineCount, 4)
	// go func() {
	// 	for range ticker.C {
	// 		secondsElapsed++
	// 		printWordAt(secondsElapsed)
	// 	}
	// }()
}

func formtLyricsAsHTML() {
	output := htmlTop

	output += "  <div id=\"lyrics\">\n"
	output += formatLyricsBlockAsHTML()
	output += "  </div>\n"

	// output += "  <div id=\"lyrics-background\">"
	// output += formatLyricsBlockAsHTML()
	// output += "  </div>"

	output += htmlBottom

	err := ioutil.WriteFile("html/index.html", []byte(output), 0644)
	if err != nil {
		panic(err)
	}
}

func formatLyricsBlockAsHTML() (output string) {
	timecodeCounter := 0
	output += "    <div class=\"page visible\">\n"
	twolineOffset := timeCodes[timecodeCounter]*10 - 1000
	output += fmt.Sprintf(
		"      <div class=\"two-lines\" data-appear-at=\"%d\">\n",
		twolineOffset,
	)
	for i, lyricLine := range lyrics {
		if i > 0 && i%2 == 0 && i%4 != 0 {
			output += "      </div>\n"
			twolineOffset = timeCodes[timecodeCounter]*10 - 1000
			output += fmt.Sprintf(
				"      <div class=\"two-lines\" data-appear-at=\"%d\">\n",
				twolineOffset,
			)
		}
		if i > 0 && i%4 == 0 {
			output += "      </div>\n"
			output += "    </div>\n"
			output += "    <div class=\"page\">\n"
			twolineOffset = timeCodes[timecodeCounter]*10 - 1000
			output += fmt.Sprintf(
				"      <div class=\"two-lines\" data-appear-at=\"%d\">\n",
				twolineOffset,
			)
		}
		lyricLine = strings.ReplaceAll(lyricLine, "/ ", "/_ ")
		lyricLine = strings.ReplaceAll(lyricLine, " /", "/_ ")
		dl := processLine(lyricLine)
		output += fmt.Sprintf("        ")
		for _, displayWord := range dl.words {
			delay := timeCodes[timecodeCounter]*10 - twolineOffset
			var duration int
			if timecodeCounter+1 < len(timeCodes) {
				duration = ((timeCodes[timecodeCounter+1] * 10) - twolineOffset) - delay
			} else {
				duration = 100
			}

			output += fmt.Sprintf("<span title=\"%s\" class=\"lyric\" style=\"", displayWord.output)
			output += fmt.Sprintf("animation-duration: %dms; animation-delay: %dms\">", duration, delay)
			output += fmt.Sprintf("%s</span>", displayWord.output)
			if displayWord.space {
				output += " "
			}
			timecodeCounter++
		}
		output += fmt.Sprintf("<br />\n")
	}
	output += "    </div>\n"
	return output
}

func printNextLines(lineCount int, numLines int) {
	for i := 0; i < numLines; i++ {
		lyricLine := lyrics[lineCount+i]
		lyricLine = strings.ReplaceAll(lyricLine, "/ ", "/_ ")
		lyricLine = strings.ReplaceAll(lyricLine, " /", "/_ ")
		dl := processLine(lyricLine)

		for _, displayWord := range dl.words {
			fmt.Print(displayWord.output)
			if displayWord.space {
				fmt.Print(" ")
			}
		}

		fmt.Print("\n")
	}
	moveCursorUpLines(numLines)
}

func moveCursorUpLines(numLines int) {
	for i := 0; i < numLines; i++ {
		fmt.Print("\033[1A")
	}
}

func printWordAt(timeCode int) {
	if timeCode == timeCodes[wordCount] {
		lyricLine := lyrics[lineCount]
		lyricLine = strings.ReplaceAll(lyricLine, "/ ", "/_ ")
		lyricLine = strings.ReplaceAll(lyricLine, " /", "/_ ")
		dl := processLine(lyricLine)

		fmt.Print(yellow(dl.words[wordInLineCount].output))
		if dl.words[wordInLineCount].space {
			fmt.Print(" ")
		}

		wordCount++
		wordInLineCount++

		if wordInLineCount >= len(dl.words) {
			wordInLineCount = 0
			fmt.Printf("\n")
			lineCount++
			if math.Mod(float64(lineCount), 4) == 0 {
				printNextLines(lineCount, 4)
			}
		}
	}
}

func yellow(message string) string {
	return fgYellow + message + reset
}

const fgYellow = "\x1b[33m"
const reset = "\x1b[0m"

func processLine(line string) (dl displayLine) {
	var tmpHolder string

	for _, rune := range line {
		if unicode.IsSpace(rune) {
			dl.addWord(tmpHolder)
			tmpHolder = ""
		} else if string(rune) == "/" {
			dl.addSyllable(tmpHolder)
			tmpHolder = ""
		} else if string(rune) == "_" {
			tmpHolder += ""
		} else {
			tmpHolder += string(rune)
		}
	}
	dl.addWord(tmpHolder)

	return dl
}

type displayWord struct {
	output string
	space  bool
}

type displayLine struct {
	words []displayWord
}

func (dl *displayLine) addSyllable(word string) {
	dl.words = append(dl.words, displayWord{output: word, space: false})
}

func (dl *displayLine) addWord(word string) {
	dl.words = append(dl.words, displayWord{output: word, space: true})
}

var wordCount = 0
var lineCount = 0
var wordInLineCount = 0

func splitRune(r rune) bool {
	return r == ' ' || r == '/'
}

var lyrics = []string{
	"I feel/ so /un/sure/_",
	"/as I take your hand/_",
	"and lead you to/_",
	"the/ dance floor/_",
	"/As the mu/sic dies/_",
	"some/thing in your eyes/_",
	"calls to mind the sil/ver screen/_",
	"and all its sad good-//byes/_",
	"I'm ne/ver go/nna dance a/gain/_",
	"guil/ty feet have got no rhy/thm/_",
	"Though it's ea/sy to pre/tend/_",
	"I know you're not a fool/_",
	"Shoul/d've known be/tter than",
	"to cheat a friend/_",
	"And waste the chance that I've/_",
	"been gi/ven/_",
	"So I'm ne/ver",
	"go/nna dance a/gain/_",
	"the way I danced with/ you/_",
	"/oh oh oh oh/_",
	"Time can ne/ver mend/_",
	"/the care/less whis/pers/_",
	"of a /good/ friend/_",
	"_",
	"To the heart and mind/_",
	"i/gno/rance is kind/_",
	"There's no com/fort in the truth/_",
	"pain is all you'll find/_",
	"I'm ne/ver go/nna dance a/gain/_",
	"guil/ty feet have got no rhy/thm/_",
	"Though it's ea/sy to pre/tend/_",
	"I know you're not a fool/_",
	"Should/'ve known be/tter than",
	"to cheat a friend/_",
	"And waste this chance that I've/_",
	"been gi/ven/_",
	"So I'm ne/ver",
	"go/nna dance a/gain/_",
	"the way I danced with you/_",
	"oh oh oh/_",
	"Ne/ver wi/thout your love/_",
	"To/night the mu/sic",
	"seems so loud/_",
	"I wish that we",
	"could lose this crowd/_",
	"May/be/ it's be/tter this way/_",
	"We'd hurt each o/ther",
	"with the things we want to say/_",
	"_",
	"We could have been",
	"so good to/ge/ther/_",
	"We could have lived",
	"this dance fo/re/ver/_",
	"But now",
	"who's go/nna dance with me?/_",
	"Please stay/_",
	"_",
	"And I'm ne/ver",
	"go/nna dance a/gain/_",
	"guil/ty feet have got no rhy/thm/_",
	"_",
	"_",
	"Though it's ea/sy to pre/tend/_",
	"I know you're not a fool/_",
	"_",
	"Should/'ve known be/tter than",
	"to cheat a friend/_",
	"And waste the chance that I've/_",
	"been gi/ven/_",
	"So I'm ne/ver",
	"go/nna dance a/gain/_",
	"the way I danced with/ you/_",
	"oh oh oh oh oh/_",
	"_",
	"Now/ that you're gone/_",
	"_",
	"Ah ah oh/_",
	"_",
	"_",
	"Was what I did so wrong/_",
	"so wrong/_",
	"that you had to leave/_",
	"me a/lone?/_",
}

var timeCodes = []int{
	2631, 2673, 2719, 2752, 2807, 2826, 2907, 3011, 3197, 3217, 3236, 3257, 3334, 3375, 3433, 3451, 3492, 3529, 3569, 3608, 3647, 3686, 3729, 3806, 3873, 3901, 3920, 3940, 3959, 3998, 4018, 4142, 4233, 4252, 4272, 4311, 4334, 4460, 4506, 4545, 4584,
	4628, 4665, 4701, 4723, 4768, 4805, 4824, 4897, 4975, 5054, 5101, 5116, 5167, 5190, 5210, 5230, 5249, 5269, 5295, 5328, 5348, 5426, 5446, 5485, 5524, 5544, 5583, 5642, 5681, 5721, 5761, 5799, 5819, 5838, 5878, 5917, 5956, 5976, 6034, 6054, 6074,
	6094, 6133, 6172, 6192, 6349, 6388, 6408, 6427, 6467, 6486, 6506, 6526, 6545, 6584, 6604, 6663, 6683, 6702, 6741, 6781, 6820, 6842, 6879, 6899, 6938, 6977, 7017, 7056, 7076, 7095, 7115, 7135, 7155, 7174, 7214, 7233, 7289, 7312, 7332, 7351, 7391,
	7430, 7468, 7489, 7841, 7883, 7909, 7925, 7967, 8009, 8158, 8907, 8985, 9025, 9103, 9181, 9299, 9494, 9513, 9533, 9611, 9650, 9708, 9795, 9844, 9922, 9961, 10000, 10059, 10060, 10132, 10176, 10195, 10215, 10234, 10273, 10292, 10387, 10507,
	10526, 10546, 10585, 10604, 10705, 10780, 10820, 10859, 10898, 10937, 10977, 10996, 11035, 11094, 11173, 11212, 11329, 11369, 11465, 11466, 11486, 11505, 11525, 11545, 11565, 11603, 11623, 11682, 11721, 11760, 11800, 11819, 11858, 11917,
	11956, 11996, 12036, 12074, 12094, 12113, 12152, 12192, 12231, 12251, 12309, 12329, 12349, 12368, 12407, 12447, 12466, 12643, 12662, 12682, 12702, 12741, 12760, 12780, 12800, 12819, 12858, 12878, 12937, 12956, 12976, 13015, 13054, 13094,
	13116, 13152, 13172, 13211, 13251, 13290, 13329, 13349, 13368, 13388, 13407, 13427, 13447, 13486, 13505, 13561, 13585, 13604, 13624, 13663, 13703, 13743, 13801, 13802, 13865, 13925, 14004, 14472, 14492, 14512, 14535, 14596, 14638, 14732,
	16430, 16453, 16490, 16513, 16549, 16595, 16648, 16672, 16728, 16747, 16772, 16795, 16832, 16871, 16911, 16973, 17008, 17062, 17084, 17127, 17162, 17188, 17224, 17262, 17282, 17321, 17377, 17381, 17400, 17420, 17440, 17460, 17480, 17499,
	17519, 17539, 17559, 17598, 17618, 17691, 17693, 17697, 17717, 17738, 17777, 17816, 17856, 17915, 17935, 17955, 17991, 18014, 18034, 18054, 18093, 18133, 18172, 18231, 18251, 18271, 18304, 18330, 18350, 18448, 18488, 18547, 18567, 18626,
	18646, 18818, 18823, 18891, 19001, 19010, 19021, 19040, 19060, 19079, 19099, 19118, 19139, 19178, 19198, 19257, 19297, 19336, 19376, 19395, 19435, 19494, 19533, 19573, 19613, 19626, 19633, 19652, 19672, 19691, 19731, 19770, 19810, 19829,
	19889, 19908, 19928, 19948, 19987, 20047, 20067, 20242, 20243, 20244, 20264, 20283, 20323, 20343, 20362, 20382, 20402, 20441, 20461, 20520, 20540, 20560, 20599, 20639, 20678, 20701, 20737, 20757, 20797, 20836, 20876, 20915, 20935, 20954,
	20974, 20994, 21014, 21033, 21073, 21093, 21142, 21172, 21191, 21211, 21251, 21290, 21328, 21349, 21725, 21744, 21764, 21799, 21837, 21854, 21913, 23152, 23164, 23292, 23321, 23380, 23399, 23538, 23949, 23969, 23988, 24008, 24063, 24081,
	24089, 24302, 24342, 24364, 24400, 24440, 24481, 24556, 24597, 24636, 24841, 24852, 24872, 24891, 24931, 24950, 25041, 25048, 25088, 25110, 25377,
}

const htmlTop = `<!DOCTYPE html>
	<html>
	<head>
	  <meta charset="UTF-8">
	  <link rel="stylesheet" type="text/css" href="site.css">  <title>Karaoke!</title>
	  <script src="https://ajax.googleapis.com/ajax/libs/jquery/3.4.0/jquery.min.js"></script>
	  <script>
		  function play() {
			var song = document.getElementById("song");
			var lyrics = document.getElementById("lyrics");
			song.play();
			lyrics.style.animationPlayState = "running";

			$('.two-lines').each(function (index) {
                $(this).delay($(this).data("appear-at")).fadeIn(500);
			});
		  };
	  </script>
	</head>
	<body>

	<div class="button" onclick="play();">Play!</div>
	<audio src="careless.ogg" id="song"></audio>

	<div id="lyrics-wrapper">
	`

const htmlBottom = `</div>
</body>
</html>`
