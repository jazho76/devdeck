-- Host packages
--
-- Fedora: cargo rust openssl-devel pkgconf-pkg-config
-- Debian/Ubuntu: cargo rustc libssl-dev pkg-config
return {
  lsp = {
    asm_lsp = {},
  },
  treesitter = {
    'asm',
  },
}
