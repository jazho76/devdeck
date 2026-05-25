-- Host packages:
--
-- Fedora: python3 python3-pip
-- Debian/Ubuntu: python3 python3-pip
return {
  lsp = {
    pyright = {},
  },
  mason = {
    'black',
    'isort',
  },
  formatters_by_ft = {
    python = { 'isort', 'black' },
  },
  treesitter = {
    'python',
  },
}
