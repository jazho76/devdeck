return {
  'github/copilot.vim',
  config = function()
    vim.g.copilot_assume_mapped = true
    vim.keymap.set('i', '<C-S>', 'copilot#Accept("\\<CR>")', {
      silent = true,
      expr = true,
      replace_keycodes = false,
      desc = 'Accept the current Copilot inline suggestion',
    })
  end,
}
