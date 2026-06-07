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
    html = { filetypes = { 'html', 'twig', 'hbs' } },
    cssls = {},
  },
  mason = {
    'prettier',
  },
  formatters_by_ft = {
    html = { 'prettier' },
    css = { 'prettier' },
    scss = { 'prettier' },
  },
  treesitter = {
    'html',
    'css',
  },
}
