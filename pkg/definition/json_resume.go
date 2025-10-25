package definition

// JSONResume is the top-level structure for a JSON Resume.
// It follows the schema defined at https://jsonresume.org/schema/
type JSONResume struct {
	Basics       Basics          `json:"basics"`
	Work         []Work          `json:"work"`
	Volunteer    []Volunteer     `json:"volunteer"`
	Education    []JSONEducation `json:"education"`
	Awards       []Award         `json:"awards"`
	Publications []Publication   `json:"publications"`
	Skills       []Skill         `json:"skills"`
	Languages    []Language      `json:"languages"`
	Interests    []Interest      `json:"interests"`
	References   []Reference     `json:"references"`
	Projects     []JSONProject   `json:"projects"`
}

// Basics contains fundamental information about the person.
// This includes contact details, website, and a summary.
type Basics struct {
	Name     string    `json:"name"`
	Label    string    `json:"label"`
	Image    string    `json:"image"`
	Email    string    `json:"email"`
	Phone    string    `json:"phone"`
	URL      string    `json:"url"`
	Summary  string    `json:"summary"`
	Location Location  `json:"location"`
	Profiles []Profile `json:"profiles"`
}

// Profile represents a social media profile.
// It includes the network name, username, and URL.
type Profile struct {
	Network  string `json:"network"`
	Username string `json:"username"`
	URL      string `json:"url"`
}

// Work details a professional position.
// It includes the company, role, dates, and summary of responsibilities.
type Work struct {
	Name       string   `json:"name"`
	Position   string   `json:"position"`
	URL        string   `json:"url"`
	StartDate  string   `json:"startDate"`
	EndDate    string   `json:"endDate"`
	Summary    string   `json:"summary"`
	Highlights []string `json:"highlights"`
}

// Volunteer describes a volunteer role.
// It includes the organization, position, dates, and a summary of the work.
type Volunteer struct {
	Organization string   `json:"organization"`
	Position     string   `json:"position"`
	URL          string   `json:"url"`
	StartDate    string   `json:"startDate"`
	EndDate      string   `json:"endDate"`
	Summary      string   `json:"summary"`
	Highlights   []string `json:"highlights"`
}

// Award represents an award or honor received.
type Award struct {
	Title   string `json:"title"`
	Date    string `json:"date"`
	Awarder string `json:"awarder"`
	Summary string `json:"summary"`
}

// Publication represents a published work.
type Publication struct {
	Name        string `json:"name"`
	Publisher   string `json:"publisher"`
	ReleaseDate string `json:"releaseDate"`
	URL         string `json:"url"`
	Summary     string `json:"summary"`
}

// Language proficiency.
type Language struct {
	Language string `json:"language"`
	Fluency  string `json:"fluency"`
}

// Interest in a particular topic.
type Interest struct {
	Name     string   `json:"name"`
	Keywords []string `json:"keywords"`
}

// Reference from a professional contact.
type Reference struct {
	Name      string `json:"name"`
	Reference string `json:"reference"`
}

// Skill represents a skill category and associated keywords.
type Skill struct {
	Name     string   `json:"name"`
	Level    string   `json:"level"`
	Keywords []string `json:"keywords"`
}

// JSONEducation entry for JSON Resume format (different from main Education)
type JSONEducation struct {
	Institution string   `json:"institution"`
	Area        string   `json:"area"`
	StudyType   string   `json:"studyType"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	GPA         string   `json:"gpa"`
	Courses     []string `json:"courses"`
}

// JSONProject entry for JSON Resume format (different from main Project)
type JSONProject struct {
	Name        string   `json:"name"`
	Description string   `json:"description"`
	Highlights  []string `json:"highlights"`
	Keywords    []string `json:"keywords"`
	StartDate   string   `json:"startDate"`
	EndDate     string   `json:"endDate"`
	URL         string   `json:"url"`
	Roles       []string `json:"roles"`
	Entity      string   `json:"entity"`
	Type        string   `json:"type"`
}

// Note: Location type is now defined in definition.go to avoid duplication
