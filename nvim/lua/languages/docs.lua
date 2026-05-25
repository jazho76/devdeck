-- Host packages:
--
-- Fedora: nodejs npm
-- Debian/Ubuntu: nodejs npm
return {
  mason = {
    'prettier',
  },
  formatters_by_ft = {
    markdown = { 'prettier' },
  },
  treesitter = {
    'markdown',
    'markdown_inline',
  },
}
