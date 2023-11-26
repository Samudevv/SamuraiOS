local function bootstrap_pckr()
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
end
local function remember_cursor_position()
    vim.api.nvim_create_augroup('UserGroup', {})
    -- From vim defaults.vim
    -- ---
    -- When editing a file, always jump to the last known cursor position.
    -- Don't do it when the position is invalid, when inside an event handler
    -- (happens when dropping a file on gvim) and for a commit message (it's
    -- likely a different one than last time).
    vim.api.nvim_create_autocmd('BufReadPost', {
        group = 'UserGroup',
        callback = function(args)
            local valid_line = vim.fn.line([['"]]) >= 1 and vim.fn.line([['"]]) < vim.fn.line('$')
            local not_commit = vim.b[args.buf].filetype ~= 'commit'

        if valid_line and not_commit then
            vim.cmd([[normal! g`"]])
        end
      end,
    })
end

local require_failed = false
local function prequire(module)
    local ok, m = pcall(require, module)
    if not ok then
        require_failed = true
        return nil
    end
    return m
end

local function cmd(args)
    local value = 'command! -nargs=0 ' .. args
    vim.cmd(value)
end
local function cf(pattern, command)
    vim.api.nvim_create_autocmd({'BufEnter', 'BufWinEnter'}, {
        pattern = pattern,
        command = command,
    })
end
local function ac(event, func)
    vim.api.nvim_create_autocmd(event, {
        callback = func,
    })
end

bootstrap_pckr()

local bind = vim.keymap.set
local c = vim.cmd
local slt = { silent = true }
local normal = 'n'
local all = {'n', 'i', 'v'}
local pckr = require('pckr')
local lang_c = {'*.c', '*.cpp', '*.cxx', '*.h', '*.hpp', '*.hxx'}
local lang_go = {'*.go'}

--  ____  _             _
-- |  _ \| |_   _  __ _(_)_ __  ___
-- | |_) | | | | |/ _` | | '_ \/ __|
-- |  __/| | |_| | (_| | | | | \__ \
-- |_|   |_|\__,_|\__, |_|_| |_|___/
--                |___/
--
pckr.add{
    'Mofiqul/dracula.nvim';
    'mg979/vim-visual-multi';
    'ibhagwan/fzf-lua';
    'akinsho/bufferline.nvim';
    {'nvim-treesitter/nvim-treesitter', run = function() require('nvim-treesitter.install').update({with_sync = true})() end};
    'cappyzawa/trim.nvim';
    'mhinz/vim-signify';
}

local fzf = prequire('fzf-lua')
local bufferline = prequire('bufferline')
local treesitter = prequire('nvim-treesitter.configs')
local trim = prequire('trim')

if not require_failed then
    fzf.setup({'fzf-native'})
    vim.opt.termguicolors = true
    bufferline.setup{}
    treesitter.setup{
        ensure_installed = {
            'bash',
            'bibtex',
            'c',
            'cpp',
            'cmake',
            'go',
            'lua',
            'odin',
        },
        sync_install = false,
        auto_install = false,
        highlight = {
            enable = true
        }
    }
    trim.setup{
        trim_last_line = false,
        trim_first_line = false,
    }
--  _  __          _                         _
-- | |/ /___ _   _| |__   ___   __ _ _ __ __| |
-- | ' // _ \ | | | '_ \ / _ \ / _` | '__/ _` |
-- | . \  __/ |_| | |_) | (_) | (_| | | | (_| |
-- |_|\_\___|\__, |_.__/ \___/ \__,_|_|  \__,_|
--           |___/
--  ____  _           _ _
-- | __ )(_)_ __   __| (_)_ __   __ _ ___
-- |  _ \| | '_ \ / _` | | '_ \ / _` / __|
-- | |_) | | | | | (_| | | | | | (_| \__ \
-- |____/|_|_| |_|\__,_|_|_| |_|\__, |___/
--                              |___/
    bind(normal, '<c-P>', fzf.git_files,                                 slt)
    bind(normal, '<Tab>', '<cmd>bnext<CR>',                              slt)
    bind(all,    '<c-W>', '<cmd>bd<CR>',                                 slt)
    bind(all,    '<c-Q>', '<cmd>q<CR>',                                  slt)
    bind(all,    '<c-S>', '<cmd>w<CR>',                                  slt)
--   ____          _
--  / ___|   _ ___| |_ ___  _ __ ___
-- | |  | | | / __| __/ _ \| '_ ` _ \
-- | |__| |_| \__ \ || (_) | | | | | |
--  \____\__,_|___/\__\___/|_| |_| |_|
--
--   ____                                          _
--  / ___|___  _ __ ___  _ __ ___   __ _ _ __   __| |___
-- | |   / _ \| '_ ` _ \| '_ ` _ \ / _` | '_ \ / _` / __|
-- | |__| (_) | | | | | | | | | | | (_| | | | | (_| \__ \
--  \____\___/|_| |_| |_|_| |_| |_|\__,_|_| |_|\__,_|___/
--
    cmd('W write')
    cmd('Q quit')
    cmd('SRC source %')
--   ____                                          _
--  / ___|___  _ __ ___  _ __ ___   __ _ _ __   __| |___
-- | |   / _ \| '_ ` _ \| '_ ` _ \ / _` | '_ \ / _` / __|
-- | |__| (_) | | | | | | | | | | | (_| | | | | (_| \__ \
--  \____\___/|_| |_| |_|_| |_| |_|\__,_|_| |_|\__,_|___/
--
    c('set relativenumber')
    c('set number')
    c('colorscheme dracula')
    c('set listchars=space:·,tab:-->')
    c('set list')
    c('set tabstop=8 softtabstop=0')
    c('set shiftwidth=4 smarttab')
    c('set expandtab')

    cf(lang_c,  'set shiftwidth=2')
    cf(lang_go, 'set noexpandtab tabstop=4')

    remember_cursor_position()
end

