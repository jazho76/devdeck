# DevDeck

An opinionated Neovim + tmux development environment, managed by a single CLI.

DevDeck installs a self-contained Neovim (plugins, LSP, formatters, debuggers, treesitter) and a preconfigured tmux, then keeps them updated and out of your way. Language support is modular: enable only the toolsets you need.

## Requirements

- A [Nerd Font](https://www.nerdfonts.com/), set as your terminal font, for glyphs to render correctly.

The installer checks for these system dependencies and refuses to proceed if a required one is missing:

| Dependency      | Required | Notes                                    |
| --------------- | -------- | ---------------------------------------- |
| git             | yes      |                                          |
| tmux            | yes      | 3.5a or newer                            |
| gcc / cc        | yes      | builds native plugins                    |
| ripgrep (`rg`)  | yes      | project-wide text search                 |
| tree-sitter-cli | yes      |                                          |
| make            | no       | compiles telescope's native fuzzy finder |
| fd              | no       | faster file finding                      |
| wl-clipboard    | no       | system clipboard on Wayland              |

## Install

```bash
curl -fsSL https://raw.githubusercontent.com/jazho76/devdeck/main/scripts/get.sh | sh
```

This downloads the `devdeck` CLI to `~/.local/bin` and runs `devdeck install`, which sets up Neovim and tmux. Make sure `~/.local/bin` is on your `PATH`:

```bash
export PATH="$HOME/.local/bin:$PATH"
```

During install you will be prompted to pick the language toolsets you want. You can change this at any time later with `devdeck toolsets`.

## Commands

```text
devdeck install      Install the environment (Neovim + tmux)
devdeck toolsets     Choose which language toolsets are enabled
devdeck workspace    Save and restore tmux sessions
devdeck update       Update the config and plugins
devdeck upgrade      Update the devdeck binary itself
devdeck doctor       Check that the installation is healthy
devdeck uninstall    Remove the environment
devdeck version      Print the version
```

Run any command with `--help` for its full set of flags.

## Toolsets

A toolset is a per-language bundle: LSP servers, formatters, debug adapters, and treesitter parsers, wired together so a language just works when enabled. Disabled toolsets cost nothing.

Open the interactive picker:

```bash
devdeck toolsets
```

Or set them non-interactively:

```bash
devdeck toolsets --toolsets go,rust,python   # enable exactly these
devdeck toolsets --all-toolsets              # enable everything
devdeck toolsets --no-toolsets               # disable everything
```

Available toolsets:

```text
angular   asm      c_cpp     config   csharp
eslint    go       graphql   javascript   javascript-debug
markdown  python   react     rust     typescript   web-markup
```

Each toolset may need extra tools on your `PATH` (compilers, language runtimes). The picker shows these requirements per toolset, and warns about anything missing after you select.

## Workspaces

A workspace is a saved tmux session: its windows, panes, layouts, and working directories. Restore one to pick up exactly where you left off.

The intended way to use workspaces is from inside tmux, without leaving your session. The prefix is `C-a`:

- `prefix C-s` opens a popup to save the current session.
- `prefix Tab` opens a popup to restore a saved one.

The same operations are available from the CLI, mostly for scripting or managing saved workspaces outside tmux:

```bash
devdeck workspace save [name]     # save the current session
devdeck workspace restore <name>  # replace the tmux server with a saved session
devdeck workspace list            # list saved workspaces, most recent first
devdeck workspace rename <old> <new>
devdeck workspace delete <name>   # add --force to skip the prompt
```

## Updating

DevDeck has two update paths:

```bash
devdeck update    # pull the latest config and plugin updates
devdeck upgrade   # update the devdeck binary itself
```

Run `devdeck doctor` anytime to verify the installation is intact and dependencies are satisfied.

## Uninstall

```bash
devdeck uninstall
```

This removes the DevDeck Neovim, tmux config, and managed source. Your Neovim plugin data and caches are left in place. To remove those as well:

```bash
devdeck uninstall --purge
```
