return {
  'neovim/nvim-lspconfig',
  dependencies = {
    {
      'williamboman/mason.nvim',
    },
    {
      'williamboman/mason-lspconfig.nvim',
    },
    {
      'WhoIsSethDaniel/mason-tool-installer.nvim',
    },
    {
      'j-hui/fidget.nvim',
      opts = {
        notification = {
          window = {
            winblend = 0,
          },
        },
      },
    },
    {
      'folke/lazydev.nvim',
      ft = 'lua',
      opts = {
        library = {
          { path = '${3rd}/luv/library', words = { 'vim%.uv' } },
        },
      },
    },
  },
  config = function()
    local on_attach = function(_, bufnr)
      local nmap = function(keys, func, desc)
        if desc then
          desc = 'LSP: ' .. desc
        end

        vim.keymap.set('n', keys, func, { buffer = bufnr, desc = desc })
      end

      nmap('<leader>cr', vim.lsp.buf.rename, 'Rename symbol')
      nmap('<leader>ca', vim.lsp.buf.code_action, 'Code action')

      nmap('<leader>lr', '<cmd>lsp restart<CR>', 'Restart server')
      nmap('<leader>li', '<cmd>checkhealth vim.lsp<CR>', 'Server info')
      nmap('<leader>ld', function()
        vim.diagnostic.enable(not vim.diagnostic.is_enabled())
      end, 'Toggle diagnostics')

      nmap('gd', require('telescope.builtin').lsp_definitions, 'Goto definition')
      nmap('gr', require('telescope.builtin').lsp_references, 'Goto references')
      nmap('gI', require('telescope.builtin').lsp_implementations, 'Goto implementation')
      nmap('gy', require('telescope.builtin').lsp_type_definitions, 'Goto type definition')
      nmap('gD', vim.lsp.buf.declaration, 'Goto declaration')

      nmap('<leader>ss', require('telescope.builtin').lsp_document_symbols, 'Document symbols')
      nmap('<leader>sS', require('telescope.builtin').lsp_dynamic_workspace_symbols, 'Workspace symbols')
      nmap('K', vim.lsp.buf.hover, 'Hover documentation')
    end

    vim.o.winborder = 'rounded'

    require('mason').setup()

    local toolsets = require('config.toolset-registry')
    local servers = vim.tbl_deep_extend('force', {
      lua_ls = {
        settings = {
          Lua = {
            workspace = { checkThirdParty = false },
            telemetry = { enable = false },
            diagnostics = { disable = { 'missing-fields' } },
          },
        },
      },
    }, toolsets.lsp)

    local capabilities = vim.lsp.protocol.make_client_capabilities()
    capabilities = require('cmp_nvim_lsp').default_capabilities(capabilities)

    local mason_lspconfig = require('mason-lspconfig')
    mason_lspconfig.setup({
      ensure_installed = vim.tbl_keys(servers),
      automatic_enable = false,
    })

    for server_name, server_config in pairs(servers) do
      local config = vim.tbl_deep_extend('force', {
        capabilities = capabilities,
        on_attach = on_attach,
      }, server_config)

      vim.lsp.config(server_name, config)
      vim.lsp.enable(server_name)
    end

    require('mason-tool-installer').setup({
      ensure_installed = toolsets.mason,
      run_on_start = true,
    })
  end,
}
