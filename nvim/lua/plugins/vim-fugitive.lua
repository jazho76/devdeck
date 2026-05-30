return {
  'tpope/vim-fugitive',
  config = function()
    vim.keymap.set('n', '<leader>gg', ':G | only<CR>', { desc = 'Open Fugitive status' })
    vim.keymap.set('n', '<leader>gp', ':Git push<CR>', { desc = 'Git push' })
    vim.keymap.set('n', '<leader>gP', ':Git pull<CR>', { desc = 'Git pull' })
    vim.keymap.set('n', '<leader>gF', ':Git fetch<CR>', { desc = 'Git fetch' })
    vim.keymap.set('n', '<leader>gB', ':Git blame<CR>', { desc = 'Git blame' })
  end
}
