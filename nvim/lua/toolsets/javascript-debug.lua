--[[devdeck
{
  "requires": [
    { "bin": ["node"] },
    { "bin": ["npm"] }
  ]
}
]]
local js_based_languages = {
  'typescript',
  'javascript',
  'typescriptreact',
  'javascriptreact',
}

local js_debug_adapter = vim.fn.stdpath('data') .. '/mason/bin/js-debug-adapter'

local dap_adapters = {}
for _, adapter in ipairs({ 'pwa-node', 'pwa-chrome', 'pwa-msedge', 'node-terminal' }) do
  dap_adapters[adapter] = {
    type = 'server',
    host = 'localhost',
    port = '${port}',
    executable = {
      command = js_debug_adapter,
      args = { '${port}' },
    },
  }
end

local dap_configurations = {}
for _, language in ipairs(js_based_languages) do
  dap_configurations[language] = {
    {
      type = 'pwa-node',
      request = 'launch',
      name = 'Launch file',
      program = '${file}',
      cwd = '${workspaceFolder}',
    },
    {
      type = 'pwa-node',
      request = 'attach',
      name = 'Attach',
      processId = function(...)
        return require('dap.utils').pick_process(...)
      end,
      cwd = '${workspaceFolder}',
    },
    {
      type = 'pwa-chrome',
      request = 'launch',
      name = 'Start Chrome with localhost',
      url = 'http://localhost:3000',
      webRoot = '${workspaceFolder}',
      userDataDir = '${workspaceFolder}/.vscode/vscode-chrome-debug-userdatadir',
    },
  }
end

return {
  mason = {
    'js-debug-adapter',
  },
  dap_adapters = dap_adapters,
  dap_configurations = dap_configurations,
}
