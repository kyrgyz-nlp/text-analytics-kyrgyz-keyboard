package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"sort"
	"strings"
	"unicode"
)

type pair struct {
	r     rune
	count uint64
}

type bp struct {
	k     [2]rune
	count uint64
}

type tp struct {
	k     [3]rune
	count uint64
}

func main() {
	path := flag.String("input", "all_texts.txt", "Path to input text file or '-' for stdin")
	topN := flag.Int("top", 50, "Number of most frequent characters/ngrams to print")
	cyrOnly := flag.Bool("cyrillic-only", false, "Count only Cyrillic letters (and build ngrams from them)")
	doBigram := flag.Bool("bigram", false, "Compute bigram stats (letters only when -cyrillic-only)")
	doTrigram := flag.Bool("trigram", false, "Compute trigram stats (letters only when -cyrillic-only)")
	writeTSV := flag.Bool("tsv", false, "Write full sorted lists to TSV files")
	outPrefix := flag.String("out-prefix", "stats", "Output file prefix for TSV files")
	stopwordsPath := flag.String("stopwords", "config/stopwords.txt", "Path to stopwords file (one token per line)")
	excludeWordsPath := flag.String("exclude-words", "config/foreign_words.txt", "Path to extra words to exclude (one token per line)")
	flag.Parse()

	var f *os.File
	if *path == "-" {
		f = os.Stdin
	} else {
		var err error
		f, err = os.Open(*path)
		if err != nil {
			fmt.Fprintf(os.Stderr, "open input: %v\n", err)
			os.Exit(1)
		}
		defer f.Close()
	}

	// Lowercase, then remove tokens from the stopwords + excluded words lists before counting.
	removeTokens, err := loadStopwords(*stopwordsPath)
	if err != nil {
		fmt.Fprintf(os.Stderr, "load stopwords: %v\n", err)
		os.Exit(1)
	}
	if *excludeWordsPath != "" {
		extra, err := loadStopwords(*excludeWordsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load exclude words: %v\n", err)
			os.Exit(1)
		}
		for k := range extra {
			removeTokens[k] = true
		}
	}

	counts := make(map[rune]uint64, 1024)
	var total uint64

	bigrams := make(map[[2]rune]uint64, 2048)
	trigrams := make(map[[3]rune]uint64, 4096)
	letters := make([]rune, 0, 3)

	reader := bufio.NewReaderSize(f, 1<<20)
	token := make([]rune, 0, 64)

	isCountable := func(r rune) bool {
		if *cyrOnly {
			return unicode.In(r, unicode.Cyrillic) && unicode.IsLetter(r)
		}
		return true
	}

	flushToken := func() {
		if len(token) == 0 {
			return
		}
		cleaned := make([]rune, 0, len(token))
		for _, r := range token {
			if unicode.IsLetter(r) {
				cleaned = append(cleaned, r)
			}
		}
		if len(cleaned) == 0 {
			token = token[:0]
			return
		}
		if removeTokens[string(cleaned)] {
			token = token[:0]
			return
		}
		for _, r := range cleaned {
			if !isCountable(r) {
				continue
			}
			counts[r]++
			total++
			if *doBigram || *doTrigram {
				letters = append(letters, r)
				if len(letters) > 3 {
					letters = letters[len(letters)-3:]
				}
				if *doBigram && len(letters) >= 2 {
					var key [2]rune
					key[0] = letters[len(letters)-2]
					key[1] = letters[len(letters)-1]
					bigrams[key]++
				}
				if *doTrigram && len(letters) >= 3 {
					var key [3]rune
					key[0] = letters[len(letters)-3]
					key[1] = letters[len(letters)-2]
					key[2] = letters[len(letters)-1]
					trigrams[key]++
				}
			}
		}
		token = token[:0]
	}

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			flushToken()
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "read rune: %v\n", err)
			os.Exit(1)
		}
		r = unicode.ToLower(r)

		if unicode.IsSpace(r) {
			flushToken()
			if !*cyrOnly {
				counts[r]++
				total++
			} else {
				// Reset ngram window across word boundaries in cyrillic-only mode.
				letters = letters[:0]
			}
			continue
		}
		token = append(token, r)
	}

	list := make([]pair, 0, len(counts))
	for r, c := range counts {
		list = append(list, pair{r: r, count: c})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].count == list[j].count {
			return list[i].r < list[j].r
		}
		return list[i].count > list[j].count
	})

	fmt.Printf("total_chars: %d\n", total)
	fmt.Printf("unique_chars: %d\n", len(list))
	fmt.Printf("top_%d_chars:\n", *topN)
	for i := 0; i < *topN && i < len(list); i++ {
		r := list[i].r
		display := string(r)
		if r == '\n' {
			display = "\\n"
		} else if r == '\t' {
			display = "\\t"
		} else if r == ' ' {
			display = "space"
		}
		percent := 0.0
		if total > 0 {
			percent = (float64(list[i].count) / float64(total)) * 100.0
		}
		fmt.Printf("%4d. %s\t%d\t(%.4f%%)\n", i+1, display, list[i].count, percent)
	}

	if *writeTSV {
		if err := writeCharTSV(*outPrefix, list, total); err != nil {
			fmt.Fprintf(os.Stderr, "write char tsv: %v\n", err)
			os.Exit(1)
		}
	}

	if *doBigram {
		bpList := make([]bp, 0, len(bigrams))
		var bTotal uint64
		for k, c := range bigrams {
			bpList = append(bpList, bp{k: k, count: c})
			bTotal += c
		}
		sort.Slice(bpList, func(i, j int) bool {
			if bpList[i].count == bpList[j].count {
				return string(bpList[i].k[:]) < string(bpList[j].k[:])
			}
			return bpList[i].count > bpList[j].count
		})
		fmt.Printf("top_%d_bigrams:\n", *topN)
		for i := 0; i < *topN && i < len(bpList); i++ {
			percent := 0.0
			if bTotal > 0 {
				percent = (float64(bpList[i].count) / float64(bTotal)) * 100.0
			}
			fmt.Printf("%4d. %s\t%d\t(%.4f%%)\n", i+1, string(bpList[i].k[:]), bpList[i].count, percent)
		}

		if *writeTSV {
			if err := writeBigramTSV(*outPrefix, bpList, bTotal); err != nil {
				fmt.Fprintf(os.Stderr, "write bigram tsv: %v\n", err)
				os.Exit(1)
			}
		}
	}

	if *doTrigram {
		tpList := make([]tp, 0, len(trigrams))
		var tTotal uint64
		for k, c := range trigrams {
			tpList = append(tpList, tp{k: k, count: c})
			tTotal += c
		}
		sort.Slice(tpList, func(i, j int) bool {
			if tpList[i].count == tpList[j].count {
				return string(tpList[i].k[:]) < string(tpList[j].k[:])
			}
			return tpList[i].count > tpList[j].count
		})
		fmt.Printf("top_%d_trigrams:\n", *topN)
		for i := 0; i < *topN && i < len(tpList); i++ {
			percent := 0.0
			if tTotal > 0 {
				percent = (float64(tpList[i].count) / float64(tTotal)) * 100.0
			}
			fmt.Printf("%4d. %s\t%d\t(%.4f%%)\n", i+1, string(tpList[i].k[:]), tpList[i].count, percent)
		}

		if *writeTSV {
			if err := writeTrigramTSV(*outPrefix, tpList, tTotal); err != nil {
				fmt.Fprintf(os.Stderr, "write trigram tsv: %v\n", err)
				os.Exit(1)
			}
		}
	}
}

func displayRune(r rune) string {
	switch r {
	case '\n':
		return "\\n"
	case '\t':
		return "\\t"
	case ' ':
		return "SPACE"
	default:
		return string(r)
	}
}

func writeCharTSV(prefix string, list []pair, total uint64) error {
	out, err := os.Create(prefix + "_chars.tsv")
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	_, _ = fmt.Fprintln(w, "rank\tchar\tcount\tpercent")
	for i, p := range list {
		percent := 0.0
		if total > 0 {
			percent = (float64(p.count) / float64(total)) * 100.0
		}
		_, _ = fmt.Fprintf(w, "%d\t%s\t%d\t%.6f\n", i+1, displayRune(p.r), p.count, percent)
	}
	return w.Flush()
}

func writeBigramTSV(prefix string, list []bp, total uint64) error {
	out, err := os.Create(prefix + "_bigrams.tsv")
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	_, _ = fmt.Fprintln(w, "rank\tbigram\tcount\tpercent")
	for i, p := range list {
		percent := 0.0
		if total > 0 {
			percent = (float64(p.count) / float64(total)) * 100.0
		}
		_, _ = fmt.Fprintf(w, "%d\t%s%s\t%d\t%.6f\n", i+1, displayRune(p.k[0]), displayRune(p.k[1]), p.count, percent)
	}
	return w.Flush()
}

func writeTrigramTSV(prefix string, list []tp, total uint64) error {
	out, err := os.Create(prefix + "_trigrams.tsv")
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	_, _ = fmt.Fprintln(w, "rank\ttrigram\tcount\tpercent")
	for i, p := range list {
		percent := 0.0
		if total > 0 {
			percent = (float64(p.count) / float64(total)) * 100.0
		}
		_, _ = fmt.Fprintf(w, "%d\t%s%s%s\t%d\t%.6f\n", i+1, displayRune(p.k[0]), displayRune(p.k[1]), displayRune(p.k[2]), p.count, percent)
	}
	return w.Flush()
}

func loadStopwords(path string) (map[string]bool, error) {
	f, err := os.Open(path)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	stop := make(map[string]bool, 64)
	scanner := bufio.NewScanner(f)
	for scanner.Scan() {
		line := strings.TrimSpace(scanner.Text())
		if line == "" || strings.HasPrefix(line, "#") {
			continue
		}
		line = strings.ToLower(line)
		stop[line] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return stop, nil
}
