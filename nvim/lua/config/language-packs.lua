local enabled_languages = require('config.languages')

local enabled_language_set = {}
for _, language in ipairs(enabled_languages) do
  enabled_language_set[language] = true
end

local function append_unique(target, values)
  for _, value in ipairs(values or {}) do
    if not vim.tbl_contains(target, value) then
      table.insert(target, value)
    end
  end
end

local function merge_tables(target, values)
  for key, value in pairs(values or {}) do
    target[key] = value
  end
end

local packs = {
  lsp = {},
  mason = {},
  formatters_by_ft = {},
  dap_adapters = {},
  dap_configurations = {},
  treesitter = {
    'lua',
    'vimdoc',
    'vim',
    'bash',
  },
}

for _, language in ipairs(enabled_languages) do
  local language_pack = require('languages.' .. language)

  merge_tables(packs.lsp, language_pack.lsp)
  append_unique(packs.mason, language_pack.mason)
  merge_tables(packs.formatters_by_ft, language_pack.formatters_by_ft)
  merge_tables(packs.dap_adapters, language_pack.dap_adapters)
  merge_tables(packs.dap_configurations, language_pack.dap_configurations)
  append_unique(packs.treesitter, language_pack.treesitter)
end

function packs.has_language(language)
  return enabled_language_set[language] == true
end

return packs
