local parsers = require('config.toolset-registry').treesitter

local function start_treesitter()
  local bufnr = vim.api.nvim_get_current_buf()
  local filetype = vim.bo[bufnr].filetype
  local language = vim.treesitter.language.get_lang(filetype) or filetype

  if not language or language == '' then
    return
  end

  if vim.treesitter.language.add(language) then
    vim.treesitter.start(bufnr, language)
    vim.bo[bufnr].indentexpr = "v:lua.require'nvim-treesitter'.indentexpr()"
  end
end

local function map_select(keys, query, desc)
  vim.keymap.set({ 'x', 'o' }, keys, function()
    require('nvim-treesitter-textobjects.select').select_textobject(query, 'textobjects')
  end, { desc = desc })
end

local function map_move(keys, query, direction, desc)
  vim.keymap.set({ 'n', 'x', 'o' }, keys, function()
    require('nvim-treesitter-textobjects.move')[direction](query, 'textobjects')
  end, { desc = desc })
end

return {
  'nvim-treesitter/nvim-treesitter',
  dependencies = {
    'nvim-treesitter/nvim-treesitter-textobjects',
  },
  build = function()
    require('nvim-treesitter').install(parsers):wait(300000)
  end,
  config = function()
    require('nvim-treesitter').setup()

    if #vim.api.nvim_list_uis() > 0 then
      require('nvim-treesitter').install(parsers)
    end

    vim.api.nvim_create_autocmd('FileType', {
      group = vim.api.nvim_create_augroup('TreesitterFeatures', { clear = true }),
      pattern = parsers,
      callback = start_treesitter,
    })

    require('nvim-treesitter-textobjects').setup({
      select = {
        lookahead = true,
      },
      move = {
        set_jumps = true,
      },
    })

    map_select('aa', '@parameter.outer', 'Select the outer parameter textobject')
    map_select('ia', '@parameter.inner', 'Select the inner parameter textobject')
    map_select('af', '@function.outer', 'Select the outer function textobject')
    map_select('if', '@function.inner', 'Select the inner function textobject')
    map_select('ac', '@class.outer', 'Select the outer class textobject')
    map_select('ic', '@class.inner', 'Select the inner class textobject')
    map_select('ar', '@parameter.outer', 'Select the outer parameter textobject')
    map_select('ab', '@block.outer', 'Select the outer block textobject')
    map_select('ib', '@block.inner', 'Select the inner block textobject')
    map_select('ak', '@assignment.lhs', 'Select the assignment left-hand side textobject')
    map_select('ik', '@assignment.lhs', 'Select the assignment left-hand side textobject')
    map_select('av', '@assignment.rhs', 'Select the assignment right-hand side textobject')
    map_select('iv', '@assignment.rhs', 'Select the assignment right-hand side textobject')

    map_move(']a', '@parameter.inner', 'goto_next_start', 'Jump to the next start of the inner parameter textobject')
    map_move(']b', '@block.inner', 'goto_next_start', 'Jump to the next start of the inner block textobject')
    map_move(']f', '@function.outer', 'goto_next_start', 'Jump to the next start of the outer function textobject')
    map_move(']k', '@assignment.lhs', 'goto_next_start', 'Jump to the next start of the assignment left-hand side textobject')
    map_move(']v', '@assignment.rhs', 'goto_next_start', 'Jump to the next start of the assignment right-hand side textobject')

    map_move(']A', '@parameter.inner', 'goto_next_end', 'Jump to the next end of the inner parameter textobject')
    map_move(']B', '@block.inner', 'goto_next_end', 'Jump to the next end of the inner block textobject')
    map_move(']F', '@function.outer', 'goto_next_end', 'Jump to the next end of the outer function textobject')
    map_move(']K', '@assignment.lhs', 'goto_next_end', 'Jump to the next end of the assignment left-hand side textobject')
    map_move(']V', '@assignment.rhs', 'goto_next_end', 'Jump to the next end of the assignment right-hand side textobject')

    map_move('[a', '@parameter.inner', 'goto_previous_start', 'Jump to the previous start of the inner parameter textobject')
    map_move('[b', '@block.inner', 'goto_previous_start', 'Jump to the previous start of the inner block textobject')
    map_move('[f', '@function.outer', 'goto_previous_start', 'Jump to the previous start of the outer function textobject')
    map_move('[k', '@assignment.lhs', 'goto_previous_start', 'Jump to the previous start of the assignment left-hand side textobject')
    map_move('[v', '@assignment.rhs', 'goto_previous_start', 'Jump to the previous start of the assignment right-hand side textobject')

    map_move('[A', '@parameter.inner', 'goto_previous_end', 'Jump to the previous end of the inner parameter textobject')
    map_move('[B', '@block.inner', 'goto_previous_end', 'Jump to the previous end of the inner block textobject')
    map_move('[F', '@function.outer', 'goto_previous_end', 'Jump to the previous end of the outer function textobject')
    map_move('[K', '@assignment.lhs', 'goto_previous_end', 'Jump to the previous end of the assignment left-hand side textobject')
    map_move('[V', '@assignment.rhs', 'goto_previous_end', 'Jump to the previous end of the assignment right-hand side textobject')

    vim.keymap.set('n', '<leader>wa', function()
      require('nvim-treesitter-textobjects.swap').swap_next('@parameter.inner')
    end, { desc = 'Swap the current parameter with the next parameter' })
    vim.keymap.set('n', '<leader>wf', function()
      require('nvim-treesitter-textobjects.swap').swap_next('@function.outer')
    end, { desc = 'Swap the current function with the next function' })
    vim.keymap.set('n', '<leader>wA', function()
      require('nvim-treesitter-textobjects.swap').swap_previous('@parameter.inner')
    end, { desc = 'Swap the current parameter with the previous parameter' })
    vim.keymap.set('n', '<leader>wF', function()
      require('nvim-treesitter-textobjects.swap').swap_previous('@function.outer')
    end, { desc = 'Swap the current function with the previous function' })

    local ts_repeat_move = require('nvim-treesitter-textobjects.repeatable_move')
    vim.keymap.set({ 'n', 'x', 'o' }, ';', ts_repeat_move.repeat_last_move_next, { desc = 'Repeat the last Treesitter textobject move forward' })
    vim.keymap.set({ 'n', 'x', 'o' }, ',', ts_repeat_move.repeat_last_move_previous, { desc = 'Repeat the last Treesitter textobject move backward' })
    vim.keymap.set({ 'n', 'x', 'o' }, 'f', ts_repeat_move.builtin_f_expr, { expr = true, desc = 'Use repeatable forward character search' })
    vim.keymap.set({ 'n', 'x', 'o' }, 'F', ts_repeat_move.builtin_F_expr, { expr = true, desc = 'Use repeatable backward character search' })
    vim.keymap.set({ 'n', 'x', 'o' }, 't', ts_repeat_move.builtin_t_expr, { expr = true, desc = 'Use repeatable forward till-character search' })
    vim.keymap.set({ 'n', 'x', 'o' }, 'T', ts_repeat_move.builtin_T_expr, { expr = true, desc = 'Use repeatable backward till-character search' })
  end,
}
