--             _
--  _ __   ___| | ___ __
-- | '_ \ / __| |/ / '__|
-- | |_) | (__|   <| |
-- | .__/ \___|_|\_\_|
-- |_|
--
-- bootstrap packr
local pckr_path = vim.fn.stdpath('data') .. '/pckr/pckr.nvim'
if not vim.loop.fs_stat(pckr_path) then
    vim.fn.system({
        'git',
        'clone',
        '--filter=blob:none',
        'https://github.com/lewis6991/pckr.nvim',
        pckr_path
    })
end
vim.opt.rtp:prepend(pckr_path)

pckr = require('pckr')

