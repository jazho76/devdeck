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
    javascriptreact = { 'prettier' },
    typescriptreact = { 'prettier' },
  },
  treesitter = {
    'tsx',
  },
}
