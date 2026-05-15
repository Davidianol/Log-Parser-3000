package sections

import (
	"encoding/csv"
	"fmt"
	"strings"
)

func detectDelimiter(header string) (rune, error) {
	if strings.Contains(header, ",") {
		return ',', nil
	}
	// Насколько я знаю, ';' не может быть по стандарту в ibdiagnet2.db_csv, поэтому проверки на нее здесь нет
	return 0, fmt.Errorf("cannot detect delimiter in header: %q", header)
}

func parseCSVLine(line string, delimiter rune) ([]string, error) {
	r := csv.NewReader(strings.NewReader(line))
	r.Comma = delimiter
	r.LazyQuotes = false
	r.TrimLeadingSpace = true

	record, err := r.Read()
	if err != nil {
		return nil, err
	}
	return record, nil

}

func validateHeader(actual string, expected string, delimiter rune) error {
	got, err := parseCSVLine(actual, delimiter)
	if err != nil {
		return fmt.Errorf("parse actual header: %w", err)
	}
	want, err := parseCSVLine(expected, delimiter)
	if err != nil {
		return fmt.Errorf("parse expected header: %w", err)
	}
	if len(got) != len(want) {
		return fmt.Errorf("invalid header column count: got=%d want=%d", len(got), len(want))
	}
	for i := range want {
		if strings.TrimSpace(got[i]) != strings.TrimSpace(want[i]) {
			return fmt.Errorf("invalid header at col=%d: got=%q want=%q", i, got[i], want[i])
		}
	}
	return nil
}
