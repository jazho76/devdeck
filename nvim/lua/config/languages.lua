local ok, override = pcall(require, 'config.languages-local')
if ok and type(override) == 'table' then
  return override
end

return require('config.languages-default')
