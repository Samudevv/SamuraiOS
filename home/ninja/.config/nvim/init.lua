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

local bind = vim.keymap.set

bootstrap_pckr()

require('pckr').add{
	'Mofiqul/dracula.nvim',
	'mg979/vim-visual-multi',
	'ibhagwan/fzf-lua',
	'akinsho/bufferline.nvim',
}

local fzf = require('fzf-lua')
local bufferline = require('bufferline')

vim.cmd('set number')
vim.cmd('colorscheme dracula')
fzf.setup({'fzf-native'})
vim.opt.termguicolors = true
bufferline.setup{}

bind('n', '<c-P>', '<cmd>lua require("fzf-lua").git_files()<CR>', { silent = true })
bind('n', '<Tab>', '<cmd>bnext<CR>',                              { silent = true })
bind('n', '<c-W>', '<cmd>bd<CR>',                                 { silent = true })
bind('n', '<c-Q>', '<cmd>q<CR>',                                  { silent = true })
bind('n', '<c-S>', '<cmd>w<CR>',                                  { silent = true })

