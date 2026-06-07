--[[devdeck
{
  "requires": [
    { "bin": ["go"] }
  ]
}
]]
local go_adapter = function(callback, _)
  local port = 38697
  local handle
  handle = vim.uv.spawn('dlv', {
    args = { 'dap', '-l', '127.0.0.1:' .. port },
    detached = true,
  }, function(code)
    if handle then
      handle:close()
    end
    print('Delve exited with code', code)
  end)
  vim.defer_fn(function()
    callback({ type = 'server', host = '127.0.0.1', port = port })
  end, 100)
end

return {
  lsp = {
    gopls = {},
  },
  mason = {
    'delve',
  },
  formatters_by_ft = {
    go = { 'gofmt' },
  },
  dap_adapters = {
    go = go_adapter,
  },
  dap_configurations = {
    go = {
      {
        type = 'go',
        name = 'Debug',
        request = 'launch',
        program = '${file}',
      },
      {
        type = 'go',
        name = 'Debug Package',
        request = 'launch',
        program = '${fileDirname}',
      },
    },
  },
  treesitter = {
    'go',
  },
}
