-- Host packages:
--
-- Fedora: nodejs npm
-- Debian/Ubuntu: nodejs npm
local js_based_languages = {
  'typescript',
  'javascript',
  'typescriptreact',
  'javascriptreact',
}

local js_debug_path = vim.fn.stdpath('data') .. '/lazy/vscode-js-debug/src/dapDebugServer.js'

local dap_adapters = {}
for _, adapter in ipairs({ 'pwa-node', 'pwa-chrome', 'pwa-msedge', 'node-terminal' }) do
  dap_adapters[adapter] = {
    type = 'server',
    host = 'localhost',
    port = '${port}',
    executable = {
      command = 'node',
      args = { js_debug_path, '${port}' },
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
  lsp = {
    ts_ls = {},
    html = { filetypes = { 'html', 'twig', 'hbs' } },
    angularls = {},
    cssls = {},
    eslint = {},
  },
  mason = {
    'prettier',
  },
  formatters_by_ft = {
    javascript = { 'prettier' },
    javascriptreact = { 'prettier' },
    typescript = { 'prettier' },
    typescriptreact = { 'prettier' },
    css = { 'prettier' },
    scss = { 'prettier' },
    html = { 'prettier' },
    graphql = { 'prettier' },
  },
  dap_adapters = dap_adapters,
  dap_configurations = dap_configurations,
  treesitter = {
    'tsx',
    'javascript',
    'typescript',
    'graphql',
  },
}
