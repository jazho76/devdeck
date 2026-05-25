-- Host packages:
--
-- Fedora: cargo rust
-- Debian/Ubuntu: cargo rustc
return {
  lsp = {
    rust_analyzer = {},
  },
  treesitter = {
    'rust',
  },
}
