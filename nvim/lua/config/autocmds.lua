-- [[ No line number on terminals ]]
vim.api.nvim_create_autocmd('TermOpen', {
  pattern = '*',
  command = 'setlocal nonumber norelativenumber'
})

-- [[ Highlight on yank ]]
local highlight_group = vim.api.nvim_create_augroup('YankHighlight', { clear = true })
vim.api.nvim_create_autocmd('TextYankPost', {
  callback = function()
    vim.highlight.on_yank()
  end,
  group = highlight_group,
  pattern = '*',
})

-- [[ Wipe stray empty buffers ]]
local function wipe_empty_buffers()
  for _, buf in ipairs(vim.api.nvim_list_bufs()) do
    if vim.api.nvim_buf_is_valid(buf)
        and vim.bo[buf].buflisted
        and vim.bo[buf].buftype == ''
        and not vim.bo[buf].modified
        and vim.api.nvim_buf_get_name(buf) == ''
        and #vim.fn.win_findbuf(buf) == 0
        and vim.api.nvim_buf_line_count(buf) == 1
        and vim.api.nvim_buf_get_lines(buf, 0, 1, false)[1] == ''
    then
      pcall(vim.api.nvim_buf_delete, buf, { force = false })
    end
  end
end

local empty_buffer_group = vim.api.nvim_create_augroup('WipeEmptyBuffers', { clear = true })
vim.api.nvim_create_autocmd('BufEnter', {
  group = empty_buffer_group,
  pattern = '*',
  callback = function()
    vim.schedule(wipe_empty_buffers)
  end,
})
