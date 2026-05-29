return {
  'stevearc/conform.nvim',
  event = { 'BufWritePre' },
  cmd = { 'ConformInfo', 'Format' },
  init = function()
    vim.g.autoformat_enabled = true

    local toggle_autoformat = function()
      vim.g.autoformat_enabled = not vim.g.autoformat_enabled
      local status = vim.g.autoformat_enabled and 'enabled' or 'disabled'
      vim.notify('Autoformatting is ' .. status, vim.log.levels.INFO, { title = 'Toggle Autoformat' })
    end

    vim.keymap.set('n', '<Leader>tf', toggle_autoformat, { desc = 'Toggle autoformat' })
  end,
  opts = {
    default_format_opts = {
      lsp_format = 'fallback',
      quiet = true,
      timeout_ms = 3000,
    },
    format_on_save = function()
      if not vim.g.autoformat_enabled then
        return nil
      end

      return {
        lsp_format = 'fallback',
        quiet = true,
        timeout_ms = 3000,
      }
    end,
    notify_on_error = false,
    notify_no_formatters = false,
    formatters_by_ft = require('config.toolset-registry').formatters_by_ft,
    formatters = {
      prettier = {
        prepend_args = function(_, ctx)
          if vim.bo[ctx.buf].filetype == 'scss' then
            return { '--parser', 'scss' }
          end
          return {}
        end,
      },
    },
  },
  config = function(_, opts)
    local conform = require('conform')
    conform.setup(opts)

    vim.api.nvim_create_user_command('Format', function()
      conform.format({ bufnr = 0 })
    end, { desc = 'Format current buffer' })
  end,
}
