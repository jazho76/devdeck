local ok, override = pcall(require, 'config.toolsets-local')
if ok and type(override) == 'table' then
  return override
end

return require('config.toolsets-default')
