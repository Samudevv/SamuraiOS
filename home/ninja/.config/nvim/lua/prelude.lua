--  ____           _           _
-- |  _ \ _ __ ___| |_   _  __| | ___
-- | |_) | '__/ _ \ | | | |/ _` |/ _ \
-- |  __/| | |  __/ | |_| | (_| |  __/
-- |_|   |_|  \___|_|\__,_|\__,_|\___|
--
function remember_cursor_position()
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

require_failed = false
function prequire(module)
    local ok, m = pcall(require, module)
    if not ok then
        require_failed = true
        return nil
    end
    return m
end

--  ____  _                _             _
-- / ___|| |__   ___  _ __| |_ ___ _   _| |_ ___
-- \___ \| '_ \ / _ \| '__| __/ __| | | | __/ __|
--  ___) | | | | (_) | |  | || (__| |_| | |_\__ \
-- |____/|_| |_|\___/|_|   \__\___|\__,_|\__|___/
--
function cmd(args)
    local value = 'command! -nargs=0 ' .. args
    vim.cmd(value)
end
function cf(pattern, command)
    vim.api.nvim_create_autocmd({'BufEnter', 'BufWinEnter'}, {
        pattern = pattern,
        command = command,
    })
end
function ac(event, func)
    opts = {}
    if type(func) == 'string' then
        opts['command'] = func
    else
        opts['callback'] = func
    end
    vim.api.nvim_create_autocmd(event, opts)
end
function vscode_c(args)
    if is_vscode then
        c(args)
    end
end
function vscode_cmd(args)
    if is_vscode then
        cmd(args)
    end
end
function vscode_ac(args)
    if is_vscode then
        ac(args)
    end
end
function term_c(args)
    if not is_vscode then
        c(args)
    end
end
function term_cmd(args)
    if not is_vscode then
        cmd(args)
    end
end
function term_ac(event, func)
    if not is_vscode then
        ac(event, func)
    end
end

bind = vim.keymap.set
c = vim.cmd
slt = { silent = true }
normal = 'n'
all = {'n', 'i', 'v'}
lang_c = {'*.c', '*.cpp', '*.cxx', '*.h', '*.hpp', '*.hxx'}
lang_go = {'*.go'}
is_vscode = vim.g.vscode

require("packager")
