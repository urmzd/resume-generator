# Resume Generator — CLI Demo
# Requires: vhs (brew install vhs)
# Run: vhs e2e/demo.tape

Output assets/demo-cli.gif
Set Shell "zsh"
Set FontFamily "MonaspiceNe Nerd Font"
Set FontSize 13
Set Width 1200
Set Height 600
Set Padding 4
Set Theme {
  "name": "cyberdream",
  "black": "#16181a",
  "red": "#ff6e5e",
  "green": "#5eff6c",
  "yellow": "#f1ff5e",
  "blue": "#5ea1ff",
  "purple": "#bd5eff",
  "cyan": "#5ef1ff",
  "white": "#ffffff",
  "brightBlack": "#3c4048",
  "brightRed": "#ff6e5e",
  "brightGreen": "#5eff6c",
  "brightYellow": "#f1ff5e",
  "brightBlue": "#5ea1ff",
  "brightPurple": "#bd5eff",
  "brightCyan": "#5ef1ff",
  "brightWhite": "#ffffff",
  "background": "#16181a",
  "foreground": "#ffffff",
  "selectionBackground": "#3c4048",
  "cursor": "#ffffff",
  "cursorAccent": "#16181a"
}

# --- Build the CLI binary ---
Hide
Type "go build -o resume-generator . 2>/dev/null"
Enter
Sleep 10s
Show

# --- List available templates ---
Type "echo 'List available templates'"
Enter
Sleep 500ms
Type "./resume-generator templates list"
Enter
Sleep 3s

# --- Generate a resume ---
Type "echo 'Generate a resume with the modern-html template'"
Enter
Sleep 500ms
Type "./resume-generator run -i assets/example_resumes/software_engineer.yml -t modern-html -o /tmp/demo-output"
Enter
Sleep 8s

# --- Show the output structure ---
Type "echo 'Output structure'"
Enter
Sleep 500ms
Type "tree /tmp/demo-output"
Enter
Sleep 3s

# --- Done ---
Sleep 2s
