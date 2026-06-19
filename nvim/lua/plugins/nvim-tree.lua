local function default_tree_width()
  return math.max(30, math.floor(vim.o.columns * 0.20))
end

return {
  'nvim-tree/nvim-tree.lua',
  lazy = false,
  dependencies = {
    {
      'nvim-tree/nvim-web-devicons',
    }
  },
  config = function()
    local remembered_tree_width

    require('nvim-tree').setup {
      update_focused_file = {
        enable = true,
      },
      view = {
        width = default_tree_width(),
        preserve_window_proportions = true,
      },
      actions = {
        open_file = {
          quit_on_open = true
        }
      }
    }

    local api = require('nvim-tree.api')
    local Event = api.events.Event

    local function remember_tree_width()
      local tree_window = api.tree.winid()
      if tree_window and vim.api.nvim_win_is_valid(tree_window) then
        remembered_tree_width = vim.api.nvim_win_get_width(tree_window)
      end
    end

    local function restore_tree_width()
      api.tree.resize({ absolute = remembered_tree_width or default_tree_width() })
    end

    local function toggle_tree()
      if api.tree.is_visible() then
        remember_tree_width()
        api.tree.close()
      else
        api.tree.open()
        restore_tree_width()
      end
    end

    api.events.subscribe(Event.TreeClose, remember_tree_width)

    vim.keymap.set('n', '<A-e>', toggle_tree, { noremap = true, silent = true, desc = 'Toggle the NvimTree file explorer' })
  end,
}
