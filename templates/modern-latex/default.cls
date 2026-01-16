\NeedsTeXFormat{LaTeX2e}
\ProvidesClass{default}[2025/01/01 Modern Resume Class]
\LoadClass[11pt]{article}

% ============================================================================
% PACKAGE DEPENDENCIES
% ============================================================================
\RequirePackage{fix-cm}           % Fix computer modern font sizes
\RequirePackage{geometry}         % Page layout
\RequirePackage{titlesec}         % Section formatting
\RequirePackage{enumitem}         % List formatting
\RequirePackage{xparse}           % Advanced command definitions
\RequirePackage{array}            % Better tables
\RequirePackage{tabularx}         % Auto-width tables
\RequirePackage{hyperref}         % Hyperlinks
\RequirePackage{xcolor}           % Colors
\RequirePackage{microtype}        % Typography improvements
\RequirePackage{parskip}          % Paragraph spacing

% ============================================================================
% PAGE GEOMETRY
% ============================================================================
\geometry{
    a4paper,
    left=0.5in,
    right=0.5in,
    top=0.5in,
    bottom=0.5in
}

% ============================================================================
% TYPOGRAPHY
% ============================================================================
\RequirePackage{helvet}
\renewcommand{\familydefault}{\sfdefault}

% Disable page numbers
\pagestyle{empty}

% Paragraph formatting
\setlength{\parindent}{0pt}
\setlength{\parskip}{0pt}

% Microtype for better typography (expansion disabled for xetex compatibility)
\microtypesetup{
    protrusion=true,
    expansion=false
}

% Font setup
% User request: "Restore LaTeX look" -> Use Latin Modern
\RequirePackage{lmodern}
% User request: "Restore LaTeX look" -> Use Latin Modern
\RequirePackage{lmodern}
% \renewcommand{\familydefault}{\sfdefault} % REMOVED: Default to Serif (Roman) for classical look
\renewcommand{\familydefault}{\rmdefault} % Explicitly set to Roman (Serif)


% ============================================================================
% COLORS
% ============================================================================
\definecolor{linkcolor}{RGB}{0,102,204}
\definecolor{rulecolor}{RGB}{0,0,0}

% Hyperref setup
\hypersetup{
    colorlinks=true,
    linkcolor=linkcolor,
    urlcolor=linkcolor,
    allcolors=linkcolor
}

% ============================================================================
% SECTION FORMATTING
% ============================================================================
\titleformat{\section}
    {\normalfont\large\bfseries}          % Format
    {}                                      % Label
    {0pt}                                   % Sep
    {\MakeUppercase}                        % Before (uppercase titles)
    [\vspace{-0.5ex}\titlerule]            % After

\titlespacing*{\section}{0pt}{1.5ex plus 0.5ex minus 0.2ex}{0.8ex plus 0.2ex}

% ============================================================================
% LIST FORMATTING
% ============================================================================
\setlist[itemize]{
    leftmargin=1.5em,
    labelwidth=1em,
    labelsep=0.5em,
    itemsep=0.2ex,
    topsep=2pt,
    parsep=0pt,
    partopsep=0pt,
    label=\textbullet
}

\setlist[description]{
    leftmargin=0pt,
    itemsep=0pt,
    topsep=3pt,
    parsep=0pt,
    labelwidth=!,
    font=\normalfont\bfseries
}

% ============================================================================
% CUSTOM LENGTHS
% ============================================================================
\newlength{\resumeitemsep}
\setlength{\resumeitemsep}{0.8ex plus 0.2ex minus 0.1ex}

\newlength{\headersep}
\setlength{\headersep}{1.5ex plus 0.3ex minus 0.2ex}

% ============================================================================
% HEADER COMMANDS
% ============================================================================

% Name command
\newcommand{\resumename}[1]{%
    \begin{center}
        {\Huge\bfseries #1}%
    \end{center}%
}

% Title/Tagline command
\newcommand{\resumetitle}[1]{%
    \begin{center}
        {\large\itshape #1}%
    \end{center}%
    \vspace{1.5ex}
}

% Contact information with evenly distributed spacing
% User requirement: "space before the start of the first link and space after the last link"
% and "separated by pipes".
\NewDocumentCommand{\resumecontact}{m}{%
    \vspace{1.0ex} % Space between name and links
    \begin{center}
        {%
            % We will rely on template.tex to provide separators, 
            % but we need to ensure the container allows spacing.
            \small#1%
        }%
    \end{center}%
    \vspace{\headersep}
}

% Separator for contact info (visible pipe with spacing)
% User req: "pipes between the link and ultimately having space before the start of the first link"
% We will redefine sep to be just the pipe with spacing.
\newcommand{\sep}{\hspace{1em}|\hspace{1em}}

% Summary block
\newcommand{\resumesummary}[1]{%
    \vspace{-0.5ex}
    {\small\itshape #1}%
    \par\vspace{\headersep}
}

% ============================================================================
% EXPERIENCE/PROJECT ENTRY COMMANDS
% ============================================================================

% Three-column header (title, organization, date)
% User req: Title, Company, Date bolded. Company centered.
\NewDocumentCommand{\resumeentry}{m m m}{%
    \vspace{\resumeitemsep}%
    \noindent%
    \begin{tabularx}{\textwidth}{@{} >{\raggedright\arraybackslash}X >{\centering\arraybackslash}c >{\raggedleft\arraybackslash}X @{}}
        \textbf{#1} & \textbf{#2} & \textbf{#3} \\
    \end{tabularx}%
    \par\vspace{-0.5ex}
}

% Three-column with location (title, organization, dates, location)
% User req: Location hidden by default, but if this is used we should bold elements too.
% However, plan says I will likely remove usage of this or make it default to Title/Company/Date.
% I'll keep it but align with bolding requirements in case it IS used.
\NewDocumentCommand{\resumeentryloc}{m m m m}{%
    \vspace{\resumeitemsep}%
    \noindent%
    \begin{tabularx}{\textwidth}{@{} >{\raggedright\arraybackslash}X >{\centering\arraybackslash}c >{\raggedleft\arraybackslash}X @{}}
        \textbf{#1} & \textbf{#2} & \textbf{#3} \\
                    &              & \textbf{#4} \\
    \end{tabularx}%
    \par\vspace{-0.5ex}
}

% Education entry
% User req: Degree below. School, degree, date, location bolded.
\NewDocumentCommand{\resumeeducation}{m m m m}{%
    \vspace{\resumeitemsep}%
    \noindent%
    \begin{tabularx}{\textwidth}{@{}X r@{}}
        \textbf{#1} & \textbf{#3} \\
        \textbf{#2} & \textbf{#4} \\
    \end{tabularx}%
    \par\vspace{0.2ex} % Slight space before description/items
}

% Subentry (for degrees, etc.)
\newcommand{\resumesubentry}[1]{%
    \vspace{0.3ex}%
    \noindent\textbf{#1}%
    \par\vspace{-0.3ex}
}

% ============================================================================
% CONVENIENCE MACROS
% ============================================================================

% Separator for contact info (visible pipe with spacing)
% \newcommand{\sep}{\hspace{1em}{\textbf{|}\hspace{1em}}}

% Email link
\newcommand{\email}[1]{\href{mailto:#1}{#1}}

% Phone link
\newcommand{\phone}[1]{\href{tel:#1}{#1}}

% Location formatter
\newcommand{\location}[1]{\textit{#1}}
