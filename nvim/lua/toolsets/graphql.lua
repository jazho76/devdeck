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
    graphql = { 'prettier' },
  },
  treesitter = {
    'graphql',
  },
}
