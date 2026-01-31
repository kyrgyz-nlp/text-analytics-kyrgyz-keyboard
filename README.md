# Kyrgyz Text Analytics for Keyboard Layouts

This repo contains corpus-based statistics and tooling to justify a Kyrgyz-first keyboard layout.
It focuses on character frequencies, n-grams, and a data-backed layout rationale.

## Stats and tools

Stats files (TSV):
- `kyr_chars.tsv` (character frequencies)
- `kyr_bigrams.tsv` (bigram frequencies)
- `kyr_trigrams.tsv` (trigram frequencies)
- `raw_chars.tsv` (raw character frequencies, no stopwords)
- `raw_bigrams.tsv` (raw bigram frequencies)
- `raw_trigrams.tsv` (raw trigram frequencies)
- `word_freq.tsv` (word frequencies)

Scripts:
- `scripts/char_frequency_go.go` (characters + n-grams; optional TSV output)
- `scripts/stopword_finder_go.go` (word frequency tool for stopword discovery)

Stopwords:
- `config/stopwords.txt` (one token per line)

## Text analytics (character + n-gram stats)

The character tool lowercases text and removes tokens listed in a stopwords file (default:
`config/stopwords.txt`, one token per line), then computes character, bigram, and trigram
statistics. It can output full sorted TSV files for analysis.

Example usage:
```
go run scripts/char_frequency_go.go -input all_texts.txt -cyrillic-only -bigram -trigram -tsv -out-prefix kyr
```

Inputs:
- `-input /path/to/file.txt` (or `-input -` to read from stdin)
- `-stopwords /path/to/stopwords.txt` (optional; defaults to `config/stopwords.txt`)

## Stopword finder

This Go script finds the most frequent words (useful for stopword discovery).

Example usage:
```
go run scripts/stopword_finder_go.go -input all_texts.txt -top 200 -cyrillic-only -out word_freq.tsv -tsv
```

Inputs:
- `-input /path/to/file.txt` (or `-input -` to read from stdin)
- `-stopwords /path/to/stopwords.txt` (optional; remove known stopwords before counting)
- `-min-len 2` (optional; default 2)

### Steps taken (English)
1) Build a stopword list by running the word-frequency tool on the corpus.
2) Update `config/stopwords.txt` based on the most frequent words.
3) Re-run the character + n-gram tool to produce updated TSVs.

Example flow:
```
go run scripts/stopword_finder_go.go -input all_texts.txt -top 200 -cyrillic-only -out word_freq.tsv -tsv
go run scripts/char_frequency_go.go -input all_texts.txt -cyrillic-only -bigram -trigram -tsv -out-prefix kyr
```

### Кадамдар (Кыргызча)
1) Корпустагы эң көп колдонулган сөздөрдү табуу үчүн сөз жыштыгын эсептөөчү куралды иштетүү.
2) Натыйжага жараша `config/stopwords.txt` файлын жаңылоо.
3) Символдор жана n-граммалар боюнча статистиканы кайра эсептеп TSV файлдарын чыгаруу.

Мисал:
```
go run scripts/stopword_finder_go.go -input all_texts.txt -top 200 -cyrillic-only -out word_freq.tsv -tsv
go run scripts/char_frequency_go.go -input all_texts.txt -cyrillic-only -bigram -trigram -tsv -out-prefix kyr
```

## Kyrgyz keyboard layout rationale

This corpus shows that Kyrgyz-specific letters are frequent enough to deserve their own visible
places on the keyboard. Two of them, **ү** and **ө**, are high-frequency letters and should not be
hidden. **ң** is also more common than several Russian-only letters in Kyrgyz text, so it belongs
on the main layout as well.

Proposed placement decisions:
- **ү** takes the place of **ц** (ц becomes a long-press on **с**).
- **ө** takes the place of **щ** (щ becomes a long-press on **ш**).
- **ң** takes the place of **ф** (ф becomes a long-press on **в**).

ASCII layout (letters only):
```
й у к е н г ш з х
ф ы в а п р о л д ж э
я ч с м и т ь б ю
```

## Kyrgyz keyboard layout (proposed)
```
й ү у к е ё н ң г ш ө з х
ф ы в а п р о л д ж э
я ч с м и т ь б ю
```

Notes:
- Long-presses (examples): с→ц, ш→щ, в→ф
- Keep digits and punctuation as in the existing keyboard.

## Кыргызча негиздеме

Бул корпустун статистикасы кыргыз тилиндеги өзгөчө тамгалар өзүнчө көрүнүктүү орунга татыктуу
экенин көрсөтөт. **ү** жана **ө** эң көп колдонулган тамгалардын катарына кирет, ошондуктан аларды
жашыруу туура эмес. **ң** тамгасы да көп колдонулат жана кыргыз текстинде айрым орусча тамгалардан
да көп учурайт.

Сунушталган орун алмаштыруулар:
- **ү** тамгасы **ц** ордун алат (ц — **с** тамгасынын long‑press варианты).
- **ө** тамгасы **щ** ордун алат (щ — **ш** тамгасынын long‑press варианты).
- **ң** тамгасы **ф** ордун алат (ф — **в** тамгасынын long‑press варианты).

Эскертүү:
- long‑press мисалдары: с→ц, ш→щ, в→ф
- Сандар жана белгилер азыркыдай калат.
