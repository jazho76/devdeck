FROM fedora:44

RUN dnf install -y \
      sudo git curl tar gzip tmux gcc ripgrep nodejs npm tree-sitter-cli \
 && dnf clean all

RUN useradd --create-home --shell /bin/bash dev \
 && echo 'dev ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/dev \
 && chmod 0440 /etc/sudoers.d/dev

COPY scripts/sandbox-entry /usr/local/bin/sandbox-entry
RUN chmod 0755 /usr/local/bin/sandbox-entry

USER dev
WORKDIR /home/dev

RUN mkdir -p /home/dev/.local/bin
ENV PATH="/home/dev/.local/bin:${PATH}"

CMD ["sandbox-entry"]
