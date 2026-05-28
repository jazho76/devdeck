return {
  'christoomey/vim-tmux-navigator',
  cmd = {
    'TmuxNavigateLeft',
    'TmuxNavigateDown',
    'TmuxNavigateUp',
    'TmuxNavigateRight',
    'TmuxNavigatePrevious',
  },
  keys = {
    { '<c-h>',  '<cmd><C-U>TmuxNavigateLeft<cr>',     desc = 'Move focus to the pane or Neovim split on the left' },
    { '<c-j>',  '<cmd><C-U>TmuxNavigateDown<cr>',     desc = 'Move focus to the pane or Neovim split below' },
    { '<c-k>',  '<cmd><C-U>TmuxNavigateUp<cr>',       desc = 'Move focus to the pane or Neovim split above' },
    { '<c-l>',  '<cmd><C-U>TmuxNavigateRight<cr>',    desc = 'Move focus to the pane or Neovim split on the right' },
    { '<c-\\>', '<cmd><C-U>TmuxNavigatePrevious<cr>', desc = 'Move focus to the previously active tmux pane or Neovim split' },
  },
}
