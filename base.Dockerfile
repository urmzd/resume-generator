# Use Debian bookworm-slim as the base image
FROM debian:bookworm-slim

# Update the package lists and install necessary packages
# - libfontconfig: Font configuration and customization library
# - wget: Utility to retrieve files from the web
# - perl: Programming language, used for the TeX Live installation script
RUN apt-get -y update && apt-get -y install libfontconfig wget perl

# Set the working directory to /tmp
WORKDIR /tmp

# Download the TeX Live installer
RUN wget https://mirror.ctan.org/systems/texlive/tlnet/install-tl-unx.tar.gz

# Unpack the TeX Live installer
RUN zcat < install-tl-unx.tar.gz | tar xf -

# Install TeX Live
# The installation is non-interactive and uses the 'small' scheme to reduce size
# TODO: Update year when the next millennium comes install-tl-2 => installt-tl-3
RUN cd install-tl-2* && perl install-tl --no-interaction --scheme=small

# Update PATH to include the TeX Live binaries
ENV PATH="/usr/local/texlive/2024/bin/x86_64-linux:${PATH}"
ENV PATH="/usr/local/texlive/2024/bin/aarch64-linux:${PATH}"

# Install additional TeX Live packages
# - enumitem: Control layout of itemize, enumerate, description
# - titlesec: Select alternative section titles
RUN tlmgr install enumitem titlesec

# Set the entrypoint to bash to allow interactive use of the container
# When running the container, it will start a bash shell
ENTRYPOINT [ "/bin/bash" ]
