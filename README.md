# [Kloop Kyrgyz](https://ky.kloop.asia) crawler.

Crawled content is included in the sqlite3 DB.

# Credits:

-   Thanks to [Bektour Iskender](https://twitter.com/bektour), we were allowed to crawl and use Kloop articles.
-   Thanks to [Henriette Brand's awesome tutorial](https://blog.theodo.com/2019/01/data-scraping-scrapy-django-integration/). My previous attempts to couple Django and Scrapy were [ridiculous](https://github.com/kyrgyz-nlp/readthedocs_cleaned_projects_list/). Her article gave me hints on how to organize the project to make it possible to call Scrapy from within Django.

# Notes:

-   This is currently a work in progress (written in ~8 hours or so).
-   There are 30 934 articles crawled from 2011 to 2024 (as of September 8, 2024). Due to network failures some of the articles were not crawled.
-   It took more than 12 hours to crawl the articles, because of the more or less gentle crawler settings (I didn't want to stress kloop's servers):

```
CONCURRENT_REQUESTS = 3
DOWNLOAD_DELAY = 1
```

## Technology

-   Django 4 (for robust admin panel and awesome ORM)
-   Scrapy (for, surprise, scraping)
-   Other usefult libs (see the project's `requirements` file).

## Explore the corpus

Unpack [all_texts.txt.zip](https://github.com/kyrgyz-nlp/kloop-corpus/blob/main/all_texts.txt.zip) and have fun.

## Text analytics (character + n-gram stats)

There is a Go script that lowercases text and removes tokens listed in a stopwords file (default:
`config/stopwords.txt`, one token per line),
and computes character, bigram, and trigram statistics. It can output full sorted TSV files for
further analysis.

Script:
- `scripts/char_frequency_go.go`

Example usage:
```
go run scripts/char_frequency_go.go -input all_texts.txt -cyrillic-only -bigram -trigram -tsv -out-prefix kyr
```

Inputs:
- `-input /path/to/file.txt` (or `-input -` to read from stdin)
- `-stopwords /path/to/stopwords.txt` (optional; defaults to `config/stopwords.txt`)

### Steps taken (English)
1) Build a stopword list by running the word-frequency tool on the corpus.
2) Update `config/stopwords.txt` based on the most frequent words.
3) Re-run the character + n-gram tool to produce updated TSVs.

Example flow:
```
go run scripts/stopword_finder_go.go -input all_texts.txt -top 200 -cyrillic-only -out word_freq.tsv -tsv
go run scripts/char_frequency_go.go -input all_texts.txt -cyrillic-only -bigram -trigram -tsv -out-prefix kyr
```

Outputs (when `-tsv` is used):
- `kyr_chars.tsv`
- `kyr_bigrams.tsv`
- `kyr_trigrams.tsv`

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

## Stats and tools

Stats files (TSV):
- `kyr_chars.tsv` (character frequencies)
- `kyr_bigrams.tsv` (bigram frequencies)
- `kyr_trigrams.tsv` (trigram frequencies)
- `word_freq.tsv` (word frequencies)

Scripts:
- `scripts/char_frequency_go.go` (characters + n-grams; optional TSV output)
- `scripts/stopword_finder_go.go` (word frequency tool for stopword discovery)

## Stopword finder

This Go script finds the most frequent words (useful for stopword discovery).

Script:
- `scripts/stopword_finder_go.go`

Example usage:
```
go run scripts/stopword_finder_go.go -input all_texts.txt -top 100 -cyrillic-only -tsv -out word_freq.tsv
```

Inputs:
- `-input /path/to/file.txt` (or `-input -` to read from stdin)
- `-stopwords /path/to/stopwords.txt` (optional; remove known stopwords before counting)
- `-min-len 2` (optional; default 2)

### Кадамдар (Кыргызча)
1) Корпустагы эң көп колдонулган сөздөрдү табуу үчүн сөз жыштыгын эсептөөчү куралды иштетүү.
2) Натыйжага жараша `config/stopwords.txt` файлын жаңылоо.
3) Символдор жана n-граммалар боюнча статистиканы кайра эсептеп TSV файлдарын чыгаруу.

Мисал:
```
go run scripts/stopword_finder_go.go -input all_texts.txt -top 200 -cyrillic-only -out word_freq.tsv -tsv
go run scripts/char_frequency_go.go -input all_texts.txt -cyrillic-only -bigram -trigram -tsv -out-prefix kyr
```
## Dev setup prerequisites

We assume you have the following packages are installed in your system:

-   git
-   Python 3.6 or above
-   venv

## Install dependencies and run

-   Clone the project:
    `git clone https://github.com/kyrgyz-nlp/kloop-corpus.git`

-   Go to the project folder:
    `cd kloop-corpus`

-   Create and activate a virtual env (sorry, Windows users, I don't know how you do this on your machine):
    `python -m venv env && source env/bin/activate`

-   Install project dependecies:
    `pip install -r requirements.txt`

-   Run the server:
    `./manage.py runserver`

-   To run the crawler
    `python manage.py crawl`

## Changelog
-   Sep, 08 2024: re-crawled from the beginning to update all_texts.txt.zip because the articles didn't contain valuable metadata.

## TODO

-   [x] Introduce `--start-year=2020` kind of args to the crawler
-   [ ] Extract metadata using NER models or LLM
-   [ ] Remove whitespaces and remove empty articles
-   [ ] Push the current version of the corpus to Hugging Face
-   [ ] Introduce --start-from='2020-04' kind of args to the crawler
-   [ ] Introduce upsert logic: if the article is not in the DB, then crawl and save
-   [ ] Add webpages with basic corpus statistics: `frequency dictionary`, `most frequent n-grams` etc.

## Citing

If you are using this dataset in your work, please cite it.
```bibtex
@misc{kyrgyz-nlp_kloop-corpus,
  author = {{kyrgyz-nlp}},
  title = {Kloop Corpus},
  year = {2024},
  publisher = {GitHub},
  journal = {GitHub repository},
  howpublished = {\url{https://github.com/kyrgyz-nlp/kloop-corpus}},
}
```
