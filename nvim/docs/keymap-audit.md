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
parsers into the registry. Keymaps live in `lua/config/keymaps.lua`, the plugin
files, and a handful of plugin defaults (Comment.nvim, nvim-surround) that ship
mappings without any local configuration.

Scope legend used throughout: **global** maps load at startup; **buffer-local**
maps are created in an LSP/gitsigns `on_attach` and exist only in matching
buffers (so they also only surface in which-key inside those buffers).

## 1. which-key group registry

Defined in `lua/plugins/which-key.lua`.

| Prefix             | Label           | Members                            | Notes                                  |
| ------------------ | --------------- | ---------------------------------- | -------------------------------------- |
| `<leader>c`        | [C]ode          | `cr`, `ca`                         | LSP rename and action (buffer-local)   |
| `<leader>d`        | [D]ebug         | `db`, `dB`, `du`                   |                                        |
| `<leader>g`        | [G]it           | `gf`, `gc`, `gb`, `gs`, `gh`, `gg` | telescope-git + fugitive               |
| `<leader>h` (n, v) | Git [H]unk      | gitsigns set (buffer-local)        |                                        |
| `<leader>l`        | [L]SP           | `lr`                               | group of one (buffer-local)            |
| `<leader>n`        | [N]otifications | `nh`, `nd`                         |                                        |
| `<leader>s`        | [S]earch        | telescope + LSP symbols            | `ss`/`sS` are buffer-local             |
| `<leader>t`        | [T]oggle        | `tf`, `tb`, `td`                   |                                        |
| `<leader>T`        | [T]ests         | `Tr`                               | group of one                           |
| `<leader>w`        | S[w]ap          | `wa`, `wf`, `wA`, `wF`             | all members lack `desc`                |

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

`hs` and `hr` exist in both normal and visual mode. Wording is inconsistent:
visual mode uses "stage git hunk" / "reset git hunk", normal mode uses "git
stage hunk" / "git reset hunk" for the same action, and `hS` / `hR` use
mid-sentence capitals ("git Stage buffer", "git Reset buffer").

### `<leader>l` [L]SP (LSP, buffer-local, `nvim-lspconfig.lua`)

| Key  | Action     | desc |
| ---- | ---------- | ---- |
| `lr` | LspRestart | yes  |

### `<leader>n` [N]otifications (`vim-notify.lua`)

| Key  | Action      | desc |
| ---- | ----------- | ---- |
| `nh` | History     | yes  |
| `nd` | Dismiss all | yes  |

### `<leader>s` [S]earch (`telescope.lua`, LSP symbols in `nvim-lspconfig.lua`)

| Key  | Action                   | desc                   | Scope        |
| ---- | ------------------------ | ---------------------- | ------------ |
| `s/` | Grep open files          | yes                    | global       |
| `st` | Telescope builtin picker | yes (typo: "Tlescope") | global       |
| `sf` | Find files               | yes                    | global       |
| `sh` | Help tags                | yes                    | global       |
| `sw` | Grep current word        | yes                    | global       |
| `sg` | Live grep                | yes                    | global       |
| `sG` | Live grep in git root    | yes                    | global       |
| `sd` | Diagnostics              | yes                    | global       |
| `sm` | Marks                    | yes                    | global       |
| `sr` | Resume                   | yes                    | global       |
| `sb` | Buffers                  | yes                    | global       |
| `ss` | LSP document symbols     | yes                    | buffer-local |
| `sS` | LSP workspace symbols    | yes                    | buffer-local |

### `<leader>t` [T]oggle (`gitsigns.lua`, `conform.lua`)

| Key  | Action              | desc | Scope                  |
| ---- | ------------------- | ---- | ---------------------- |
| `tf` | Toggle autoformat   | yes  | global (`conform`)     |
| `tb` | Toggle line blame   | yes  | buffer-local (gitsigns)|
| `td` | Toggle show deleted | yes  | buffer-local (gitsigns)|

### `<leader>T` [T]ests (`nvim-test.lua`)

| Key  | Action           | desc |
| ---- | ---------------- | ---- |
| `Tr` | Run nearest test | yes  |

Only `:TestNearest` is mapped. `nvim-test` also exposes `:TestFile`,
`:TestSuite`, `:TestLast`, and `:TestVisit`, none of which have keymaps.

### `<leader>w` S[w]ap (treesitter, `nvim-treesitter.lua`)

| Key  | Action                       | desc |
| ---- | ---------------------------- | ---- |
| `wa` | Swap parameter with next     | no   |
| `wf` | Swap function with next      | no   |
| `wA` | Swap parameter with previous | no   |
| `wF` | Swap function with previous  | no   |

### Standalone leader leafs (no group)

| Key       | Action                   | Source             | Scope        | desc |
| --------- | ------------------------ | ------------------ | ------------ | ---- |
| `e`       | Open floating diagnostic | keymaps.lua        | global       | yes  |
| `q`       | Diagnostics to loclist   | keymaps.lua        | global       | yes  |
| `?`       | Recent files             | telescope.lua      | global       | yes  |
| `<space>` | Buffers                  | telescope.lua      | global       | yes  |
| `/`       | Fuzzy find in buffer     | telescope.lua      | global       | yes  |
| `D`       | LSP type definition      | nvim-lspconfig.lua | buffer-local | yes  |

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
Note that `<C-S>` and `<C-s>` are the same keycode to Neovim; copilot's accept
therefore occupies the slot Neovim 0.11+ assigns by default to insert-mode
signature help (`vim.lsp.buf.signature_help`). See finding 4.

### LSP goto (buffer-local, `nvim-lspconfig.lua`)

| Key  | Action         | desc |
| ---- | -------------- | ---- |
| `gd` | Definition     | yes  |
| `gr` | References     | yes  |
| `gI` | Implementation | yes  |
| `gD` | Declaration    | yes  |
| `K`  | Hover          | yes  |

`gd`, `gr`, and `gI` use telescope pickers; `gD` and `K` use the native
`vim.lsp.buf.*` calls. There is no signature-help map (see finding 4).

### Comment.nvim (plugin defaults, no `desc`)

`Comment.nvim` is loaded with `opts = {}`, which enables its full default
mapping set. None carry a `desc`, so they are invisible in which-key.

| Key            | Mode | Action                       |
| -------------- | ---- | ---------------------------- |
| `gcc`          | n    | Toggle line comment          |
| `gbc`          | n    | Toggle block comment         |
| `gc{motion}`   | n    | Linewise comment operator    |
| `gb{motion}`   | n    | Blockwise comment operator   |
| `gc`           | x    | Comment selection (linewise) |
| `gb`           | x    | Comment selection (blockwise)|
| `gco`          | n    | Add comment line below       |
| `gcO`          | n    | Add comment line above       |
| `gcA`          | n    | Add comment at end of line   |

### nvim-surround (plugin defaults, no `desc`)

`nvim-surround` is loaded with `opts = {}`, enabling its default mappings. These
override the native `ys`, `ds`, `cs`, and visual `S` behaviors.

| Key            | Mode | Action                         |
| -------------- | ---- | ------------------------------ |
| `ys{motion}c`  | n    | Add surround                   |
| `yss c`        | n    | Add surround to whole line     |
| `yS{motion}c`  | n    | Add surround on new lines      |
| `ds{char}`     | n    | Delete surround                |
| `cs{char}c`    | n    | Change surround                |
| `cS{char}c`    | n    | Change surround on new lines   |
| `S{char}`      | x    | Surround selection             |
| `gS{char}`     | x    | Surround selection, new lines  |
| `<C-g>s`       | i    | Add surround                   |
| `<C-g>S`       | i    | Add surround on new lines      |

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

Note: `ak`/`ik` both map to `@assignment.lhs` (inner and outer are identical),
and `ar` duplicates `aa` (both `@parameter.outer`). Likely copy-paste leftovers.

### DAP function keys (`nvim-dap.lua`)

| Key    | Action    | desc |
| ------ | --------- | ---- |
| `<F5>` | Continue  | no   |
| `<F6>` | Step over | no   |
| `<F7>` | Step into | no   |
| `<F8>` | Step out  | no   |

### Gitsigns navigation and text object (`gitsigns.lua`, buffer-local)

| Key         | Mode | Action             | desc |
| ----------- | ---- | ------------------ | ---- |
| `]c` / `[c` | n, v | Next/previous hunk | yes  |
| `ih`        | o, x | Select hunk        | yes  |

## 4. Findings

### Conflicts

1. `gr` shadows Neovim defaults. Neovim 0.11+ ships `grn` (rename), `gra`
   (action), `grr` (references), `gri` (implementation) under the `gr` prefix.
   Mapping `gr` directly makes it a complete mapping, so the `gr` prefix now
   waits for `timeoutlen` and the built-ins become awkward to reach. The custom
   `<leader>cr`, `<leader>ca`, `gd`, and `gI` also duplicate these built-ins.
2. Plugin namespace overrides (intentional, documented here for completeness).
   `nvim-surround` overrides native `ys`, `ds`, `cs`, `cS`, and visual `S`;
   `Comment.nvim` claims the `gc*`/`gb*` operator namespace. These are expected
   plugin behaviors, not bugs, but they occupy the `g`, `c`, `d`, and `y`
   namespaces and constrain any future custom maps there.

### Overlaps and redundancy

3. `<leader><space>` and `<leader>sb` both open the telescope buffer picker.
4. `K` (hover) duplicates the Neovim built-in `K`. Signature help is now
   provided only by the `lsp_signature.nvim` auto-popup: Neovim's default
   insert-mode `<C-s>` is shadowed by copilot's `<C-S>` accept (same keycode),
   and there is no normal-mode signature-help map. Fine if the auto-popup is the
   intended source; worth a manual normal-mode map otherwise.
5. Three diagnostic entry points: `<leader>e` (float), `<leader>q` (loclist),
   `<leader>sd` (telescope).

### Gaps (installed but no keymap)

6. diffview.nvim: no keymaps. Reachable only via `:DiffviewOpen`. Natural fit
   under `<leader>g`.
7. todo-comments.nvim: no jump maps (`]t` / `[t`) and no `:TodoTelescope` /
   `:TodoQuickFix` map. Effectively passive.
8. vim-dadbod-ui: reachable only via the alpha dashboard `d` button or `:DBUI`.
   No leader map.
9. nvim-test: only `<leader>Tr` (`:TestNearest`). `:TestFile`, `:TestSuite`,
   `:TestLast`, and `:TestVisit` are unmapped, leaving the [T]ests group at one
   member.
10. conform.nvim: format-on-save plus `<leader>tf` toggle and a `:Format`
    command, but no "format now" keymap. Minor.
11. vim-obsession: command-only (likely acceptable).
12. No visual-mode code action. `<leader>ca` is normal mode only. Range code
    actions are common.

### Documentation gaps (no `desc`, so invisible in which-key)

13. `<leader>w` swap maps: the group exists but its four children are unlabeled.
14. Treesitter textobjects and moves (about 40 maps): no `desc`. Mostly
    operator/visual, but normal-mode moves such as `]f` and `]a` do surface.
15. Comment.nvim and nvim-surround default maps: no `desc`. The `g`-prefixed
    Comment maps (`gcc`, `gc`, `gb`, etc.) surface in the which-key `g` popup
    unlabeled.
16. DAP `<F5>` through `<F8>`: no `desc`. Function keys do not show in which-key
    regardless, but this is inconsistent with `<leader>d*`.
17. Most non-leader editor, buffer, and window maps lack `desc`. Lower priority
    since they are not leader-driven.

### Grouping problems

18. LSP functionality is spread across five namespaces: rename and action in
    [C]ode, restart in [L]SP, symbols in [S]earch, goto in bare `g*`, and type
    definition in standalone `<leader>D`. The [L]SP group then holds only `lr`.
19. `<leader>D` (type definition) sits next to the lowercase `<leader>d` [D]ebug
    group. A case-only distinction adjacent to a group prefix is error-prone.
20. Groups of one: [L]SP (`lr`) and [T]ests (`Tr`).

### Convention inconsistencies

21. Mixed `desc` styles: telescope and LSP use bracket hints (for example
    "[S]earch [F]iles"), gitsigns uses lowercase prose ("preview git hunk"), DAP
    uses a mixed form.
22. `<leader>st` description typo: "Tlescope".
23. Normal and visual hunk descriptions differ for the same conceptual action,
    and `hS`/`hR` use mid-sentence capitals.

### Likely copy-paste errors

24. Treesitter `ak`/`ik` both map to `@assignment.lhs` (no distinct inner/outer),
    and `ar` duplicates `aa` (`@parameter.outer`). The intended targets were
    probably `@assignment.outer`/`@assignment.inner` and `@return`/regex.

## 5. Candidates to cut or consolidate

These are proposals pending confirmation of actual usage — the call on what is
unused is the user's.

- `<leader>sb` (duplicate of `<leader><space>`).
- `<leader>D` (move away from the debug-group prefix, or fold type definition
  under [C]ode or the goto cluster).
- `gr`, `gd`, `gI`, `<leader>cr`, `<leader>ca` given the Neovim 0.11+ built-ins
  (decision pending; depends on whether telescope pickers are preferred over the
  native list handlers).
- `K` hover (duplicate of the built-in; only worth keeping if the explicit map
  is doing something the default does not).
