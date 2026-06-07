--[[devdeck
{
  "requires": [
    { "bin": ["cargo"] },
    { "bin": ["rustc"] },
    { "bin": ["c++", "g++"] },
    { "bin": ["pkg-config", "pkgconf"] },
    { "pc": "openssl" }
  ]
}
]]
return {
  lsp = {
    asm_lsp = {},
  },
  treesitter = {
    'asm',
  },
}
