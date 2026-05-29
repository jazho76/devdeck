return {
  'tpope/vim-obsession',
  config = function()
    vim.keymap.set('n', '<leader>to', '<cmd>Obsession<CR>', { desc = 'Toggle Obsession session' })
  end,
}
