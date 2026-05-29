return {
  'tpope/vim-fugitive',
  config = function()
    vim.keymap.set('n', '<leader>gg', ':G | only<CR>', { desc = 'Open Fugitive status' })
  end
}
