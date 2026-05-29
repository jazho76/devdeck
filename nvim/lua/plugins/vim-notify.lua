return {
  'rcarriga/nvim-notify',
  config = function()
    vim.notify = require('notify')

    vim.keymap.set('n', '<leader>nh', function()
      require('notify')._print_history()
    end, { desc = 'Show notification history' })

    vim.keymap.set('n', '<leader>nd', function()
      require('notify').dismiss({})
    end, { desc = 'Dismiss all notifications' })
  end
}
