return {
  'folke/which-key.nvim',
  dependencies = {
    {
      'nvim-tree/nvim-web-devicons',
    },
    {
      'echasnovski/mini.nvim',
    }
  },
  config = function()
    local wk = require('which-key')
    wk.add({
      { '<leader>c',  group = 'Code' },
      { '<leader>d',  group = 'Debug' },
      { '<leader>g',  group = 'Git' },
      { '<leader>h',  group = 'Git hunk', mode = 'n' },
      { '<leader>h',  group = 'Git hunk', mode = 'v' },
      { '<leader>l',  group = 'LSP' },
      { '<leader>n',  group = 'Notifications' },
      { '<leader>s',  group = 'Search' },
      { '<leader>t',  group = 'Toggle' },
      { '<leader>T',  group = 'Tests' },
      { '<leader>w',  group = 'Swap' },
    }
    )
  end
}
