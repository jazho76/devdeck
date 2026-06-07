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
    json = { 'prettier' },
    yaml = { 'prettier' },
  },
  treesitter = {
    'json',
    'yaml',
  },
}
