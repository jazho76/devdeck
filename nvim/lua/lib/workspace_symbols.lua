-- LSP servers start from a FileType event, so with no buffer open there is no client
-- to answer workspace/symbol. A hidden buffer holding a real project file starts one;
-- the file must exist on disk, since a pathless buffer resolves no root and indexes
-- nothing. Queries then fan out, because a monorepo runs one client per project root.

local channel = require('plenary.async.control').channel
local actions = require('telescope.actions')
local conf = require('telescope.config').values
local finders = require('telescope.finders')
local make_entry = require('telescope.make_entry')
local pickers = require('telescope.pickers')

local ATTACH_TIMEOUT_MS = 10000
local ATTACH_POLL_MS = 50

local supported_filetypes
local anchor_paths = {}

local function lsp_filetypes()
  if supported_filetypes then
    return supported_filetypes
  end

  supported_filetypes = {}
  for _, config in ipairs(vim.lsp.get_configs()) do
    if vim.lsp.is_enabled(config.name) then
      for _, filetype in ipairs(config.filetypes or {}) do
        supported_filetypes[filetype] = true
      end
    end
  end
  return supported_filetypes
end

local function has_enabled_server(path, filetypes)
  return filetypes[vim.filetype.match({ filename = path })] == true
end

local function recently_edited_file(cwd, filetypes)
  for _, entry in ipairs(vim.v.oldfiles) do
    local path = vim.fs.normalize(entry)
    if vim.fs.relpath(cwd, path) and vim.uv.fs_stat(path) and has_enabled_server(path, filetypes) then
      return path
    end
  end
end

local function directory_depth(path)
  local _, separators = path:gsub('/', '')
  return separators
end

local function dominant_language_file(filetypes)
  local files = vim.fn.systemlist({ 'rg', '--files' })
  if vim.v.shell_error ~= 0 then
    return nil
  end

  local languages = {}
  for _, file in ipairs(files) do
    local extension = file:match('%.([%w_]+)$')
    if extension then
      local language = languages[extension]
      if not language then
        languages[extension] = { count = 1, shallowest_file = file }
      else
        language.count = language.count + 1
        if directory_depth(file) < directory_depth(language.shallowest_file) then
          language.shallowest_file = file
        end
      end
    end
  end

  local dominant
  for _, language in pairs(languages) do
    local outranks_dominant = not dominant or language.count > dominant.count
    if outranks_dominant and has_enabled_server(language.shallowest_file, filetypes) then
      dominant = language
    end
  end

  return dominant and vim.fn.fnamemodify(dominant.shallowest_file, ':p')
end

local function shares_tree(one, other)
  return vim.fs.relpath(one, other) ~= nil or vim.fs.relpath(other, one) ~= nil
end

local function symbol_clients(cwd)
  local clients = {}
  for _, client in ipairs(vim.lsp.get_clients()) do
    if
      client:supports_method('workspace/symbol')
      and client.root_dir
      and shares_tree(cwd, vim.fs.normalize(client.root_dir))
    then
      clients[#clients + 1] = client
    end
  end
  return clients
end

local function server_summary(clients)
  local counts, names = {}, {}
  for _, client in ipairs(clients) do
    if not counts[client.name] then
      names[#names + 1] = client.name
    end
    counts[client.name] = (counts[client.name] or 0) + 1
  end

  return table.concat(
    vim.tbl_map(function(name)
      return counts[name] > 1 and ('%s x%d'):format(name, counts[name]) or name
    end, names),
    ', '
  )
end

local function symbol_requester(cwd)
  local cancel_previous = function() end

  return function(query)
    cancel_previous()

    local clients = symbol_clients(cwd)
    if vim.tbl_isempty(clients) then
      return {}
    end

    local send, receive = channel.oneshot()
    local locations = {}
    local pending = #clients
    local inflight = {}
    local released = false

    local function release()
      if not released then
        released = true
        send()
      end
    end

    local function settle()
      pending = pending - 1
      if pending == 0 then
        release()
      end
    end

    for _, client in ipairs(clients) do
      local started, request_id = client:request('workspace/symbol', { query = query }, function(err, result)
        if not err and result then
          vim.list_extend(locations, vim.lsp.util.symbols_to_items(result, nil, client.offset_encoding))
        end
        inflight[client] = nil
        settle()
      end)

      if started then
        inflight[client] = request_id
      else
        settle()
      end
    end

    cancel_previous = function()
      for client, request_id in pairs(inflight) do
        client:cancel_request(request_id)
      end
      release()
    end

    receive()
    return locations
  end
end

local function open_symbol_picker(cwd, clients)
  local opts = {}
  pickers
    .new(opts, {
      prompt_title = 'Workspace Symbols (' .. server_summary(clients) .. ')',
      finder = finders.new_dynamic({
        entry_maker = make_entry.gen_from_lsp_symbols(opts),
        fn = symbol_requester(cwd),
      }),
      previewer = conf.qflist_previewer(opts),
      sorter = conf.generic_sorter(opts),
      attach_mappings = function(_, map)
        map('i', '<c-space>', actions.to_fuzzy_refine)
        return true
      end,
    })
    :find()
end

local function load_hidden(path)
  local bufnr = vim.fn.bufadd(path)
  vim.bo[bufnr].swapfile = false
  vim.fn.bufload(bufnr)
  vim.bo[bufnr].buflisted = false
end

local function anchor_path(cwd)
  local cached = anchor_paths[cwd]
  if cached and vim.uv.fs_stat(cached) then
    return cached
  end
  if cached == false then
    return nil
  end

  local filetypes = lsp_filetypes()
  local path = recently_edited_file(cwd, filetypes) or dominant_language_file(filetypes)
  anchor_paths[cwd] = path or false
  return path
end

local function has_running_server(path, clients)
  local filetype = vim.filetype.match({ filename = path })
  for _, client in ipairs(clients) do
    if vim.tbl_contains(client.config.filetypes or {}, filetype) then
      return true
    end
  end
  return false
end

local function await_serving_client(cwd, path, on_result)
  local deadline = vim.uv.now() + ATTACH_TIMEOUT_MS

  local function poll()
    local clients = symbol_clients(cwd)
    if has_running_server(path, clients) or vim.uv.now() >= deadline then
      on_result(clients)
    else
      vim.defer_fn(poll, ATTACH_POLL_MS)
    end
  end

  poll()
end

return function()
  local cwd = vim.fs.normalize(vim.fn.getcwd())
  local clients = symbol_clients(cwd)
  local path = anchor_path(cwd)

  if path and not has_running_server(path, clients) then
    vim.notify('Starting LSP server from ' .. vim.fn.fnamemodify(path, ':t'))
    load_hidden(path)
    return await_serving_client(cwd, path, function(ready)
      if vim.tbl_isempty(ready) then
        vim.notify('No LSP server attached for workspace symbols', vim.log.levels.WARN)
      else
        open_symbol_picker(cwd, ready)
      end
    end)
  end

  if vim.tbl_isempty(clients) then
    return vim.notify('No file in this project matches an enabled LSP server', vim.log.levels.WARN)
  end

  open_symbol_picker(cwd, clients)
end
