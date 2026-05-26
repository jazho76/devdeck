# Keymap Audit

Audit of all keymaps defined in this configuration, with overlaps, gaps, and
documentation issues identified.

- Date: 2026-05-26
- Target: Neovim 0.12.x (per `nvim/README.md`)
- Leader: `Space` (localleader also `Space`)
- which-key: v3 `wk.add` spec. It registers groups only. Per-key labels are
  inherited from each mapping's `desc`, so any mapping without a `desc` is
  invisible in which-key.

Toolset modules (`lua/toolsets/*.lua`) contribute no keymaps. They only feed
LSP servers, Mason packages, formatters, DAP adapters/configs, and treesitter
parsers into the registry. All keymaps live in `lua/config/keymaps.lua` and the
plugin files.

## 1. which-key group registry

Defined in `lua/plugins/which-key.lua`.

| Prefix             | Label           | Members                            | Notes                    |
| ------------------ | --------------- | ---------------------------------- | ------------------------ |
| `<leader>c`        | [C]ode          | `cr`, `ca`                         | LSP rename and action    |
| `<leader>d`        | [D]ebug         | `db`, `dB`, `du`                   |                          |
| `<leader>g`        | [G]it           | `gf`, `gc`, `gb`, `gs`, `gh`, `gg` | telescope-git + fugitive |
| `<leader>h` (n, v) | Git [H]unk      | gitsigns set                       |                          |
| `<leader>l`        | [L]SP           | `lr`                               | group of one             |
| `<leader>n`        | [N]otifications | `nh`, `nd`                         |                          |
| `<leader>s`        | [S]earch        | telescope + LSP symbols            |                          |
| `<leader>t`        | [T]oggle        | `tf`, `tb`, `td`                   |                          |
| `<leader>T`        | [T]ests         | `Tr`                               | group of one             |
| `<leader>w`        | S[w]ap          | `wa`, `wf`, `wA`, `wF`             | all members lack `desc`  |

## 2. Leader maps

Legend: `desc` column marks whether the mapping carries a description that
which-key can display.

### `<leader>c` [C]ode (LSP, buffer-local, `nvim-lspconfig.lua`)

| Key  | Action                         | desc |
| ---- | ------------------------------ | ---- |
| `cr` | Rename                         | yes  |
| `ca` | Code action (normal mode only) | yes  |

### `<leader>d` [D]ebug (`nvim-dap.lua`)

| Key  | Action                 | desc |
| ---- | ---------------------- | ---- |
| `db` | Toggle breakpoint      | yes  |
| `dB` | Conditional breakpoint | yes  |
| `du` | Toggle DAP UI          | yes  |

### `<leader>g` [G]it (`telescope.lua`, `vim-fugitive.lua`)

| Key  | Action                       | desc |
| ---- | ---------------------------- | ---- |
| `gf` | Git files                    | yes  |
| `gc` | Git commits                  | yes  |
| `gb` | Git branches                 | yes  |
| `gs` | Git status                   | yes  |
| `gh` | Git buffer commits (history) | yes  |
| `gg` | Fugitive (`:G \| only`)      | yes  |

### `<leader>h` Git [H]unk (gitsigns, buffer-local, n + v)

| Key  | Action                   | desc |
| ---- | ------------------------ | ---- |
| `hs` | Stage hunk               | yes  |
| `hr` | Reset hunk               | yes  |
| `hS` | Stage buffer             | yes  |
| `hu` | Undo stage hunk          | yes  |
| `hR` | Reset buffer             | yes  |
| `hp` | Preview hunk             | yes  |
| `hb` | Blame line               | yes  |
| `hd` | Diff against index       | yes  |
| `hD` | Diff against last commit | yes  |

Wording is inconsistent (for example "git stage hunk" in normal mode vs "stage
git hunk" in visual mode for the same action).

### `<leader>l` [L]SP (`nvim-lspconfig.lua`)

| Key  | Action     | desc |
| ---- | ---------- | ---- |
| `lr` | LspRestart | yes  |

### `<leader>n` [N]otifications (`vim-notify.lua`)

| Key  | Action      | desc |
| ---- | ----------- | ---- |
| `nh` | History     | yes  |
| `nd` | Dismiss all | yes  |

### `<leader>s` [S]earch (`telescope.lua`, LSP symbols in `nvim-lspconfig.lua`)

| Key  | Action                   | desc                   |
| ---- | ------------------------ | ---------------------- |
| `s/` | Grep open files          | yes                    |
| `st` | Telescope builtin picker | yes (typo: "Tlescope") |
| `sf` | Find files               | yes                    |
| `sh` | Help tags                | yes                    |
| `sw` | Grep current word        | yes                    |
| `sg` | Live grep                | yes                    |
| `sG` | Live grep in git root    | yes                    |
| `sd` | Diagnostics              | yes                    |
| `sm` | Marks                    | yes                    |
| `sr` | Resume                   | yes                    |
| `sb` | Buffers                  | yes                    |
| `ss` | LSP document symbols     | yes                    |
| `sS` | LSP workspace symbols    | yes                    |

### `<leader>t` [T]oggle (`gitsigns.lua`, `conform.lua`)

| Key  | Action              | desc |
| ---- | ------------------- | ---- |
| `tf` | Toggle autoformat   | yes  |
| `tb` | Toggle line blame   | yes  |
| `td` | Toggle show deleted | yes  |

### `<leader>T` [T]ests (`nvim-test.lua`)

| Key  | Action           | desc |
| ---- | ---------------- | ---- |
| `Tr` | Run nearest test | yes  |

### `<leader>w` S[w]ap (treesitter, `nvim-treesitter.lua`)

| Key  | Action                       | desc |
| ---- | ---------------------------- | ---- |
| `wa` | Swap parameter with next     | no   |
| `wf` | Swap function with next      | no   |
| `wA` | Swap parameter with previous | no   |
| `wF` | Swap function with previous  | no   |

### Standalone leader leafs (no group)

| Key       | Action                   | Source             | desc |
| --------- | ------------------------ | ------------------ | ---- |
| `e`       | Open floating diagnostic | keymaps.lua        | yes  |
| `q`       | Diagnostics to loclist   | keymaps.lua        | yes  |
| `?`       | Recent files             | telescope.lua      | yes  |
| `<space>` | Buffers                  | telescope.lua      | yes  |
| `/`       | Fuzzy find in buffer     | telescope.lua      | yes  |
| `D`       | LSP type definition      | nvim-lspconfig.lua | yes  |

## 3. Non-leader maps

### Editor core (`keymaps.lua`)

| Key               | Mode | Action                      | desc |
| ----------------- | ---- | --------------------------- | ---- |
| `<Space>`         | n, v | Nop (frees leader)          | no   |
| `j` / `k`         | n    | Wrap-aware down/up          | no   |
| `p`               | x    | Paste keeping yank register | no   |
| `<A-j>` / `<A-k>` | n, v | Move line/selection down/up | no   |

### Buffers (`keymaps.lua`)

| Key     | Action                  | desc |
| ------- | ----------------------- | ---- |
| `<A-a>` | Previous buffer         | no   |
| `<A-d>` | Next buffer             | no   |
| `<A-q>` | Delete buffer           | no   |
| `<A-Q>` | Force delete buffer     | no   |
| `<A-o>` | Close all other buffers | no   |

### Windows and tree

| Key                      | Action           | Source        | desc |
| ------------------------ | ---------------- | ------------- | ---- |
| `<C-Up>` / `<C-Down>`    | Resize height    | keymaps.lua   | no   |
| `<C-Left>` / `<C-Right>` | Resize width     | keymaps.lua   | no   |
| `<A-e>`                  | Toggle file tree | nvim-tree.lua | no   |

### Pane navigation (`vim-tmux-navigator.lua`)

| Key     | Action        | desc |
| ------- | ------------- | ---- |
| `<C-h>` | Pane left     | no   |
| `<C-j>` | Pane down     | no   |
| `<C-k>` | Pane up       | no   |
| `<C-l>` | Pane right    | no   |
| `<C-\>` | Previous pane | no   |

### Terminal (`keymaps.lua`)

| Key     | Mode | Action             | desc |
| ------- | ---- | ------------------ | ---- |
| `<C-n>` | t    | Exit terminal mode | no   |

### Diagnostics and quickfix (`keymaps.lua`)

| Key         | Action                   | desc |
| ----------- | ------------------------ | ---- |
| `[d` / `]d` | Previous/next diagnostic | yes  |
| `[q` / `]q` | Previous/next quickfix   | yes  |

### Completion and Copilot (insert mode)

| Key                 | Action                         | Source       |
| ------------------- | ------------------------------ | ------------ |
| `<C-n>` / `<C-p>`   | Next/prev completion item      | nvim-cmp.lua |
| `<C-d>` / `<C-f>`   | Scroll docs                    | nvim-cmp.lua |
| `<C-Space>`         | Trigger completion             | nvim-cmp.lua |
| `<CR>`              | Confirm                        | nvim-cmp.lua |
| `<Tab>` / `<S-Tab>` | Next/prev item or snippet jump | nvim-cmp.lua |
| `<C-S>`             | Accept Copilot suggestion      | copilot.lua  |

Telescope disables `<C-u>` and `<C-d>` in its insert mode (`telescope.lua`).

### LSP goto (buffer-local, `nvim-lspconfig.lua`)

| Key  | Action         | desc |
| ---- | -------------- | ---- |
| `gd` | Definition     | yes  |
| `gr` | References     | yes  |
| `gI` | Implementation | yes  |
| `gD` | Declaration    | yes  |
| `K`  | Hover          | yes  |

### Treesitter textobjects and moves (`nvim-treesitter.lua`)

None of these carry a `desc`.

- Select (modes x, o): `aa`, `ia`, `af`, `if`, `ac`, `ic`, `ar`, `ab`, `ib`,
  `ak`, `ik`, `av`, `iv`
- Move to next start: `]a`, `]b`, `]f`, `]k`, `]v`
- Move to next end: `]A`, `]B`, `]F`, `]K`, `]V`
- Move to previous start: `[a`, `[b`, `[f`, `[k`, `[v`
- Move to previous end: `[A`, `[B`, `[F`, `[K`, `[V`
- Repeatable motions: `;`, `,`, `f`, `F`, `t`, `T` (override native behavior with
  repeat-aware versions)

### DAP function keys (`nvim-dap.lua`)

| Key    | Action    | desc |
| ------ | --------- | ---- |
| `<F5>` | Continue  | no   |
| `<F6>` | Step over | no   |
| `<F7>` | Step into | no   |
| `<F8>` | Step out  | no   |

### Gitsigns navigation and text object (`gitsigns.lua`)

| Key         | Mode | Action             | desc |
| ----------- | ---- | ------------------ | ---- |
| `]c` / `[c` | n, v | Next/previous hunk | yes  |
| `ih`        | o, x | Select hunk        | yes  |

## 4. Findings

### Conflicts

1. RESOLVED (2026-05-26). `<C-k>` collision: tmux-navigator mapped `<C-k>` to
   pane-up globally while LSP `on_attach` mapped it to signature help
   buffer-locally, and the buffer-local map won, breaking upward pane navigation
   on every code buffer. Fixed by removing the LSP signature-help binding;
   `<C-k>` is now pane-up everywhere. Signature inspection remains via `K`
   (hover) and the `lsp_signature.nvim` auto-popup.
2. `gr` shadows Neovim defaults. Neovim 0.11+ ships `grn` (rename), `gra`
   (action), `grr` (references), `gri` (implementation) under the `gr` prefix.
   Mapping `gr` directly makes it a complete mapping, so the `gr` prefix now
   waits for `timeoutlen` and the built-ins become awkward to reach. The custom
   `<leader>cr`, `<leader>ca`, `gd`, and `gI` also duplicate these built-ins.

### Overlaps and redundancy

3. `<leader><space>` and `<leader>sb` both open the telescope buffer picker.
4. Three diagnostic entry points: `<leader>e` (float), `<leader>q` (loclist),
   `<leader>sd` (telescope).
5. `K` (hover) duplicates the Neovim built-in `K`. Signature help is now covered
   by the `lsp_signature.nvim` auto-popup and the built-in insert-mode `<C-s>`.

### Gaps (installed but no keymap)

6. diffview.nvim: no keymaps. Reachable only via `:DiffviewOpen`. Natural fit
   under `<leader>g`.
7. todo-comments.nvim: no jump maps (`]t` / `[t`) and no `:TodoTelescope` map.
   Effectively passive.
8. vim-dadbod-ui: reachable only via the alpha dashboard `d` button or `:DBUI`.
   No leader map.
9. vim-obsession: command-only (likely acceptable).
10. No visual-mode code action. `<leader>ca` is normal mode only. Range code
    actions are common.

### Documentation gaps (no `desc`, so invisible in which-key)

11. `<leader>w` swap maps: the group exists but its four children are unlabeled.
12. Treesitter textobjects and moves (about 40 maps): no `desc`. Mostly
    operator/visual, but normal-mode moves such as `]f` and `]a` do surface.
13. DAP `<F5>` through `<F8>`: no `desc`. Function keys do not show in which-key
    regardless, but this is inconsistent with `<leader>d*`.
14. Most non-leader editor, buffer, and window maps lack `desc`. Lower priority
    since they are not leader-driven.

### Grouping problems

15. LSP functionality is spread across five namespaces: rename and action in
    [C]ode, restart in [L]SP, symbols in [S]earch, goto in bare `g*`, and type
    definition in standalone `<leader>D`. The [L]SP group then holds only `lr`.
16. `<leader>D` (type definition) sits next to the lowercase `<leader>d` [D]ebug
    group. A case-only distinction adjacent to a group prefix is error-prone.
17. Groups of one: [L]SP (`lr`) and [T]ests (`Tr`).

### Convention inconsistencies

18. Mixed `desc` styles: telescope and LSP use bracket hints (for example
    "[S]earch [F]iles"), gitsigns uses lowercase prose ("preview git hunk"), DAP
    uses a mixed form.
19. `<leader>st` description typo: "Tlescope".
20. Normal and visual hunk descriptions differ for the same conceptual action.

## 5. Candidates to cut or consolidate

- `<leader>sb` (duplicate of `<leader><space>`).
- `<leader>D` (move away from the debug-group prefix, or fold type definition
  under [C]ode or the goto cluster).
- `gr`, `gd`, `gI`, `<leader>cr`, `<leader>ca` given the Neovim 0.11+ built-ins.
