--[[devdeck
{
  "requires": [
    { "bin": ["python3"] },
    { "bin": ["node"] },
    { "bin": ["npm"] }
  ]
}
]]
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
