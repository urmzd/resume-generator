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
    left=0.4in,
    right=0.4in,
    top=0.2in,
    bottom=0.3in
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

% Line spacing (balanced)
\linespread{1.05}

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
    {\normalfont\small\bfseries}          % Format (small since already bold+uppercase+underlined)
    {}                                      % Label
    {0pt}                                   % Sep
    {\MakeUppercase}                        % Before (uppercase titles)
    [\vspace{-0.5ex}\titlerule]            % After

\titlespacing*{\section}{0pt}{8pt plus 1pt minus 1pt}{3pt plus 1pt}

% ============================================================================
% LIST FORMATTING
% ============================================================================
\setlist[itemize]{
    leftmargin=0.25in,
    labelwidth=1em,
    labelsep=0.3em,
    itemsep=2pt,
    topsep=3pt,
    parsep=0pt,
    partopsep=0pt,
    label=\textbullet
}

\setlist[description]{
    leftmargin=0pt,
    itemsep=2pt,
    topsep=2pt,
    parsep=0pt,
    labelwidth=!,
    font=\normalfont\bfseries
}

% ============================================================================
% CUSTOM LENGTHS
% ============================================================================
\newlength{\resumeitemsep}
\setlength{\resumeitemsep}{6pt plus 1pt minus 1pt}

\newlength{\headersep}
\setlength{\headersep}{0.5ex plus 0.1ex minus 0.1ex}

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
    \vspace{2pt} % Space between name and links
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
\newcommand{\sep}{\hspace{0.5em}|\hspace{0.5em}}

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
    \par\vspace{1pt}
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
    \par\vspace{1pt}
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
    \par\vspace{1pt}
}

% Subentry (for degrees, etc.)
\newcommand{\resumesubentry}[1]{%
    \vspace{0.3ex}%
    \noindent\textbf{#1}%
    \par\vspace{-0.3ex}
}

% ============================================================================
% DENSITY COMMANDS
% ============================================================================
\newcommand{\setdensitycompact}{%
    \newgeometry{left=0.3in,right=0.3in,top=0.15in,bottom=0.2in}%
    \setlength{\resumeitemsep}{3pt plus 1pt minus 1pt}%
    \titlespacing*{\section}{0pt}{5pt plus 1pt minus 1pt}{2pt plus 1pt}%
    \setlist[itemize]{itemsep=1pt,topsep=2pt}%
    \linespread{1.0}\selectfont%
}

\newcommand{\setdensityrelaxed}{%
    \newgeometry{left=0.6in,right=0.6in,top=0.4in,bottom=0.5in}%
    \setlength{\resumeitemsep}{10pt plus 2pt minus 2pt}%
    \titlespacing*{\section}{0pt}{12pt plus 2pt minus 2pt}{5pt plus 1pt}%
    \setlist[itemize]{itemsep=3pt,topsep=5pt}%
    \linespread{1.15}\selectfont%
}

% Typography switching
\newcommand{\settypographymodern}{%
    \RequirePackage{helvet}%
    \renewcommand{\familydefault}{\sfdefault}%
}

\newcommand{\settypographyelegant}{%
    \RequirePackage{charter}%
    \renewcommand{\familydefault}{\rmdefault}%
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
