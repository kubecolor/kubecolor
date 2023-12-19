package tablescan

import (
	"strings"
	"testing"

	"github.com/kubecolor/kubecolor/testutil"
)

func TestScanner_emptyLines(t *testing.T) {
	const input = "\n\n\n"

	s := NewScanner(strings.NewReader(input))

	mustScanCells(t, s)
	mustScanCells(t, s)
	mustScanCells(t, s)

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %q", s.Bytes())
	}
}

func TestScanner_regularTable(t *testing.T) {
	const input = "" +
		"NAME    READY   STATUS      RESTARTS         AGE\n" +
		"pod-a   1/1     Running     21 (7d23h ago)   250d\n" +
		"pod-b   0/1     Error       0                13d\n"

	s := NewScanner(strings.NewReader(input))

	mustScanCells(t, s, "NAME", "READY", "STATUS", "RESTARTS", "AGE")
	mustScanCells(t, s, "pod-a", "1/1", "Running", "21 (7d23h ago)", "250d")
	mustScanCells(t, s, "pod-b", "0/1", "Error", "0", "13d")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %q", s.Bytes())
	}
}

func TestScanner_describeTable(t *testing.T) {
	const input = "" +
		"Name    Ready   Status      Restarts         Age\n" +
		"----    -----   ------      --------         ---\n" +
		"pod-a   1/1     Running     21 (7d23h ago)   250d\n" +
		"pod-b   0/1     Error       0                13d\n"

	s := NewScanner(strings.NewReader(input))

	mustScanCells(t, s, "Name", "Ready", "Status", "Restarts", "Age")
	mustScanCells(t, s, "----", "-----", "------", "--------", "---")
	mustScanCells(t, s, "pod-a", "1/1", "Running", "21 (7d23h ago)", "250d")
	mustScanCells(t, s, "pod-b", "0/1", "Error", "0", "13d")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %q", s.Bytes())
	}
}

func TestScanner_emptyCells(t *testing.T) {
	const input = "" +
		"NAME    READY   STATUS      RESTARTS         AGE\n" +
		"pod-a   1/1     Running                      250d\n" +
		"pod-b           Error       0\n"

	s := NewScanner(strings.NewReader(input))

	mustScanCells(t, s, "NAME", "READY", "STATUS", "RESTARTS", "AGE")
	mustScanCells(t, s, "pod-a", "1/1", "Running", "", "250d")
	mustScanCells(t, s, "pod-b", "", "Error", "0", "")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %q", s.Bytes())
	}
}

func TestScanner_multipleTables(t *testing.T) {
	const input = "" +
		"NAME    READY   STATUS      RESTARTS         AGE\n" +
		"pod-a   1/1     Running     21 (7d23h ago)   250d\n" +
		"pod-b   0/1     Error       0                13d\n" +
		"\n" +
		"NAME       READY   UP-TO-DATE   AVAILABLE   AGE\n" +
		"deploy-a   1/1     1            1           250d\n" +
		"deploy-b   0/1     0            0           13d\n"

	s := NewScanner(strings.NewReader(input))

	mustScanCells(t, s, "NAME", "READY", "STATUS", "RESTARTS", "AGE")
	mustScanCells(t, s, "pod-a", "1/1", "Running", "21 (7d23h ago)", "250d")
	mustScanCells(t, s, "pod-b", "0/1", "Error", "0", "13d")
	mustScanCells(t, s) // should detect that headers needs to be recalculated
	mustScanCells(t, s, "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE")
	mustScanCells(t, s, "deploy-a", "1/1", "1", "1", "250d")
	mustScanCells(t, s, "deploy-b", "0/1", "0", "0", "13d")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %q", s.Bytes())
	}
}

func TestScanner_multipleTightlyPackedTables(t *testing.T) {
	const input = "" +
		"NAME    READY   STATUS      RESTARTS         AGE\n" +
		"pod-a   1/1     Running     21 (7d23h ago)   250d\n" +
		"pod-b   0/1     Error       0                13d\n" +
		"NAME       READY   UP-TO-DATE   AVAILABLE   AGE\n" +
		"deploy-a   1/1     1            1           250d\n" +
		"deploy-b   0/1     0            0           13d\n"

	s := NewScanner(strings.NewReader(input))

	mustScanCells(t, s, "NAME", "READY", "STATUS", "RESTARTS", "AGE")
	mustScanCells(t, s, "pod-a", "1/1", "Running", "21 (7d23h ago)", "250d")
	mustScanCells(t, s, "pod-b", "0/1", "Error", "0", "13d")
	// on the next line it should detect that columns are not lining up
	// and should reevalutate the column indices
	mustScanCells(t, s, "NAME", "READY", "UP-TO-DATE", "AVAILABLE", "AGE")
	mustScanCells(t, s, "deploy-a", "1/1", "1", "1", "250d")
	mustScanCells(t, s, "deploy-b", "0/1", "0", "0", "13d")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %q", s.Bytes())
	}
}

func TestScanner_noHeaderAndEmptyCells(t *testing.T) {
	const input = "" +
		"pod-a           Running                      250d\n" +
		"pod-b   0/1                 0\n" +
		"pod-c   2/2     Error       0                13d\n"

	s := NewScanner(strings.NewReader(input))

	mustScanCells(t, s, "pod-a", "", "Running", "", "250d")
	mustScanCells(t, s, "pod-b", "0/1", "", "0", "")
	mustScanCells(t, s, "pod-c", "2/2", "Error", "0", "13d")

	if s.Scan() {
		t.Fatalf("Expected no more scans, but got: %q", s.Bytes())
	}
}

func mustScanCells(t *testing.T, s *Scanner, cells ...string) {
	t.Helper()
	if !s.Scan() {
		t.Fatalf("Failed scan; Expected value %q", strings.Join(cells, " | "))
	}

	var gotCells []string
	for _, c := range s.Cells() {
		gotCells = append(gotCells, c.Trimmed)
	}

	testutil.MustEqual(t, cells, gotCells)
}
