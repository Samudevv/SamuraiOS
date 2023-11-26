function! neoformat#formatters#odin#enabled() abort
    return ['odinfmt']
endfunction

function! neoformat#formatters#odin#odinfmt() abort
    return {
        \ 'exe': 'odinfmt',
        \ 'args': ['-stdin', getcwd()],
        \ 'stdin': 1,
        \ 'no_append': 1,
        \ }
endfunction

