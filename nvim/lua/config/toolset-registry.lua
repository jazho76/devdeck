local enabled_toolsets = require('config.toolsets')

local enabled_toolset_set = {}
for _, name in ipairs(enabled_toolsets) do
  enabled_toolset_set[name] = true
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

local registry = {
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

for _, name in ipairs(enabled_toolsets) do
  local module = require('toolsets.' .. name)

  merge_tables(registry.lsp, module.lsp)
  append_unique(registry.mason, module.mason)
  merge_tables(registry.formatters_by_ft, module.formatters_by_ft)
  merge_tables(registry.dap_adapters, module.dap_adapters)
  merge_tables(registry.dap_configurations, module.dap_configurations)
  append_unique(registry.treesitter, module.treesitter)
end

function registry.has_toolset(name)
  return enabled_toolset_set[name] == true
end

return registry
