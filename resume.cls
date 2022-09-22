\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{resume}
\LoadClass{article}

\RequirePackage{fix-cm}
\RequirePackage{amsmath}


\newcommand{\name}[1]{%
    \begin{center}
        \Huge{#1}
    \end{center}%
}

\newcommand{\contact}[4]{%
    \hfill #1 \hfill \textbar \hfill #2 \hfill \textbar \hfill #3 \hfill \textbar \hfill #4 \hfill
    \vspace{1mm}
    \hrule%
}


