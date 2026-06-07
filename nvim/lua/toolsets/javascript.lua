--[[devdeck
{
  "requires": [
    { "bin": ["node"] },
    { "bin": ["npm"] }
  ]
}
]]
return {
  lsp = {
    ts_ls = {},
  },
  mason = {
    'prettier',
  },
  formatters_by_ft = {
    javascript = { 'prettier' },
  },
  treesitter = {
    'javascript',
  },
}
