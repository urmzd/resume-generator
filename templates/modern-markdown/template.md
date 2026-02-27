{{define "section-summary"}}
{{- if .Summary}}

## Summary

{{.Summary}}
{{end}}
{{- end}}

{{define "section-certifications"}}
{{- if .Certifications}}{{if .Certifications.Items}}

## {{default "Certifications" .Certifications.Title}}

{{range .Certifications.Items}}- **{{.Name}}**{{if .Issuer}} — {{.Issuer}}{{end}}{{if .Notes}} ({{.Notes}}){{end}}
{{end}}
{{- end}}{{end}}
{{- end}}

{{define "section-education"}}
{{- if .Education.Institutions}}

## {{default "Education" .Education.Title}}

{{range .Education.Institutions}}### {{.Institution}}{{if .Degree.Name}} — {{.Degree.Name}}{{end}}

{{fmtDateRange .Dates}}{{if .Location}} | {{fmtLocation .Location}}{{end}}{{if .GPA}}{{$gpa := formatGPA .GPA}}{{if $gpa}} | GPA: {{$gpa}}{{end}}{{end}}
{{$descriptions := filterEmpty .Degree.Descriptions}}{{if $descriptions}}
{{range $descriptions}}- {{.}}
{{end}}{{end}}
{{- if .Thesis}}
- **Thesis:** {{.Thesis.Title}}{{if .Thesis.Link.URI}} — [{{if .Thesis.Link.Label}}{{.Thesis.Link.Label}}{{else}}Link{{end}}]({{.Thesis.Link.URI}}){{end}}
{{range .Thesis.Highlights}}- {{.}}
{{end}}{{end}}
{{end}}
{{- end}}
{{- end}}

{{define "section-skills"}}
{{- if .Skills.Categories}}

## {{default "Skills" .Skills.Title}}

{{range .Skills.Categories}}- **{{.Category}}:** {{formatList .Items}}
{{end}}
{{- end}}
{{- end}}

{{define "section-experience"}}
{{- if .Experience.Positions}}

## {{default "Experience" .Experience.Title}}

{{range .Experience.Positions}}### {{.Title}}

**{{.Company}}** | {{fmtDateRange .Dates}}{{if .Location}} | {{fmtLocation .Location}}{{end}}
{{if .Technologies}}
*{{formatList .Technologies}}*
{{end}}
{{- $high := filterEmpty .Highlights}}{{if $high}}
{{range $high}}- {{.}}
{{end}}{{end}}
{{end}}
{{- end}}
{{- end}}

{{define "section-projects"}}
{{- if .Projects}}{{if .Projects.Projects}}

## {{default "Projects" .Projects.Title}}

{{range .Projects.Projects}}### {{.Name}}{{if .Link.URI}} — [{{if .Link.Label}}{{.Link.Label}}{{else}}Link{{end}}]({{.Link.URI}}){{end}}

{{- if .Dates}}
{{fmtOptDateRange .Dates}}
{{end}}
{{- if .Technologies}}
*{{formatList .Technologies}}*
{{end}}
{{- $high := filterEmpty .Highlights}}{{if $high}}
{{range $high}}- {{.}}
{{end}}{{end}}
{{end}}
{{- end}}{{end}}
{{- end}}

{{define "section-languages"}}
{{- if .Languages}}{{if .Languages.Languages}}

## {{default "Languages" .Languages.Title}}

{{range .Languages.Languages}}- **{{.Name}}**{{if .Proficiency}} — {{.Proficiency}}{{end}}
{{end}}
{{- end}}{{end}}
{{- end}}

# {{.Contact.Name}}

{{- $sep := false -}}
{{- if .Contact.Email -}}
{{if $sep}} | {{end}}[{{.Contact.Email}}](mailto:{{.Contact.Email}})
{{- $sep = true -}}
{{- end -}}
{{- if .Contact.Phone -}}
{{if $sep}} | {{end}}{{.Contact.Phone}}
{{- $sep = true -}}
{{- end -}}
{{- if .Contact.Credentials -}}
{{if $sep}} | {{end}}{{.Contact.Credentials}}
{{- $sep = true -}}
{{- end -}}
{{- if .Contact.Location -}}
{{- $loc := fmtLocation .Contact.Location -}}
{{- if $loc -}}
{{if $sep}} | {{end}}{{$loc}}
{{- $sep = true -}}
{{- end -}}
{{- end -}}
{{- range .Contact.Links -}}
{{- if .URI -}}
{{if $sep}} | {{end}}[{{if .Label}}{{.Label}}{{else}}{{extractDisplayURL .URI}}{{end}}]({{.URI}})
{{- $sep = true -}}
{{- end -}}
{{- end}}

---
{{if and .Layout .Layout.Sections}}
{{- $root := .}}
{{- range .Layout.Sections}}
{{- if eq . "summary"}}{{template "section-summary" $root}}
{{- else if eq . "certifications"}}{{template "section-certifications" $root}}
{{- else if eq . "education"}}{{template "section-education" $root}}
{{- else if eq . "skills"}}{{template "section-skills" $root}}
{{- else if eq . "experience"}}{{template "section-experience" $root}}
{{- else if eq . "projects"}}{{template "section-projects" $root}}
{{- else if eq . "languages"}}{{template "section-languages" $root}}
{{- end}}
{{- end}}
{{- else}}
{{- template "section-summary" .}}
{{- template "section-skills" .}}
{{- template "section-experience" .}}
{{- template "section-projects" .}}
{{- template "section-education" .}}
{{- template "section-certifications" .}}
{{- template "section-languages" .}}
{{- end}}
{{- if and .Layout .Layout.References }}

---

*References available upon request*
{{- end }}
