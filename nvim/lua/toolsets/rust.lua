--[[devdeck
{
  "requires": [
    { "bin": ["cargo"] },
    { "bin": ["rustc"] }
  ]
}
]]
return {
  lsp = {
    rust_analyzer = {},
  },
  treesitter = {
    'rust',
  },
}
