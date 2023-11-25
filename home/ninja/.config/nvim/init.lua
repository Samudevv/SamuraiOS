local function bootstrap_pckr()
	local pckr_path = vim.fn.stdpath("data") .. "/pckr/pckr.nvim"
	
	if not vim.loop.fs_stat(pckr_path) then
		vim.fn.system({
			'git',
			'clone',
			"--filter=blob:none",
			'https://github.com/lewis6991/pckr.nvim',
			pckr_path
		})
	end

	vim.opt.rtp:prepend(pckr_path)
end

local function cmd(args)
	local value = 'command! -nargs=0 ' .. args
	vim.cmd(value)
end

bootstrap_pckr()

local bind = vim.keymap.set
local c = vim.cmd
local slt = { silent = true }
local normal = 'n'
local all = {'n', 'i', 'v'}
local pckr = require('pckr')

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
}

local fzf = require('fzf-lua')
local bufferline = require('bufferline')
local treesitter = require('nvim-treesitter.configs')

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
		'markdown',
		'odin',
	},
	sync_install = false,
	auto_install = false,
	highlight = {
		enable = true
	}
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
c('set number')
c('colorscheme dracula')

