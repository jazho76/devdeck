# DevDeck

An opinionated neovim + tmux development environment.

## Prerequisites

Install the required packages first.

<details>
<summary>Fedora</summary>

```bash
sudo dnf install -y \
  git \
  curl \
  tar \
  gzip \
  tmux \
  gcc \
  ripgrep \
  nodejs \
  npm \
  tree-sitter-cli
```

</details>

<details>
<summary>Ubuntu / Debian</summary>

```bash
sudo apt install -y \
  git \
  curl \
  tar \
  gzip \
  tmux \
  gcc \
  ripgrep \
  nodejs \
  npm
```

`tree-sitter-cli` is not packaged in apt; install it with `npm install -g tree-sitter-cli` or `cargo install tree-sitter-cli`.

</details>

Optional packages enhance the setup.

<details>
<summary>Fedora</summary>

```bash
sudo dnf install -y \
  make \
  fd-find \
  wl-clipboard
```

</details>

<details>
<summary>Ubuntu / Debian</summary>

```bash
sudo apt install -y \
  make \
  fd-find \
  wl-clipboard
```

On Ubuntu/Debian, `fd-find` installs the binary as `fdfind`.

</details>

- `make` - faster Telescope fuzzy sorting (builds telescope-fzf-native).
- `fd-find` - faster file finding in Telescope.
- `wl-clipboard` - system clipboard sync on Wayland.

A [Nerd Font](https://www.nerdfonts.com/) is required for terminal glyphs to render correctly.

Toolset packages are optional. Install the packages for the toolsets you enable:

| Toolset            | Fedora packages                               | Ubuntu / Debian packages            |
| ------------------ | --------------------------------------------- | ----------------------------------- |
| `angular`          | `nodejs npm`                                  | `nodejs npm`                        |
| `asm`              | `cargo rust openssl-devel pkgconf-pkg-config` | `cargo rustc libssl-dev pkg-config` |
| `c_cpp`            | `gcc gcc-c++ make`                            | `gcc g++ make`                      |
| `config`           | `nodejs npm`                                  | `nodejs npm`                        |
| `csharp`           | `dotnet-sdk-10.0`                             | `dotnet-sdk-10.0`                   |
| `eslint`           | `nodejs npm`                                  | `nodejs npm`                        |
| `go`               | `go`                                          | `golang-go`                         |
| `graphql`          | `nodejs npm`                                  | `nodejs npm`                        |
| `javascript`       | `nodejs npm`                                  | `nodejs npm`                        |
| `javascript-debug` | `nodejs npm`                                  | `nodejs npm`                        |
| `markdown`         | `nodejs npm`                                  | `nodejs npm`                        |
| `python`           | `python3 python3-pip`                         | `python3 python3-pip`               |
| `react`            | `nodejs npm`                                  | `nodejs npm`                        |
| `rust`             | `cargo rust`                                  | `cargo rustc`                       |
| `typescript`       | `nodejs npm`                                  | `nodejs npm`                        |
| `web-markup`       | `nodejs npm`                                  | `nodejs npm`                        |

## Install

```bash
git clone https://github.com/jazho76/devdeck.git ~/.devdeck
cd ~/.devdeck
./scripts/install.sh
```

## Uninstall

```bash
./scripts/uninstall.sh
```
