FROM fedora:44

RUN dnf install -y sudo && dnf clean all

RUN useradd --create-home --shell /bin/bash dev \
 && echo 'dev ALL=(ALL) NOPASSWD:ALL' > /etc/sudoers.d/dev \
 && chmod 0440 /etc/sudoers.d/dev

USER dev
WORKDIR /home/dev

RUN mkdir -p /home/dev/.local/bin
ENV PATH="/home/dev/.local/bin:/opt/sandbox/bin:${PATH}"

CMD ["bash", "-lc", "[ -x /opt/sandbox/bin/usage ] && /opt/sandbox/bin/usage; exec bash"]
