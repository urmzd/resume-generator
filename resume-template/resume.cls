\NeedsTeXFormat{LaTeX2e}

\ProvidesClass{resume}
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

\newcommand{\contact}[4]{%
    \hfill #1 \hfill \textbar \hfill #2 \hfill \textbar \hfill #3 \hfill \textbar \hfill #4 \hfill
}

\NewDocumentCommand\project{mmm}{%
    \makebox[0pt][l]{\textbf{#1}} \hfill \makebox[0pt][c]{\textbf{#2}} \hfill \makebox[0pt][r]{\textbf{#3}}%
}

\NewDocumentCommand\subproject{m}{%
    \newline
    \textbf{#1}
}

\NewDocumentCommand\achievements{O{Did X as measured by Y resulting in Z}ooo}{%
    \begin{itemize}
        \IfNoValueTF{#1}{}{\item #1}
        \IfNoValueTF{#2}{}{\item #2}
        \IfNoValueTF{#3}{}{\item #3}
        \IfNoValueTF{#4}{}{\item #4}
    \end{itemize}%
}
