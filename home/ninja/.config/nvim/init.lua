require("prelude")

--  ____  _             _
-- |  _ \| |_   _  __ _(_)_ __  ___
-- | |_) | | | | |/ _` | | '_ \/ __|
-- |  __/| | |_| | (_| | | | | \__ \
-- |_|   |_|\__,_|\__, |_|_| |_|___/
--                |___/
--
local plugins = {
    'mg979/vim-visual-multi';
    'cappyzawa/trim.nvim';
}
local term_plugins = {
    'Mofiqul/dracula.nvim';
    'ibhagwan/fzf-lua';
    'akinsho/bufferline.nvim';
    {'nvim-treesitter/nvim-treesitter', run = function() require('nvim-treesitter.install').update({with_sync = true})() end};
    'mhinz/vim-signify';
    'sbdchd/neoformat';
}
if not is_vscode then
    for _, value in ipairs(term_plugins) do
        table.insert(plugins, value)
    end
end

pckr.add(plugins)


if not is_vscode then
    bufferline = prequire('bufferline')
    treesitter = prequire('nvim-treesitter.configs')
    fzf = prequire('fzf-lua')
end

trim = prequire('trim')

if require_failed then
    os.exit(0)
end

if not is_vscode then
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
            'fish',
        },
        sync_install = false,
        auto_install = false,
        highlight = {
            enable = true
        },
    }

    fzf.setup({'fzf-native'})
    vim.opt.termguicolors = true
end

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
if not is_vscode then
    bind(normal, '<c-P>',   function() fzf.files({ cmd = os.getenv('FZF_DEFAULT_COMMAND')}) end, slt)
    bind(normal, '<S-P>',   fzf.blines,                                                          slt)
    bind(all,    '<c-Q>',   '<cmd>q<CR>',                                                        slt)
    bind(all,    '<c-S>',   '<cmd>w<CR>',                                                        slt)
end
bind(normal, '<Tab>',   '<cmd>bnext<CR>',                                                    slt)
bind(normal, '<S-Tab>', '<cmd>bprevious<CR>',                                                slt)
bind(all,    '<c-E>',   '<cmd>bd<CR>',                                                       slt)
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
term_c('set number')
term_c('colorscheme dracula')
term_c('set listchars=space:·,tab:-->')
term_c('set list')
c('set tabstop=8 softtabstop=0')
c('set shiftwidth=4 smarttab')
c('set expandtab')

cf(lang_c,  'set shiftwidth=2')
cf(lang_go, 'set noexpandtab tabstop=4')
term_ac('BufWritePre', 'silent! undojoin | Neoformat')

if not is_vscode then
    remember_cursor_position()
end

