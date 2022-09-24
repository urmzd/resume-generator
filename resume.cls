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

\NewDocumentCommand\project{mmmO{}}{%
    #1 \hfill #2 \hfill #3
    \IfNoValueTF{#4}
        {\nolinebreak}
        {\newline #4}%
}

\NewDocumentCommand\achievements{O{11}O{22}O{33}}{%
    \begin{itemize}
        \item #1
        \item #2
        \item #3
    \end{itemize}%
}
