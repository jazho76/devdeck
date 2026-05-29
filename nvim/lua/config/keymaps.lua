-- Keymaps for better default experience
vim.keymap.set({ 'n', 'v' }, '<Space>', '<Nop>', { silent = true })

-- Remap for dealing with word wrap
vim.keymap.set('n', 'k', "v:count == 0 ? 'gk' : 'k'", { expr = true, silent = true })
vim.keymap.set('n', 'j', "v:count == 0 ? 'gj' : 'j'", { expr = true, silent = true })

-- Buffer
vim.keymap.set('n', '<A-a>', ':bprevious<CR>', { noremap = true, silent = true, desc = 'Go to the previous buffer' })
vim.keymap.set('n', '<A-d>', ':bnext<CR>', { noremap = true, silent = true, desc = 'Go to the next buffer' })
vim.keymap.set('n', '<A-q>', ':bdelete<CR>', { noremap = true, silent = true, desc = 'Delete the current buffer' })
vim.keymap.set('n', '<A-Q>', ':bdelete!<CR>', { silent = true, desc = 'Force-delete the current buffer' })
vim.keymap.set('n', '<A-o>', function()
  local current = vim.api.nvim_get_current_buf()
  for _, buf in ipairs(vim.api.nvim_list_bufs()) do
    if buf ~= current
      and vim.bo[buf].buflisted
      and vim.bo[buf].buftype == ''
      and not vim.bo[buf].modified then
      vim.api.nvim_buf_delete(buf, {})
    end
  end
end, { noremap = true, silent = true, desc = 'Close all buffers except the current buffer' })

-- Split resizing
vim.keymap.set('n', '<C-Up>', ':resize +2<CR>', { noremap = true, silent = true, desc = 'Increase the current window height' })
vim.keymap.set('n', '<C-Down>', ':resize -2<CR>', { noremap = true, silent = true, desc = 'Decrease the current window height' })
vim.keymap.set('n', '<C-Left>', ':vertical resize -2<CR>', { noremap = true, silent = true, desc = 'Decrease the current window width' })
vim.keymap.set('n', '<C-Right>', ':vertical resize +2<CR>', { noremap = true, silent = true, desc = 'Increase the current window width' })

-- Terminal
vim.keymap.set('t', '<C-n>', [[<C-\><C-n>]], { noremap = true, silent = true, desc = 'Leave terminal insert mode and return to normal mode' })

-- Diagnostic
vim.diagnostic.config({
  float = {
    border = 'rounded',
  },
  signs = {
    text = {
      [vim.diagnostic.severity.ERROR] = '󰅚',
      [vim.diagnostic.severity.WARN] = '',
      [vim.diagnostic.severity.HINT] = '',
      [vim.diagnostic.severity.INFO] = '',
    },
  },
})
vim.keymap.set('n', '[d', vim.diagnostic.goto_prev, { desc = 'Go to previous diagnostic message' })
vim.keymap.set('n', ']d', vim.diagnostic.goto_next, { desc = 'Go to next diagnostic message' })
vim.keymap.set('n', '<leader>e', vim.diagnostic.open_float, { desc = 'Open floating diagnostic message' })
vim.keymap.set('n', '<leader>q', vim.diagnostic.setloclist, { desc = 'Open diagnostics list' })

-- Quick fix
vim.keymap.set('n', '[q', ':cprevious<CR>', { desc = 'Go to previous quick fix' })
vim.keymap.set('n', ']q', ':cnext<CR>', { desc = 'Go to next quick fix' })

-- Keep yank when pasting over a selection
vim.keymap.set('x', 'p', '"_dP', { silent = true, desc = 'Paste over the selection without replacing the default register' })

-- Move selected lines up and down
vim.keymap.set('n', '<A-j>', ':m .+1<CR>==', { silent = true, desc = 'Move the current line down' })
vim.keymap.set('n', '<A-k>', ':m .-2<CR>==', { silent = true, desc = 'Move the current line up' })
vim.keymap.set('v', '<A-j>', ":m '>+1<CR>gv=gv", { silent = true, desc = 'Move the selected lines down and keep the selection' })
vim.keymap.set('v', '<A-k>', ":m '<-2<CR>gv=gv", { silent = true, desc = 'Move the selected lines up and keep the selection' })
