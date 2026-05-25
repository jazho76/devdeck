return {
  {
    'mfussenegger/nvim-dap',
    config = function()
      local dap = require('dap')

      -- Set keymaps to control the debugger
      vim.keymap.set('n', '<F5>', dap.continue)
      vim.keymap.set('n', '<F6>', dap.step_over)
      vim.keymap.set('n', '<F7>', dap.step_into)
      vim.keymap.set('n', '<F8>', dap.step_out)
      vim.keymap.set('n', '<leader>db', dap.toggle_breakpoint, { desc = '[D]ebug Toggle [B]reakpoint' })
      vim.keymap.set('n', '<leader>dB', function()
        dap.set_breakpoint(vim.fn.input('Breakpoint condition: '))
      end, { desc = '[D]ebug Toggle [B]reakpoint With Condition' })

      local language_packs = require('config.language-packs')
      for name, adapter in pairs(language_packs.dap_adapters) do
        dap.adapters[name] = adapter
      end

      for language, configurations in pairs(language_packs.dap_configurations) do
        dap.configurations[language] = configurations
      end
    end
  },
  {
    'rcarriga/nvim-dap-ui',
    dependencies = {
      {
        'nvim-neotest/nvim-nio',
      },
    },
    config = function()
      require('dapui').setup()
      local dap, dapui = require('dap'), require('dapui')

      dap.listeners.after.event_initialized['dapui_config'] = function()
        dapui.open({})
      end
      dap.listeners.before.event_terminated['dapui_config'] = function()
        dapui.close({})
      end
      dap.listeners.before.event_exited['dapui_config'] = function()
        dapui.close({})
      end

      vim.keymap.set('n', '<leader>du', require 'dapui'.toggle, { desc = '[D]ebug Toggle [U]I' })
    end
  },
  {
    -- JS debugger
    'microsoft/vscode-js-debug',
    enabled = function()
      return require('config.language-packs').has_language('web')
    end,
    build = table.concat({
      'curl -fL https://github.com/microsoft/vscode-js-debug/releases/download/v1.117.0/js-debug-dap-v1.117.0.tar.gz -o js-debug-dap.tar.gz',
      'tar xzf js-debug-dap.tar.gz --strip-components=1',
      'rm js-debug-dap.tar.gz',
    }, ' && '),
  },
}
