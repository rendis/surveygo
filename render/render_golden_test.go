package render

import (
	"encoding/json"
	"os"
	"path/filepath"
	"strings"
	"testing"

	surveygo "github.com/rendis/surveygo/v2"
)

const testdataDir = "testdata"
const outputDir = "testdata/output"

// TestRenderTreeOnly processes each survey-only JSON (no _answers in name)
// and generates tree outputs (tree.json + tree.html) into per-survey folders.
func TestRenderTreeOnly(t *testing.T) {
	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("reading testdata dir: %v", err)
	}

	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".json") || strings.Contains(name, "_answers") {
			continue
		}

		base := strings.TrimSuffix(name, ".json")
		t.Run(base, func(t *testing.T) {
			survey := parseSurvey(t, filepath.Join(testdataDir, name))

			treeResult, err := DefinitionTree(survey)
			if err != nil {
				t.Fatalf("DefinitionTree: %v", err)
			}

			outDir := filepath.Join(outputDir, base)
			mustMkdir(t, outDir)

			treeJSON, _ := json.MarshalIndent(treeResult.JSON, "", "  ")
			writeFile(t, filepath.Join(outDir, "tree.json"), treeJSON)
			writeFile(t, filepath.Join(outDir, "tree.html"), treeResult.HTML)

			t.Logf("tree output → %s", outDir)
		})
	}
}

// TestRenderExamples processes each answers JSON paired with its survey,
// generating all render outputs into per-example folders.
func TestRenderExamples(t *testing.T) {
	entries, err := os.ReadDir(testdataDir)
	if err != nil {
		t.Fatalf("reading testdata dir: %v", err)
	}

	for _, e := range entries {
		name := e.Name()
		if !strings.HasSuffix(name, ".json") || !strings.Contains(name, "_answers") {
			continue
		}

		base := strings.TrimSuffix(name, ".json")
		t.Run(base, func(t *testing.T) {
			surveyFile := deriveSurveyFile(name)
			survey := parseSurvey(t, filepath.Join(testdataDir, surveyFile))
			answers := loadAnswers(t, filepath.Join(testdataDir, name))

			outDir := filepath.Join(outputDir, base)
			mustMkdir(t, outDir)

			// Tree outputs
			treeResult, err := DefinitionTree(survey)
			if err != nil {
				t.Fatalf("DefinitionTree: %v", err)
			}
			treeJSON, _ := json.MarshalIndent(treeResult.JSON, "", "  ")
			writeFile(t, filepath.Join(outDir, "tree.json"), treeJSON)
			writeFile(t, filepath.Join(outDir, "tree.html"), treeResult.HTML)

			// All answer-based outputs
			result, err := AnswersTo(survey, answers, OutputOptions{
				CSV: true, JSON: true, HTML: true, TipTap: true,
			})
			if err != nil {
				t.Fatalf("AnswersTo: %v", err)
			}

			writeFile(t, filepath.Join(outDir, "answers.csv"), result.CSV)

			if result.JSON != nil {
				cardJSON, _ := json.MarshalIndent(result.JSON, "", "  ")
				writeFile(t, filepath.Join(outDir, "card.json"), cardJSON)
			}
			if result.HTML != nil {
				writeFile(t, filepath.Join(outDir, "card.html"), result.HTML.HTML)
				writeFile(t, filepath.Join(outDir, "card.css"), result.HTML.CSS)
			}
			if result.TipTap != nil {
				tiptapJSON, _ := json.MarshalIndent(result.TipTap, "", "  ")
				writeFile(t, filepath.Join(outDir, "card_tiptap.json"), tiptapJSON)
			}

			t.Logf("render output → %s", outDir)
		})
	}
}

// deriveSurveyFile extracts the survey filename from an answers filename.
// E.g. "sample_with_choice_answers_repeat.json" → "sample_with_choice.json"
func deriveSurveyFile(answersFile string) string {
	name := strings.TrimSuffix(answersFile, ".json")
	idx := strings.Index(name, "_answers")
	if idx < 0 {
		return answersFile
	}
	return name[:idx] + ".json"
}

func parseSurvey(t *testing.T, path string) *surveygo.Survey {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading survey %s: %v", path, err)
	}
	survey, err := surveygo.ParseFromBytes(data)
	if err != nil {
		t.Skipf("skipping %s: ParseFromBytes failed: %v", path, err)
	}
	return survey
}

func loadAnswers(t *testing.T, path string) surveygo.Answers {
	t.Helper()
	data, err := os.ReadFile(path)
	if err != nil {
		t.Fatalf("reading answers %s: %v", path, err)
	}
	var answers surveygo.Answers
	if err := json.Unmarshal(data, &answers); err != nil {
		t.Skipf("skipping %s: unmarshal answers failed: %v", path, err)
	}
	return answers
}

func mustMkdir(t *testing.T, dir string) {
	t.Helper()
	if err := os.MkdirAll(dir, 0o755); err != nil {
		t.Fatalf("mkdir %s: %v", dir, err)
	}
}

func writeFile(t *testing.T, path string, data []byte) {
	t.Helper()
	if err := os.WriteFile(path, data, 0o644); err != nil {
		t.Fatalf("writing %s: %v", path, err)
	}
}
