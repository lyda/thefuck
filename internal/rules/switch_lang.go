package rules

import (
	"strings"

	"github.com/lyda/thefuck/internal/types"
)

// targetLayout is the English QWERTY layout characters.
const targetLayout = `qwertyuiop[]asdfghjkl;'zxcvbnm,./QWERTYUIOP{}ASDFGHJKL:"ZXCVBNM<>?`

// sourceLayouts maps non-English keyboard layouts to targetLayout positions.
var sourceLayouts = []string{
	// Russian
	`йцукенгшщзхъфывапролджэячсмитьбю.ЙЦУКЕНГШЩЗХЪФЫВАПРОЛДЖЭЯЧСМИТЬБЮ,`,
	// Ukrainian
	`йцукенгшщзхїфівапролджєячсмитьбю.ЙЦУКЕНГШЩЗХЇФІВАПРОЛДЖЄЯЧСМИТЬБЮ,`,
	// Arabic
	"ضصثقفغعهخحجچشسیبلاتنمکگظطزرذدپو./ًٌٍَُِّْ][}{ؤئيإأآة»«:؛كٓژٰ\u200cٔء><؟",
	// Hebrew
	`/'קראטוןםפ][שדגכעיחלךף,זסבהנמצתץ.QWERTYUIOP{}ASDFGHJKL:"ZXCVBNM<>?`,
	// Greek
	`;ςερτυθιοπ[]ασδφγηξκλ΄ζχψωβνμ,./:΅ΕΡΤΥΘΙΟΠ{}ΑΣΔΦΓΗΞΚΛ¨"ΖΧΨΩΒΝΜ<>?`,
	// Korean
	`ㅂㅈㄷㄱㅅㅛㅕㅑㅐㅔ[]ㅁㄴㅇㄹㅎㅗㅓㅏㅣ;'ㅋㅌㅊㅍㅠㅜㅡ,./ㅃㅉㄸㄲㅆㅛㅕㅑㅒㅖ{}ㅁㄴㅇㄹㅎㅗㅓㅏㅣ:"ㅋㅌㅊㅍㅠㅜㅡ<>?`,
}

// greekToTarget handles the Greek layout which has different character counts
// (accented characters). Built from the Python source_to_target[greek] dict.
var greekToTarget = map[rune]string{
	';': "q", 'ς': "w", 'ε': "e", 'ρ': "r", 'τ': "t", 'υ': "y",
	'θ': "u", 'ι': "i", 'ο': "o", 'π': "p", '[': "[", ']': "]",
	'α': "a", 'σ': "s", 'δ': "d", 'φ': "f", 'γ': "g", 'η': "h",
	'ξ': "j", 'κ': "k", 'λ': "l", '΄': "'", 'ζ': "z", 'χ': "x",
	'ψ': "c", 'ω': "v", 'β': "b", 'ν': "n", 'μ': "m", ',': ",",
	'.': ".", '/': "/", ':': "Q", '΅': "W", 'Ε': "E", 'Ρ': "R",
	'Τ': "T", 'Υ': "Y", 'Θ': "U", 'Ι': "I", 'Ο': "O", 'Π': "P",
	'{': "{", '}': "}", 'Α': "A", 'Σ': "S", 'Δ': "D", 'Φ': "F",
	'Γ': "G", 'Η': "H", 'Ξ': "J", 'Κ': "K", 'Λ': "L", '¨': ":",
	'"': `"`, 'Ζ': "Z", 'Χ': "X", 'Ψ': "C", 'Ω': "V", 'Β': "B",
	'Ν': "N", 'Μ': "M", '<': "<", '>': ">", '?': "?",
	// Accented vowels
	'ά': "a", 'έ': "e", 'ύ': "y", 'ί': "i", 'ό': "o", 'ή': "h", 'ώ': "v",
	'Ά': "A", 'Έ': "E", 'Ύ': "Y", 'Ί': "I", 'Ό': "O", 'Ή': "H", 'Ώ': "V",
}

// greekLayout is the source layout string for Greek (used for identity comparison).
const greekLayout = `;ςερτυθιοπ[]ασδφγηξκλ΄ζχψωβνμ,./:΅ΕΡΤΥΘΙΟΠ{}ΑΣΔΦΓΗΞΚΛ¨"ΖΧΨΩΒΝΜ<>?`

// Korean decomposition data
var (
	headList = []rune{'ㄱ', 'ㄲ', 'ㄴ', 'ㄷ', 'ㄸ', 'ㄹ', 'ㅁ', 'ㅂ', 'ㅃ', 'ㅅ', 'ㅆ',
		'ㅇ', 'ㅈ', 'ㅉ', 'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
	bodyList = []rune{'ㅏ', 'ㅐ', 'ㅑ', 'ㅒ', 'ㅓ', 'ㅔ', 'ㅕ', 'ㅖ', 'ㅗ', 'ㅘ', 'ㅙ',
		'ㅚ', 'ㅛ', 'ㅜ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅠ', 'ㅡ', 'ㅢ', 'ㅣ'}
	tailList = []rune{' ', 'ㄱ', 'ㄲ', 'ㄳ', 'ㄴ', 'ㄵ', 'ㄶ', 'ㄷ', 'ㄹ', 'ㄺ', 'ㄻ',
		'ㄼ', 'ㄽ', 'ㄾ', 'ㄿ', 'ㅀ', 'ㅁ', 'ㅂ', 'ㅄ', 'ㅅ', 'ㅆ', 'ㅇ', 'ㅈ',
		'ㅊ', 'ㅋ', 'ㅌ', 'ㅍ', 'ㅎ'}
	doubleList    = []rune{'ㅘ', 'ㅙ', 'ㅚ', 'ㅝ', 'ㅞ', 'ㅟ', 'ㅢ', 'ㄳ', 'ㄵ', 'ㄶ', 'ㄺ', 'ㄻ', 'ㄼ', 'ㄽ', 'ㄾ', 'ㅀ', 'ㅄ'}
	doubleModList = []string{"ㅗㅏ", "ㅗㅐ", "ㅗㅣ", "ㅜㅓ", "ㅜㅔ", "ㅜㅣ", "ㅡㅣ", "ㄱㅅ", "ㄴㅈ", "ㄴㅎ", "ㄹㄱ", "ㄹㅁ", "ㄹㅂ", "ㄹㅅ", "ㄹㅌ", "ㄹㅎ", "ㅂㅅ"}
)

// isKorean returns true if ch is a Korean character (jamo or syllable).
func isKorean(ch rune) bool {
	return (ch >= 'ㄱ' && ch <= 'ㅎ') || (ch >= 'ㅏ' && ch <= 'ㅣ') || (ch >= '가' && ch <= '힣')
}

// targetLayoutRunes is the pre-computed rune slice for targetLayout.
var targetLayoutRunes = []rune(targetLayout)

// switchChar converts a single character from a source layout to target.
func switchChar(ch rune, layout string) rune {
	layoutRunes := []rune(layout)
	for i, r := range layoutRunes {
		if r == ch {
			if i < len(targetLayoutRunes) {
				return targetLayoutRunes[i]
			}
		}
	}
	return ch
}

// getMatchedLayout returns the source layout that matches all characters in the command.
func getMatchedLayout(script string) string {
	parts := strings.Split(script, " ")
	for _, layout := range sourceLayouts {
		allMatch := true
		for _, part := range parts {
			for _, ch := range part {
				if ch == '-' || ch == '_' {
					continue
				}
				if !strings.ContainsRune(layout, ch) {
					allMatch = false
					break
				}
			}
			if !allMatch {
				break
			}
		}
		if allMatch {
			return layout
		}
	}
	return ""
}

// switchCommand converts the command script from a source layout to English.
func switchCommand(script, layout string) string {
	if layout == greekLayout {
		var sb strings.Builder
		for _, ch := range script {
			if mapped, ok := greekToTarget[ch]; ok {
				sb.WriteString(mapped)
			} else {
				sb.WriteRune(ch)
			}
		}
		return sb.String()
	}

	var sb strings.Builder
	for _, ch := range script {
		sb.WriteRune(switchChar(ch, layout))
	}
	return sb.String()
}

// changeDouble replaces double Korean characters with their decomposed form.
func changeDouble(ch rune) string {
	for i, d := range doubleList {
		if ch == d {
			return doubleModList[i]
		}
	}
	return string(ch)
}

// decomposeKorean decomposes Korean syllable characters into jamo components.
func decomposeKorean(script string) string {
	var sb strings.Builder
	base := rune('가')
	for _, ch := range script {
		if ch >= '가' && ch <= '힣' {
			ordCh := int(ch - base)
			hd := ordCh / 588
			bd := (ordCh - 588*hd) / 28
			tl := ordCh - 588*hd - 28*bd
			for _, component := range []rune{headList[hd], bodyList[bd], tailList[tl]} {
				if component != ' ' {
					sb.WriteString(changeDouble(component))
				}
			}
		} else {
			sb.WriteString(changeDouble(ch))
		}
	}
	return sb.String()
}

func init() {
	register(Rule{
		Name: "switch_lang",
		Match: func(cmd types.Command) bool {
			if !strings.Contains(cmd.Output, "not found") {
				return false
			}
			// Check for Korean characters
			for _, ch := range cmd.Script {
				if isKorean(ch) {
					return true
				}
			}
			// Check if all characters match a known source layout
			layout := getMatchedLayout(cmd.Script)
			if layout == "" {
				return false
			}
			// Make sure switching would actually change something
			switched := switchCommand(cmd.Script, layout)
			return switched != cmd.Script
		},
		GetNewCommand: func(cmd types.Command) []types.CorrectedCommand {
			script := cmd.Script
			// Handle Korean decomposition first
			hasKorean := false
			for _, ch := range script {
				if isKorean(ch) {
					hasKorean = true
					break
				}
			}
			if hasKorean {
				script = decomposeKorean(script)
			}
			layout := getMatchedLayout(script)
			if layout == "" {
				return nil
			}
			return single(switchCommand(script, layout))
		},
	})
}
