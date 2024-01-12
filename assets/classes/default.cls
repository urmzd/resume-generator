\NeedsTeXFormat{LaTeX2e}

\ProvidesClass{default}
\LoadClass{article}

\RequirePackage{fix-cm}
\RequirePackage{amsmath}
\RequirePackage{enumitem}
\RequirePackage{xparse}

\newcommand{\name}[1]{%
    \begin{center}
        \Huge{#1}
    \end{center}%
}

\newcommand{\beforeTextLine}{%
    \vspace{0.125in}%
}

\newcommand{\afterTextLine}{%
    \vspace{0.0625in}%
}

\NewDocumentCommand{\contact}{mmmo}{%
    \hfill #1 \hfill \textbar \hfill #2 \hfill \textbar \hfill #3 \hfill \IfNoValueTF{#4}{}{\textbar \hfill #4 \hfill}
}

\NewDocumentCommand\project{mmm}{%
    \makebox[0pt][l]{\textbf{#1}} \hfill \makebox[0pt][c]{\textbf{#2}} \hfill \makebox[0pt][r]{\textbf{#3}}%
}

\NewDocumentCommand\subproject{m}{%
    \newline
    \textbf{#1}
}