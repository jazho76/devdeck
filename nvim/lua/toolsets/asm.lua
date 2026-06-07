--[[devdeck
{
  "requires": [
    { "bin": ["cargo"] },
    { "bin": ["rustc"] },
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
