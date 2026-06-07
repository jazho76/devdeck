--[[devdeck
{
  "requires": [
    { "bin": ["node"] },
    { "bin": ["npm"] }
  ]
}
]]
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
