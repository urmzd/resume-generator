# Resume Generator — CLI Demo
# Requires: vhs (https://github.com/charmbracelet/vhs)

Output assets/demo-cli.gif

Set Shell "zsh"
Set FontFamily "MonaspiceNe Nerd Font"
Set FontSize 16
Set Width 1200
Set Height 700
Set Padding 24
Set Theme { "name": "Cyberdream", "black": "#16181a", "red": "#ff6e5e", "green": "#5eff6c", "yellow": "#f1ff5e", "blue": "#5ea1ff", "magenta": "#bd5eff", "cyan": "#5ef1ff", "white": "#ffffff", "brightBlack": "#3c4048", "brightRed": "#ff6e5e", "brightGreen": "#5eff6c", "brightYellow": "#f1ff5e", "brightBlue": "#5ea1ff", "brightMagenta": "#bd5eff", "brightCyan": "#5ef1ff", "brightWhite": "#ffffff", "background": "#16181a", "foreground": "#ffffff", "selection": "#3c4048", "cursor": "#ffffff" }
Set WindowBar "Rings"
Set TypingSpeed 50ms
Set CursorBlink false

# ── Setup (hidden from recording) ──────────────────
Hide
Type "go build -o resume-generator . 2>/dev/null"
Enter
Sleep 10s
Type "export PS1='❯ '"
Enter
Sleep 500ms
Type "splash() { clear; print ''; print '  \e[90m╭──────────────────────────────────────────────╮\e[0m'; print '  \e[90m│\e[0m                                              \e[90m│\e[0m'; print '  \e[90m│\e[0m          \e[35m✦\e[0m  \e[1;36mResume Generator\e[0m  \e[35m✦\e[0m              \e[90m│\e[0m'; print '  \e[90m│\e[0m                                              \e[90m│\e[0m'; print '  \e[90m│\e[0m         \e[34mgithub.com/urmzd\e[0m \e[90m·\e[0m \e[34murmzd.com\e[0m         \e[90m│\e[0m'; print '  \e[90m│\e[0m                                              \e[90m│\e[0m'; print '  \e[90m╰──────────────────────────────────────────────╯\e[0m'; print ''; }"
Enter
Sleep 500ms
Type "splash"
Enter
Sleep 500ms
Show

# ── Splash pause ──────────────────────────────────
Sleep 3s

# ── CLI overview ──────────────────────────────────
Type "./resume-generator --help"
Enter
Sleep 4s

# ── List available templates ──────────────────────
Ctrl+L
Sleep 500ms
Type "./resume-generator templates list"
Enter
Sleep 3s

# ── Generate a resume ─────────────────────────────
Ctrl+L
Sleep 500ms
Type "./resume-generator run -i assets/example_resumes/software_engineer.yml -t modern-html -o /tmp/demo-output"
Enter
Sleep 8s

# ── Show the output structure ─────────────────────
Type "tree /tmp/demo-output"
Enter
Sleep 3s

# ── Done ──────────────────────────────────────────
Sleep 2s
