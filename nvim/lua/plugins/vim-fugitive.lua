return {
  'tpope/vim-fugitive',
  config = function()
    vim.keymap.set('n', '<leader>gg', ':Git | only<CR>', { desc = 'Open Fugitive status' })
    vim.keymap.set('n', '<leader>gp', ':Git! push<CR>', { desc = 'Git push' })
    vim.keymap.set('n', '<leader>gP', ':Git! pull<CR>', { desc = 'Git pull' })
    vim.keymap.set('n', '<leader>gF', ':Git! fetch<CR>', { desc = 'Git fetch' })
    vim.keymap.set('n', '<leader>gB', ':Git blame<CR>', { desc = 'Git blame' })
    local function browse(range)
      local remote = vim.env.SSH_TTY or vim.env.SSH_CONNECTION
      local cmd = remote and 'GBrowse!' or 'GBrowse'
      vim.cmd((range and "'<,'>" or '') .. cmd)
    end
    vim.keymap.set('n', '<leader>go', function() browse(false) end, { desc = 'Open git web URL (copy when remote)' })
    vim.keymap.set('x', '<leader>go', function() browse(true) end,
      { desc = 'Open git web URL for selection (copy when remote)' })
  end
}
