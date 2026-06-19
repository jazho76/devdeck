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
    local api = require('nvim-tree.api')
    local Event = api.events.Event
    local remembered_tree_width

    local function remember_tree_width(tree_window)
      tree_window = tree_window or api.tree.winid()
      if tree_window and vim.api.nvim_win_is_valid(tree_window) then
        remembered_tree_width = vim.api.nvim_win_get_width(tree_window)
      end
    end

    local function restore_tree_width()
      api.tree.resize({ absolute = remembered_tree_width or default_tree_width() })
    end

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

    api.events.subscribe(Event.TreeOpen, restore_tree_width)
    vim.api.nvim_create_autocmd('WinLeave', {
      callback = function(event)
        if vim.bo[event.buf].filetype == 'NvimTree' then
          remember_tree_width(vim.api.nvim_get_current_win())
        end
      end,
    })

    local function toggle_tree()
      if api.tree.is_visible() then
        remember_tree_width()
        api.tree.close()
      else
        api.tree.open()
      end
    end

    vim.keymap.set('n', '<A-e>', toggle_tree, { noremap = true, silent = true, desc = 'Toggle the NvimTree file explorer' })
  end,
}
