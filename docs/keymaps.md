# Devdeck keymaps

This document lists Devdeck-defined Neovim keymaps and their descriptions.
It does not enumerate every plugin's internal buffer or picker mappings.

`<leader>` is `Space`.

Mode legend:

- `n`: normal
- `v`: visual and select
- `x`: visual
- `o`: operator-pending
- `i`: insert
- `t`: terminal

## Which-key groups

| Prefix | Mode | Description |
| --- | --- | --- |
| `<leader>c` | all | Code |
| `<leader>d` | all | Debug |
| `<leader>g` | all | Git |
| `<leader>h` | normal | Git hunk |
| `<leader>h` | visual | Git hunk |
| `<leader>l` | all | LSP |
| `<leader>n` | all | Notifications |
| `<leader>s` | all | Search |
| `<leader>t` | all | Toggle |
| `<leader>T` | all | Tests |
| `<leader>w` | all | Swap |

## Core editing

| Key | Mode | Description |
| --- | --- | --- |
| `<Space>` | normal, visual | Disabled so Space can act only as the leader prefix. |
| `j` | normal | Move down by display line when no count is provided. |
| `k` | normal | Move up by display line when no count is provided. |
| `p` | visual | Paste over the selection without replacing the default register. |
| `<A-j>` | normal | Move the current line down. |
| `<A-k>` | normal | Move the current line up. |
| `<A-j>` | visual | Move the selected lines down and keep the selection. |
| `<A-k>` | visual | Move the selected lines up and keep the selection. |

## Buffers, windows, and terminal

| Key | Mode | Description |
| --- | --- | --- |
| `<A-a>` | normal | Go to the previous buffer. |
| `<A-d>` | normal | Go to the next buffer. |
| `<A-q>` | normal | Delete the current buffer. |
| `<A-Q>` | normal | Force-delete the current buffer. |
| `<A-o>` | normal | Close all buffers except the current buffer. |
| `<C-Up>` | normal | Increase the current window height. |
| `<C-Down>` | normal | Decrease the current window height. |
| `<C-Left>` | normal | Decrease the current window width. |
| `<C-Right>` | normal | Increase the current window width. |
| `<C-n>` | terminal | Leave terminal insert mode and return to normal mode. |

## Diagnostics and quickfix

| Key | Mode | Description |
| --- | --- | --- |
| `[d` | normal | Go to previous diagnostic message. |
| `]d` | normal | Go to next diagnostic message. |
| `<leader>e` | normal | Open floating diagnostic message. |
| `<leader>q` | normal | Open diagnostics list. |
| `[q` | normal | Go to previous quickfix item. |
| `]q` | normal | Go to next quickfix item. |

## File explorer and navigation

| Key | Mode | Description |
| --- | --- | --- |
| `<A-e>` | normal | Toggle the NvimTree file explorer. |
| `<C-h>` | normal | Move focus to the pane or Neovim split on the left. |
| `<C-j>` | normal | Move focus to the pane or Neovim split below. |
| `<C-k>` | normal | Move focus to the pane or Neovim split above. |
| `<C-l>` | normal | Move focus to the pane or Neovim split on the right. |
| `<C-\>` | normal | Move focus to the previously active tmux pane or Neovim split. |

## Search and Telescope

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>?` | normal | Find recently opened files. |
| `<leader><Space>` | normal | Find existing buffers. |
| `<leader>/` | normal | Fuzzily search in the current buffer. |
| `<leader>s/` | normal | Search in open files. |
| `<leader>st` | normal | Search Telescope builtins. |
| `<leader>sf` | normal | Search files. |
| `<leader>sh` | normal | Search help. |
| `<leader>sw` | normal | Search current word. |
| `<leader>sg` | normal | Search by grep. |
| `<leader>sG` | normal | Search by grep from the Git root. |
| `<leader>sd` | normal | Search diagnostics. |
| `<leader>sm` | normal | Search marks. |
| `<leader>sr` | normal | Resume the last Telescope picker. |
| `<leader>sb` | normal | Search buffers. |

## Git

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>gg` | normal | Open Fugitive Git status and make it the only window. |
| `<leader>gf` | normal | Search Git files. |
| `<leader>gc` | normal | Search Git commits. |
| `<leader>gb` | normal | Search Git branches. |
| `<leader>gs` | normal | Search Git status. |
| `<leader>gh` | normal | Search Git commit history for the current buffer. |
| `]c` | normal, visual | Jump to the next Git hunk. Falls back to Vim's diff mapping in diff mode. |
| `[c` | normal, visual | Jump to the previous Git hunk. Falls back to Vim's diff mapping in diff mode. |
| `<leader>hs` | normal | Stage the current Git hunk. |
| `<leader>hs` | visual | Stage the selected Git hunk range. |
| `<leader>hr` | normal | Reset the current Git hunk. |
| `<leader>hr` | visual | Reset the selected Git hunk range. |
| `<leader>hS` | normal | Stage the current buffer. |
| `<leader>hu` | normal | Undo staging for the current hunk. |
| `<leader>hR` | normal | Reset the current buffer. |
| `<leader>hp` | normal | Preview the current Git hunk. |
| `<leader>hb` | normal | Show blame for the current line. |
| `<leader>hd` | normal | Diff the current file against the index. |
| `<leader>hD` | normal | Diff the current file against the last commit. |
| `<leader>tb` | normal | Toggle Git blame for the current line. |
| `<leader>td` | normal | Toggle showing deleted Git lines. |
| `ih` | operator-pending, visual | Select the current Git hunk. |

## LSP

These mappings are buffer-local and appear after an LSP client attaches.

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>cr` | normal | LSP: code rename. |
| `<leader>ca` | normal | LSP: code action. |
| `<leader>lr` | normal | LSP: restart. |
| `gd` | normal | LSP: go to definition. |
| `gr` | normal | LSP: go to references. |
| `gI` | normal | LSP: go to implementation. |
| `<leader>D` | normal | LSP: type definition. |
| `<leader>ss` | normal | LSP: search document symbols. |
| `<leader>sS` | normal | LSP: search workspace symbols. |
| `K` | normal | LSP: hover documentation. |
| `gD` | normal | LSP: go to declaration. |

## Treesitter textobjects

### Selection

| Key | Mode | Description |
| --- | --- | --- |
| `aa` | visual, operator-pending | Select the outer parameter textobject. |
| `ia` | visual, operator-pending | Select the inner parameter textobject. |
| `af` | visual, operator-pending | Select the outer function textobject. |
| `if` | visual, operator-pending | Select the inner function textobject. |
| `ac` | visual, operator-pending | Select the outer class textobject. |
| `ic` | visual, operator-pending | Select the inner class textobject. |
| `ar` | visual, operator-pending | Select the outer parameter textobject. |
| `ab` | visual, operator-pending | Select the outer block textobject. |
| `ib` | visual, operator-pending | Select the inner block textobject. |
| `ak` | visual, operator-pending | Select the assignment left-hand side textobject. |
| `ik` | visual, operator-pending | Select the assignment left-hand side textobject. |
| `av` | visual, operator-pending | Select the assignment right-hand side textobject. |
| `iv` | visual, operator-pending | Select the assignment right-hand side textobject. |

### Movement

| Key | Mode | Description |
| --- | --- | --- |
| `]a` | normal, visual, operator-pending | Jump to the next start of the inner parameter textobject. |
| `]b` | normal, visual, operator-pending | Jump to the next start of the inner block textobject. |
| `]f` | normal, visual, operator-pending | Jump to the next start of the outer function textobject. |
| `]k` | normal, visual, operator-pending | Jump to the next start of the assignment left-hand side textobject. |
| `]v` | normal, visual, operator-pending | Jump to the next start of the assignment right-hand side textobject. |
| `]A` | normal, visual, operator-pending | Jump to the next end of the inner parameter textobject. |
| `]B` | normal, visual, operator-pending | Jump to the next end of the inner block textobject. |
| `]F` | normal, visual, operator-pending | Jump to the next end of the outer function textobject. |
| `]K` | normal, visual, operator-pending | Jump to the next end of the assignment left-hand side textobject. |
| `]V` | normal, visual, operator-pending | Jump to the next end of the assignment right-hand side textobject. |
| `[a` | normal, visual, operator-pending | Jump to the previous start of the inner parameter textobject. |
| `[b` | normal, visual, operator-pending | Jump to the previous start of the inner block textobject. |
| `[f` | normal, visual, operator-pending | Jump to the previous start of the outer function textobject. |
| `[k` | normal, visual, operator-pending | Jump to the previous start of the assignment left-hand side textobject. |
| `[v` | normal, visual, operator-pending | Jump to the previous start of the assignment right-hand side textobject. |
| `[A` | normal, visual, operator-pending | Jump to the previous end of the inner parameter textobject. |
| `[B` | normal, visual, operator-pending | Jump to the previous end of the inner block textobject. |
| `[F` | normal, visual, operator-pending | Jump to the previous end of the outer function textobject. |
| `[K` | normal, visual, operator-pending | Jump to the previous end of the assignment left-hand side textobject. |
| `[V` | normal, visual, operator-pending | Jump to the previous end of the assignment right-hand side textobject. |

### Swap and repeatable movement

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>wa` | normal | Swap the current parameter with the next parameter. |
| `<leader>wf` | normal | Swap the current function with the next function. |
| `<leader>wA` | normal | Swap the current parameter with the previous parameter. |
| `<leader>wF` | normal | Swap the current function with the previous function. |
| `;` | normal, visual, operator-pending | Repeat the last Treesitter textobject move forward. |
| `,` | normal, visual, operator-pending | Repeat the last Treesitter textobject move backward. |
| `f` | normal, visual, operator-pending | Use repeatable forward character search. |
| `F` | normal, visual, operator-pending | Use repeatable backward character search. |
| `t` | normal, visual, operator-pending | Use repeatable forward till-character search. |
| `T` | normal, visual, operator-pending | Use repeatable backward till-character search. |

## Debugging

| Key | Mode | Description |
| --- | --- | --- |
| `<F5>` | normal | Continue or start the current debug session. |
| `<F6>` | normal | Step over in the current debug session. |
| `<F7>` | normal | Step into in the current debug session. |
| `<F8>` | normal | Step out in the current debug session. |
| `<leader>db` | normal | Toggle a breakpoint. |
| `<leader>dB` | normal | Toggle a conditional breakpoint. |
| `<leader>du` | normal | Toggle the DAP UI. |

## Tests

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>Tr` | normal | Run the nearest test. |

## Notifications

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>nh` | normal | Show notification history. |
| `<leader>nd` | normal | Dismiss all notifications. |

## Formatting and AI

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>tf` | normal | Toggle autoformat. |
| `<C-S>` | insert | Accept the current Copilot inline suggestion. |

## Sessions

| Key | Mode | Description |
| --- | --- | --- |
| `<leader>o` | normal | Toggle Obsession session tracking. |
