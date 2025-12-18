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
    {\normalfont\Large\bfseries\scshape}  % Format
    {}                                      % Label
    {0pt}                                   % Sep
    {}                                      % Before
    [\vspace{-0.5ex}\titlerule]            % After

\titlespacing*{\section}{0pt}{1.5ex plus 0.5ex minus 0.2ex}{0.8ex plus 0.2ex}

% ============================================================================
% LIST FORMATTING
% ============================================================================
\setlist[itemize]{
    leftmargin=*,
    itemsep=0pt,
    topsep=3pt,
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
    \vspace{0.5ex}
}

% Contact information using tabular for clean alignment
\NewDocumentCommand{\resumecontact}{m}{%
    \begin{center}
        \begin{tabular*}{\textwidth}{@{\extracolsep{\fill}}c@{}}
            #1
        \end{tabular*}%
    \end{center}%
    \vspace{\headersep}
}

% Helper for contact line with separators
\newcommand{\contactline}[1]{%
    \begin{tabular*}{\textwidth}{@{\extracolsep{\fill}}*{20}{c}@{}}
        #1
    \end{tabular*}
}

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
\NewDocumentCommand{\resumeentry}{m m m}{%
    \vspace{\resumeitemsep}%
    \noindent%
    \begin{tabularx}{\textwidth}{@{}X r@{}}
        \textbf{#1} & \textbf{#3} \\
        \textit{#2} & \\
    \end{tabularx}%
    \par\vspace{-0.5ex}
}

% Three-column with location (title, organization, dates, location)
\NewDocumentCommand{\resumeentryloc}{m m m m}{%
    \vspace{\resumeitemsep}%
    \noindent%
    \begin{tabularx}{\textwidth}{@{}X r@{}}
        \textbf{#1} & \textbf{#3} \\
        \textit{#2} & \textit{#4} \\
    \end{tabularx}%
    \par\vspace{-0.5ex}
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

% Separator for contact info
\newcommand{\sep}{\enspace|\enspace}

% Email link
\newcommand{\email}[1]{\href{mailto:#1}{#1}}

% Phone link
\newcommand{\phone}[1]{\href{tel:#1}{#1}}

% Location formatter
\newcommand{\location}[1]{\textit{#1}}
