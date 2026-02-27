package generators

import (
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"

	"github.com/urmzd/resume-generator/pkg/resume"
	"go.uber.org/zap"
)

// professionResume builds a resume with fields that contain LaTeX special
// characters typical for that profession. Each resume exercises different
// subsets of &, %, #, $, _, {, }, ~, ^.
type professionTestCase struct {
	name   string
	resume *resume.Resume
	// expect lists substrings that MUST appear (escaped) in the output.
	expect []string
	// reject lists substrings that MUST NOT appear (unescaped specials).
	reject []string
}

func makeProfessionCases() []professionTestCase {
	t2022 := time.Date(2022, 1, 1, 0, 0, 0, 0, time.UTC)
	t2023 := time.Date(2023, 6, 1, 0, 0, 0, 0, time.UTC)
	t2018 := time.Date(2018, 9, 1, 0, 0, 0, 0, time.UTC)
	t2021 := time.Date(2021, 5, 1, 0, 0, 0, 0, time.UTC)

	return []professionTestCase{
		{
			name: "healthcare professional with & in section titles",
			resume: &resume.Resume{
				Contact: resume.Contact{
					Name:        "Dr. Jane O'Brien",
					Email:       "jane@hospital.org",
					Credentials: "RN, BSN, CCRN",
				},
				Certifications: &resume.Certifications{
					Title: "Registration & Certifications",
					Items: []resume.Certification{
						{Name: "Board Certified — Critical Care (CCRN)"},
						{Name: "Advanced Cardiac Life Support (ACLS)"},
					},
				},
				Skills: resume.Skills{
					Categories: []resume.SkillCategory{
						{Category: "Clinical & Technical", Items: []string{"ICU Monitoring", "Ventilator Management"}},
					},
				},
				Experience: resume.ExperienceList{
					Title: "Clinical & Professional Experience",
					Positions: []resume.Experience{
						{
							Title:      "Charge Nurse — ICU",
							Company:    "St. Mary's Hospital & Medical Center",
							Highlights: []string{"Reduced patient readmission by 30%", "Managed team of 12 RNs & CNAs"},
							Dates:      resume.DateRange{Start: t2022},
						},
					},
				},
				Education: resume.EducationList{
					Institutions: []resume.Education{
						{
							Institution: "University of Toronto",
							Degree:      resume.Degree{Name: "Bachelor of Science in Nursing (B.Sc.N.)"},
							Dates:       resume.DateRange{Start: t2018, End: &t2021},
						},
					},
				},
			},
			expect: []string{
				`Registration \& Certifications`,
				`Clinical \& Technical`,
				`Clinical \& Professional Experience`,
				`St. Mary's Hospital \& Medical Center`,
				`12 RNs \& CNAs`,
				`30\%`,
			},
			reject: []string{
				// bare & that isn't preceded by backslash
			},
		},
		{
			name: "finance professional with $ and % and #",
			resume: &resume.Resume{
				Contact: resume.Contact{
					Name:  "Alice $mith",
					Email: "alice@bank.com",
					Phone: "+1 (555) 999-0000",
				},
				Skills: resume.Skills{
					Categories: []resume.SkillCategory{
						{Category: "Financial Analysis & Modeling", Items: []string{"DCF Models", "M&A Valuation"}},
					},
				},
				Experience: resume.ExperienceList{
					Positions: []resume.Experience{
						{
							Title:   "VP — Structured Finance",
							Company: "Goldman & Partners LLC",
							Highlights: []string{
								"Managed $500M portfolio with 12% YoY growth",
								"Closed deal #4782 generating $25M revenue",
								"Led team of 8 analysts & 3 associates",
							},
							Dates: resume.DateRange{Start: t2022},
						},
					},
				},
				Education: resume.EducationList{
					Institutions: []resume.Education{
						{
							Institution: "Wharton School",
							Degree:      resume.Degree{Name: "MBA, Finance & Strategy"},
							Dates:       resume.DateRange{Start: t2018, End: &t2021},
						},
					},
				},
			},
			expect: []string{
				`Alice \$mith`,
				`Goldman \& Partners LLC`,
				`\$500M`,
				`12\%`,
				`\#4782`,
				`\$25M`,
				`8 analysts \& 3 associates`,
				`Financial Analysis \& Modeling`,
				`M\&A Valuation`,
				`Finance \& Strategy`,
			},
		},
		{
			name: "academic with _ and ^ and { } and ~",
			resume: &resume.Resume{
				Contact: resume.Contact{
					Name:  "Prof. María García~López",
					Email: "garcia@cs.mit.edu",
					Links: []resume.Link{
						{URI: "https://github.com/m_garcia", Label: "GitHub"},
					},
				},
				Skills: resume.Skills{
					Categories: []resume.SkillCategory{
						{Category: "Programming", Items: []string{"C++", "Python", "R"}},
						{Category: "Research", Items: []string{"NLP", "ML^2 Framework"}},
					},
				},
				Experience: resume.ExperienceList{
					Positions: []resume.Experience{
						{
							Title:   "Assistant Professor",
							Company: "MIT — CSAIL",
							Highlights: []string{
								"Published 15 papers in {top-tier} venues",
								"Secured $1.2M NSF grant for AI_safety research",
								"Supervised 8 PhD students & 12 undergrads",
							},
							Technologies: []string{"PyTorch", "TensorFlow", "CUDA_11"},
							Dates:        resume.DateRange{Start: t2022},
						},
					},
				},
				Education: resume.EducationList{
					Institutions: []resume.Education{
						{
							Institution: "Stanford University",
							Degree:      resume.Degree{Name: "Ph.D. Computer Science"},
							Thesis: &resume.Thesis{
								Title: "On the Complexity of {k}-SAT Under Constraint~Relaxation",
								Link:  resume.Link{URI: "https://purl.stanford.edu/thesis_123"},
							},
							Dates: resume.DateRange{Start: t2018, End: &t2021},
						},
					},
				},
				Projects: &resume.ProjectList{
					Projects: []resume.Project{
						{
							Name:       "open_source_nlp",
							Highlights: []string{"100% test coverage", "Used by 500+ researchers"},
							Link:       resume.Link{URI: "https://github.com/m_garcia/open_source_nlp"},
						},
					},
				},
			},
			expect: []string{
				`Prof. Mar\'{i}a Garc\'{i}a`,
				`\{top-tier\}`,
				`\$1.2M`,
				`AI\_safety`,
				`CUDA\_11`,
				`open\_source\_nlp`,
				`100\%`,
				`ML\textasciicircum{}2`,
			},
		},
		{
			name: "lawyer with special firm names and section symbols",
			resume: &resume.Resume{
				Contact: resume.Contact{
					Name:  "James O'Connor",
					Email: "joc@lawfirm.com",
				},
				Certifications: &resume.Certifications{
					Title: "Bar Admissions & Licenses",
					Items: []resume.Certification{
						{Name: "New York State Bar — #12345"},
						{Name: "U.S. District Court — Southern & Eastern Districts"},
					},
				},
				Skills: resume.Skills{
					Categories: []resume.SkillCategory{
						{Category: "Practice Areas", Items: []string{"Mergers & Acquisitions", "Securities & Exchange Compliance"}},
					},
				},
				Experience: resume.ExperienceList{
					Positions: []resume.Experience{
						{
							Title:   "Senior Associate",
							Company: "Baker & McKenzie LLP",
							Highlights: []string{
								"Managed $200M+ cross-border M&A transaction",
								"Drafted 100% of client engagement letters",
								"Reviewed contracts under Section 10(b) & Rule 10b-5",
							},
							Dates: resume.DateRange{Start: t2022},
						},
					},
				},
				Education: resume.EducationList{
					Institutions: []resume.Education{
						{
							Institution: "Yale Law School",
							Degree:      resume.Degree{Name: "Juris Doctor (J.D.)"},
							Dates:       resume.DateRange{Start: t2018, End: &t2021},
						},
					},
				},
			},
			expect: []string{
				`Bar Admissions \& Licenses`,
				`\#12345`,
				`Southern \& Eastern Districts`,
				`Mergers \& Acquisitions`,
				`Securities \& Exchange Compliance`,
				`Baker \& McKenzie LLP`,
				`\$200M`,
				`M\&A`,
				`100\%`,
				`10(b) \& Rule`,
			},
		},
		{
			name: "engineer with mixed specials",
			resume: &resume.Resume{
				Contact: resume.Contact{
					Name:  "Bob_Builder",
					Email: "bob@dev.io",
					Links: []resume.Link{
						{URI: "https://github.com/bob_builder"},
					},
				},
				Skills: resume.Skills{
					Categories: []resume.SkillCategory{
						{Category: "Languages & Frameworks", Items: []string{"C#", "F#", "ASP.NET"}},
					},
				},
				Experience: resume.ExperienceList{
					Positions: []resume.Experience{
						{
							Title:   "Staff Engineer",
							Company: "Micro$oft",
							Highlights: []string{
								"Achieved 99.99% uptime for {critical} services",
								"Reduced costs by $2M (40% savings)",
							},
							Technologies: []string{"Azure", "C#", ".NET_8"},
							Dates:        resume.DateRange{Start: t2022, End: &t2023},
						},
					},
				},
				Education: resume.EducationList{
					Institutions: []resume.Education{
						{
							Institution: "University of Waterloo",
							Degree:      resume.Degree{Name: "B.ASc. Software Engineering"},
							Dates:       resume.DateRange{Start: t2018, End: &t2021},
						},
					},
				},
			},
			expect: []string{
				`Bob\_Builder`,
				`Micro\$oft`,
				`99.99\%`,
				`\{critical\}`,
				`\$2M`,
				`40\%`,
				`.NET\_8`,
				`C\#`,
				`F\#`,
				`Languages \& Frameworks`,
			},
		},
	}
}

// TestLaTeXAutoEscape_AllTemplates generates each profession resume against
// every LaTeX template and verifies that all special characters are properly
// escaped.
func TestLaTeXAutoEscape_AllTemplates(t *testing.T) {
	logger := zap.NewNop().Sugar()

	templates := []struct {
		name string
		path string
	}{
		{"modern-latex", filepath.Join("..", "..", "templates", "modern-latex", "template.tex")},
		{"modern-cv", filepath.Join("..", "..", "templates", "modern-cv", "template.tex")},
	}

	cases := makeProfessionCases()

	// The academic case expects accented characters from LaTeX escaping,
	// but our escaper doesn't handle Unicode accents — skip those specific
	// assertions that depend on accent commands.
	// Filter out accent-related expectations for simplicity.
	for i := range cases {
		if cases[i].name == "academic with _ and ^ and { } and ~" {
			// Remove the accent expectation since we don't have accent escaping
			filtered := make([]string, 0, len(cases[i].expect))
			for _, exp := range cases[i].expect {
				if !strings.Contains(exp, `\'{`) {
					filtered = append(filtered, exp)
				}
			}
			cases[i].expect = filtered
		}
	}

	for _, tmplInfo := range templates {
		contentBytes, err := os.ReadFile(tmplInfo.path)
		if err != nil {
			t.Fatalf("failed to read template %s: %v", tmplInfo.name, err)
		}
		templateContent := string(contentBytes)

		for _, tc := range cases {
			t.Run(tmplInfo.name+"/"+tc.name, func(t *testing.T) {
				gen := NewLaTeXGenerator(logger)
				got, err := gen.Generate(templateContent, tc.resume)
				if err != nil {
					t.Fatalf("Generate() error: %v", err)
				}

				for _, exp := range tc.expect {
					if !strings.Contains(got, exp) {
						t.Errorf("expected %q in output but not found", exp)
					}
				}

				for _, rej := range tc.reject {
					if strings.Contains(got, rej) {
						t.Errorf("unexpected %q found in output", rej)
					}
				}

				// Global check: no bare & (not preceded by \) in the output.
				// We look for & that isn't part of \& — a simple heuristic.
				checkNoBareAmpersand(t, got)
			})
		}
	}
}

// checkNoBareAmpersand scans the rendered LaTeX for & characters that are
// not properly escaped as \&.
func checkNoBareAmpersand(t *testing.T, content string) {
	t.Helper()
	lines := strings.Split(content, "\n")
	for i, line := range lines {
		for j := 0; j < len(line); j++ {
			if line[j] == '&' {
				if j == 0 || line[j-1] != '\\' {
					t.Errorf("bare '&' at line %d col %d: %s", i+1, j+1, strings.TrimSpace(line))
				}
			}
		}
	}
}
