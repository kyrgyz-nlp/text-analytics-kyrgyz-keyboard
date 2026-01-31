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

type wordCount struct {
	word  string
	count uint64
}

func main() {
	path := flag.String("input", "all_texts.txt", "Path to input text file or '-' for stdin")
	topN := flag.Int("top", 50, "Number of most frequent words to print")
	cyrOnly := flag.Bool("cyrillic-only", true, "Count only Cyrillic words")
	minLen := flag.Int("min-len", 2, "Minimum token length")
	stopwordsPath := flag.String("stopwords", "", "Optional stopwords file (one token per line)")
	writeTSV := flag.Bool("tsv", false, "Write full sorted list to a TSV file")
	outFile := flag.String("out", "word_freq.tsv", "Output TSV file when -tsv is set")
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

	stopwords := map[string]bool{}
	if *stopwordsPath != "" {
		var err error
		stopwords, err = loadStopwords(*stopwordsPath)
		if err != nil {
			fmt.Fprintf(os.Stderr, "load stopwords: %v\n", err)
			os.Exit(1)
		}
	}

	counts := make(map[string]uint64, 1<<16)
	var total uint64

	reader := bufio.NewReaderSize(f, 1<<20)
	token := make([]rune, 0, 64)

	isAllowedLetter := func(r rune) bool {
		if *cyrOnly {
			return unicode.In(r, unicode.Cyrillic) && unicode.IsLetter(r)
		}
		return unicode.IsLetter(r)
	}

	flush := func() {
		if len(token) == 0 {
			return
		}
		word := strings.ToLower(string(token))
		if len([]rune(word)) < *minLen {
			token = token[:0]
			return
		}
		if stopwords[word] {
			token = token[:0]
			return
		}
		counts[word]++
		total++
		token = token[:0]
	}

	for {
		r, _, err := reader.ReadRune()
		if err == io.EOF {
			flush()
			break
		}
		if err != nil {
			fmt.Fprintf(os.Stderr, "read rune: %v\n", err)
			os.Exit(1)
		}
		if isAllowedLetter(unicode.ToLower(r)) {
			token = append(token, unicode.ToLower(r))
			continue
		}
		flush()
	}

	list := make([]wordCount, 0, len(counts))
	for w, c := range counts {
		list = append(list, wordCount{word: w, count: c})
	}
	sort.Slice(list, func(i, j int) bool {
		if list[i].count == list[j].count {
			return list[i].word < list[j].word
		}
		return list[i].count > list[j].count
	})

	fmt.Printf("total_tokens: %d\n", total)
	fmt.Printf("unique_words: %d\n", len(list))
	fmt.Printf("top_%d_words:\n", *topN)
	for i := 0; i < *topN && i < len(list); i++ {
		percent := 0.0
		if total > 0 {
			percent = (float64(list[i].count) / float64(total)) * 100.0
		}
		fmt.Printf("%4d. %s\t%d\t(%.4f%%)\n", i+1, list[i].word, list[i].count, percent)
	}

	if *writeTSV {
		if err := writeTSVFile(*outFile, list, total); err != nil {
			fmt.Fprintf(os.Stderr, "write tsv: %v\n", err)
			os.Exit(1)
		}
	}
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
		stop[strings.ToLower(line)] = true
	}
	if err := scanner.Err(); err != nil {
		return nil, err
	}
	return stop, nil
}

func writeTSVFile(path string, list []wordCount, total uint64) error {
	out, err := os.Create(path)
	if err != nil {
		return err
	}
	defer out.Close()

	w := bufio.NewWriter(out)
	_, _ = fmt.Fprintln(w, "rank\tword\tcount\tpercent")
	for i, p := range list {
		percent := 0.0
		if total > 0 {
			percent = (float64(p.count) / float64(total)) * 100.0
		}
		_, _ = fmt.Fprintf(w, "%d\t%s\t%d\t%.6f\n", i+1, p.word, p.count, percent)
	}
	return w.Flush()
}
