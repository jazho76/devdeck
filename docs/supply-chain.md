# Supply-chain notes

Devdeck does not have a project-local npm application surface. There is no
`package.json`, lockfile, npm script, or Lazy plugin hook that runs `npm install`
from this repository.

The exposed install paths are still worth separating because they behave very
differently.

## Mason-managed npm tools

Mason is the main npm entry point. The default toolsets install several
npm-distributed tools through Mason, while `lazy-lock.json` only pins Neovim
plugins and does not lock Mason registry packages.

Default npm-backed Mason tools:

- `prettier`
  - Used by `markdown`, `config`, `javascript`, `typescript`, `react`,
    `web-markup`, and `graphql`.
  - Published as `pkg:npm/prettier`.
  - No runtime npm dependency fanout in the current package metadata.
- `typescript-language-server`
  - Used by `javascript`, `typescript`, and `react` through `ts_ls`.
  - Published as `pkg:npm/typescript-language-server`.
  - No runtime npm dependency fanout in the current package metadata.
- `vscode-langservers-extracted`
  - Installed through Mason packages such as `eslint-lsp`, `html-lsp`, and
    `css-lsp`.
  - This is the largest default npm dependency fanout observed in the audit.
- `pyright`
  - Used by the Python toolset.
  - Published as `pkg:npm/pyright`.
  - Its metadata includes optional `fsevents`; Linux installs normally skip it.

Optional npm-backed tools:

- `angular` installs `@angular/language-server` when that optional toolset is
  enabled.
- `javascript-debug` uses Mason's `js-debug-adapter` package. Devdeck consumes
  a Microsoft `vscode-js-debug` release asset, not a registry `npm install`, but
  it is still JavaScript tooling from an external release channel.

Mason auto-install is configured in `nvim/lua/plugins/nvim-lspconfig.lua` with
`mason-lspconfig` and `mason-tool-installer`. That is convenient for a personal
workbench, but it means opening Neovim can fetch enabled toolchain packages.

## Lazy plugin hooks

Lazy plugin installation is separate from Mason.

Current notable plugin hooks:

- `nvim-treesitter` runs `require('nvim-treesitter').install(parsers)` and
  builds parser libraries.
- `telescope-fzf-native.nvim` runs `make`.
- `copilot.vim` expects a Node runtime for Copilot behavior, but Devdeck does
  not run npm to install it.
- `friendly-snippets` provides VS Code snippet data through a Git plugin
  dependency, not npm.

No current Lazy plugin hook runs `npm install`, `npx`, `yarn`, or `pnpm`.

## Tree-sitter CLI on Fedora

Fedora's `tree-sitter-cli` package is not an npm package and does not install a
private copy of Node.

Observed on Fedora 42:

- Package: `tree-sitter-cli-0.25.10-1.fc42`
- Source RPM: `rust-tree-sitter-cli-0.25.10-1.fc42.src.rpm`
- Installed payload: `/usr/bin/tree-sitter`, documentation, and licenses
- Dynamic dependencies: glibc, libgcc, libm, and the ELF loader
- No RPM dependency on `node`, `nodejs`, `npm`, or `libnode`

The CLI can still execute a JavaScript runtime when asked to generate a parser
from `grammar.js`:

```text
tree-sitter generate --js-runtime <EXECUTABLE>
```

The default runtime is `node`, or the value of `TREE_SITTER_JS_RUNTIME`. That is
an external process lookup, not an embedded Node copy.

For Devdeck's current parser set, the tracked `nvim-treesitter` metadata did not
mark any configured parser with `generate = true` or `generate_from_json = true`.
The normal path is therefore to clone parser repositories and build generated C
sources, not to evaluate grammar JavaScript with Node.

That still has native-code supply-chain exposure: parser repositories are cloned,
C/C++ parser code is compiled, and Neovim loads the resulting parser libraries in
process.

## Practical posture

Keep the early guardrails simple:

- Treat Mason as the npm boundary and Tree-sitter as the native parser boundary.
- Keep `nodejs npm` limited to toolsets that actually need JavaScript tooling.
- Avoid adding Lazy plugin hooks that run `npm install` unless the tradeoff is
  explicit.
- Review `nvim/lua/toolsets/*.lua` before enabling new default toolsets.
- If stricter reproducibility becomes necessary later, pin or mirror Mason
  packages separately from `lazy-lock.json`; the Lazy lockfile is not enough.
