#!/usr/bin/env python3
"""
Calculate character frequency percentages from literary source data.
Data extracted from: https://user-images.githubusercontent.com/155511/222278874-8d3c3ccb-79d9-4f33-b49e-97f775d8565c.png
"""

from pathlib import Path


def combine_sources(*sources: dict[str, int]) -> dict[str, int]:
    """Combine multiple sources by summing counts."""
    combined: dict[str, int] = {}
    for source in sources:
        for char, count in source.items():
            combined[char] = combined.get(char, 0) + count
    return combined


def load_tsv(filepath: str) -> dict[str, int]:
    """Load character counts from TSV file (rank, char, count, percent format)."""
    data: dict[str, int] = {}
    path = Path(filepath)
    with path.open(encoding="utf-8") as f:
        for line in f:
            line = line.strip()
            if not line or line.startswith("rank"):
                continue
            parts = line.split("\t")
            if len(parts) >= 3:
                char = parts[1]
                count = int(parts[2])
                data[char] = count
    return data


# Raw character counts from three sources (К and Қ merged)
MANAS = {
    'а': 143178,
    'к': 147442,  # К + Қ merged: 73721 + 46295 + 27426
    'н': 67181,
    'ы': 65134,
    'е': 59077,
    'т': 53825,
    'р': 53417,
    'л': 49676,
    'у': 40487,
    'д': 38762,
    'б': 38570,
    'п': 35414,
    'о': 32255,
    'и': 30206,
    'й': 28663,
    'г': 25216,
    'м': 24396,
    'ө': 23702,
    'ү': 22864,
    'с': 21863,
    'ж': 17959,
    'ч': 5514,
    'ш': 15281,
    'ң': 10576,
    'х': 1055,
    'э': 9822,
    'в': 14,
    'ф': 0,
}

SYNGAN_KYLYCH = {
    'а': 125337,
    'к': 138216,  # К + Қ merged: 69108 + 41662 + 27446
    'н': 67917,
    'ы': 61871,
    'е': 58713,
    'т': 55889,
    'р': 52744,
    'л': 47439,
    'у': 42967,
    'д': 42967,
    'б': 40319,
    'п': 35364,
    'о': 33342,
    'и': 30378,
    'й': 25572,
    'г': 25089,
    'ө': 27644,
    'ү': 22469,
    'с': 22033,
    'ж': 18418,
    'ч': 15458,
    'ш': 15140,
    'ң': 13699,
    'х': 7430,
    'э': 1918,
    'в': 766,
    'ф': 244,
    'м': 15665,
}

MISMILDIRIK = {
    'а': 43926,
    'н': 25672,
    'к': 22321,
    'ы': 21122,
    'т': 20210,
    'е': 20038,
    'р': 18142,
    'л': 15503,
    'у': 15314,
    'д': 14480,
    'о': 12799,
    'и': 12351,
    'б': 10585,
    'ө': 10108,
    'п': 8942,
    'ү': 8267,
    'й': 7991,
    'с': 7639,
    'м': 6852,
    'ш': 6535,
    'ж': 5567,
    'ч': 5228,
    'э': 4013,
    'ң': 2029,
    'в': 479,
    'х': 139,
    'ф': 78,
}


def calculate_percentages(data: dict[str, int]) -> dict[str, float]:
    """Calculate percentage for each character."""
    total = sum(data.values())
    return {char: (count / total) * 100 for char, count in data.items()}


def print_stats(name: str, data: dict[str, int]) -> None:
    """Print character statistics sorted by frequency."""
    percentages = calculate_percentages(data)
    total = sum(data.values())

    print(f"\n{'=' * 60}")
    print(f"{name}")
    print(f"Total characters: {total:,}")
    print(f"{'=' * 60}")
    print(f"{'Rank':<6}{'Char':<6}{'Count':>12}{'Percent':>12}")
    print(f"{'-' * 36}")

    sorted_chars = sorted(percentages.items(), key=lambda x: x[1], reverse=True)
    for rank, (char, pct) in enumerate(sorted_chars, 1):
        count = data[char]
        print(f"{rank:<6}{char:<6}{count:>12,}{pct:>11.3f}%")


def print_comparison() -> None:
    """Print side-by-side comparison of all sources including combined and corpus."""
    # Load corpus raw data
    script_dir = Path(__file__).parent
    corpus_raw = load_tsv(script_dir.parent / "raw_chars.tsv")

    # Combine literary sources
    combined = combine_sources(MANAS, SYNGAN_KYLYCH, MISMILDIRIK)

    sources = [
        ("Манас", MANAS),
        ("С.Кылыч", SYNGAN_KYLYCH),
        ("Мисмил.", MISMILDIRIK),
        ("Combined", combined),
        ("Corpus(raw)", corpus_raw),
    ]

    all_chars = set()
    for _, data in sources:
        all_chars.update(data.keys())

    print(f"\n{'=' * 100}")
    print("COMPARISON: All Sources (percentages)")
    print(f"{'=' * 100}")

    header = f"{'Char':<6}"
    for name, _ in sources:
        header += f"{name:>14}"
    print(header)
    print("-" * 76)

    # Get average percentages for sorting (use combined for consistent ordering)
    combined_pct = calculate_percentages(combined)
    sorted_chars = sorted(combined_pct.items(), key=lambda x: x[1], reverse=True)

    # Also include chars only in corpus
    corpus_only = set(corpus_raw.keys()) - set(combined.keys())
    corpus_pct = calculate_percentages(corpus_raw)
    corpus_only_sorted = sorted(
        [(c, corpus_pct[c]) for c in corpus_only], key=lambda x: x[1], reverse=True
    )

    for char, _ in list(sorted_chars) + corpus_only_sorted:
        row = f"{char:<6}"
        for _, data in sources:
            pct = calculate_percentages(data)
            if char in pct:
                row += f"{pct[char]:>13.2f}%"
            else:
                row += f"{'—':>14}"
        print(row)


def print_tsv() -> None:
    """Print TSV output for easy import."""
    # Load corpus raw data
    script_dir = Path(__file__).parent
    corpus_raw = load_tsv(script_dir.parent / "raw_chars.tsv")

    # Combine literary sources
    combined = combine_sources(MANAS, SYNGAN_KYLYCH, MISMILDIRIK)

    sources = [
        ("Манас", MANAS),
        ("С.Кылыч", SYNGAN_KYLYCH),
        ("Мисмил.", MISMILDIRIK),
        ("Combined", combined),
        ("Corpus(raw)", corpus_raw),
    ]

    all_chars = set()
    for _, data in sources:
        all_chars.update(data.keys())

    print(f"\n{'=' * 60}")
    print("TSV OUTPUT")
    print(f"{'=' * 60}")

    # Header
    print("char\t" + "\t".join(name for name, _ in sources))

    # Sort by combined percentage
    combined_pct = calculate_percentages(combined)
    sorted_chars = sorted(combined_pct.items(), key=lambda x: x[1], reverse=True)

    # Also include chars only in corpus
    corpus_only = set(corpus_raw.keys()) - set(combined.keys())
    corpus_pct = calculate_percentages(corpus_raw)
    corpus_only_sorted = sorted(
        [(c, corpus_pct[c]) for c in corpus_only], key=lambda x: x[1], reverse=True
    )

    for char, _ in list(sorted_chars) + corpus_only_sorted:
        row = [char]
        for _, data in sources:
            pct = calculate_percentages(data).get(char, 0)
            row.append(f"{pct:.3f}")
        print("\t".join(row))


if __name__ == "__main__":
    print("Kyrgyz Character Frequency Analysis")
    print("Source: Literary works frequency data")
    print("Note: К and Қ have been merged")

    print_stats("Манас (Жусуп Мамай)", MANAS)
    print_stats("Сынган Кылыч", SYNGAN_KYLYCH)
    print_stats("Мисмилдирик", MISMILDIRIK)

    print_comparison()
    print_tsv()
