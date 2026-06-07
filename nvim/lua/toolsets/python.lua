--[[devdeck
{
  "requires": [
    { "bin": ["python3"] },
    { "bin": ["pip3", "pip"] }
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
