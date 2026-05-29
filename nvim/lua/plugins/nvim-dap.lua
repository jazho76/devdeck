return {
  {
    'mfussenegger/nvim-dap',
    config = function()
      local dap = require('dap')

      -- Set keymaps to control the debugger.
      -- Stepping is bound to both <leader>d... (discoverable) and the F-keys
      -- (single-press, for fast stepping during an active session).
      vim.keymap.set('n', '<leader>dc', dap.continue, { desc = 'Continue' })
      vim.keymap.set('n', '<leader>dn', dap.step_over, { desc = 'Step over' })
      vim.keymap.set('n', '<leader>di', dap.step_into, { desc = 'Step into' })
      vim.keymap.set('n', '<leader>do', dap.step_out, { desc = 'Step out' })
      vim.keymap.set('n', '<F5>', dap.continue, { desc = 'Debug: continue' })
      vim.keymap.set('n', '<F6>', dap.step_over, { desc = 'Debug: step over' })
      vim.keymap.set('n', '<F7>', dap.step_into, { desc = 'Debug: step into' })
      vim.keymap.set('n', '<F8>', dap.step_out, { desc = 'Debug: step out' })
      vim.keymap.set('n', '<leader>db', dap.toggle_breakpoint, { desc = 'Toggle breakpoint' })
      vim.keymap.set('n', '<leader>dB', function()
        dap.set_breakpoint(vim.fn.input('Breakpoint condition: '))
      end, { desc = 'Set conditional breakpoint' })

      local toolsets = require('config.toolset-registry')
      for name, adapter in pairs(toolsets.dap_adapters) do
        dap.adapters[name] = adapter
      end

      for language, configurations in pairs(toolsets.dap_configurations) do
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

      vim.keymap.set('n', '<leader>du', require 'dapui'.toggle, { desc = 'Toggle debug UI' })
    end
  },
}
