# Neovim Configuration

Standalone Neovim configuration.

## Uninstall previous version

```bash
rm -rf ~/.local/nvim
rm -f ~/.local/bin/nvim
```

## Remove previous data and cache

```bash
rm -rf ~/.local/share/nvim
rm -rf ~/.cache/nvim
rm -rf ~/.config/nvim
```

## Install Nvim

```bash
VERSION=0.12.2
curl -L https://github.com/neovim/neovim/releases/download/v$VERSION/nvim-linux-x86_64.tar.gz -o nvim-linux-x86_64.tar.gz
tar xzvf nvim-linux-x86_64.tar.gz
rm ./nvim-linux-x86_64.tar.gz
mkdir -p ~/.local ~/.local/bin ~/.config
rm -rf ~/.local/nvim
mv ./nvim-linux-x86_64 ~/.local/nvim
ln -sfn ~/.local/nvim/bin/nvim ~/.local/bin/nvim
```

Make sure `~/.local/bin` is on your `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

## Install dependencies

Install Fedora host packages:

Core:

```bash
sudo dnf install git curl tar gzip unzip gcc gcc-c++ make ripgrep fd-find wl-clipboard fontconfig tree-sitter-cli
```

Toolsets:

```bash
sudo dnf install nodejs npm python3 python3-pip cargo rust openssl-devel pkgconf-pkg-config go dotnet-sdk-10.0
```

Install a Nerd Font for icons:

```bash
NERD_FONT_VERSION=3.4.0
curl -L https://github.com/ryanoasis/nerd-fonts/releases/download/v$NERD_FONT_VERSION/DroidSansMono.zip -o DroidSansMono.zip
unzip DroidSansMono.zip -d DroidSansMono/
mkdir -p ~/.local/share/fonts/
cp ./DroidSansMono/*.otf ~/.local/share/fonts/
rm -rf ./DroidSansMono
rm ./DroidSansMono.zip
fc-cache -fv
fc-list | rg 'DroidSansM'
```

Then configure your terminal to use `DroidSansM Nerd Font`.

## Install configuration

Clone directly to the Neovim config directory:

```bash
git clone https://github.com/jazho76/nvim.git ~/.config/nvim
```

Open Neovim and let Lazy install plugins:

```bash
nvim
```
