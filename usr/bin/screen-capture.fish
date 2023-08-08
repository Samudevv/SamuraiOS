#! /usr/bin/fish

set screen_capture_pid (ps -e | grep wf-recorder | cut -d ' ' -f1)

if test $screen_capture_pid
  kill -s SIGINT $screen_capture_pid
  notify-send 'Recording Stopped'
else
  set geo (slurp)

  if test $geo
    notify-send -t 1000 3
    sleep 1
    notify-send -t 1000 2
    sleep 1
    notify-send -t 1000 1
    sleep 1
    notify-send -t 500 'Start Recording'

    cd /tmp
    rm recording.mp4
    wf-recorder -g $geo

    if test $status -eq 0
      notify-send 'Converting to GIF ...'
      gifski -o recording.gif recording.mp4 && rm recording.mp4
      notify-send 'Recording Done'

      eog recording.gif
    else
      notify-send -u critical 'Recording Failed'
    end
  else
    notify-send 'Recording Cancelled'
  end
end
