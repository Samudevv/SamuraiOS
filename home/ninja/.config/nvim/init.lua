require("prelude")

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
    'sbdchd/neoformat';
}

fzf = prequire('fzf-lua')
bufferline = prequire('bufferline')
treesitter = prequire('nvim-treesitter.configs')
trim = prequire('trim')

if require_failed then
    os.exit(0)
end

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
        'fish',
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
bind(normal, '<c-P>',   function() fzf.files({ cmd = os.getenv('FZF_DEFAULT_COMMAND')}) end, slt)
bind(normal, '<S-P>',   fzf.blines,                                                          slt)
bind(normal, '<Tab>',   '<cmd>bnext<CR>',                                                    slt)
bind(normal, '<S-Tab>', '<cmd>bprevious<CR>',                                                slt)
bind(all,    '<c-E>',   '<cmd>bd<CR>',                                                       slt)
bind(all,    '<c-Q>',   '<cmd>q<CR>',                                                        slt)
bind(all,    '<c-S>',   '<cmd>w<CR>',                                                        slt)
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

ac('BufWritePre', 'silent! undojoin | Neoformat')

remember_cursor_position()

