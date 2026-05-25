-- Host packages:
--
-- Fedora: nodejs npm
-- Debian/Ubuntu: nodejs npm
return {
  mason = {
    'prettier',
  },
  formatters_by_ft = {
    json = { 'prettier' },
    yaml = { 'prettier' },
  },
  treesitter = {
    'json',
    'yaml',
  },
}
