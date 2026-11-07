#!/bin/sh
# Test main.go's actual route wiring without production storage or model calls.
set -eu
cd "$(dirname "$0")/.."
test_dir=$(mktemp -d)
server_pid=
cleanup() {
  if [ -n "$server_pid" ]; then kill "$server_pid" 2>/dev/null || true; wait "$server_pid" 2>/dev/null || true; fi
  rm -rf "$test_dir"
}
trap cleanup EXIT INT TERM
go build -buildvcs=false -o "$test_dir/reminder" .
mkdir -p "$test_dir/home/.reminder"
printf '%s' '{"2000-01-01":{"verse":"Fixture verse","hadith":"Fixture hadith","name":"Fixture name","hijri":"Fixture date"}}' > "$test_dir/home/.reminder/daily.json"
date -u +%F > "$test_dir/home/.reminder/last_push_date.txt"
# Fail if another server owns the port rather than testing the wrong instance.
python3 -c 'import socket; s=socket.socket(); s.bind(("127.0.0.1",8080)); s.close()'
env HOME="$test_dir/home" ANTHROPIC_API_KEY= OPENAI_API_KEY= OLLAMA_LLM_MODEL= OLLAMA_BASE_URL=http://127.0.0.1:1/v1 "$test_dir/reminder" --serve > "$test_dir/server.log" 2>&1 &
server_pid=$!
if ! python3 - "$server_pid" <<'PY'
import os,sys,time,urllib.request
for _ in range(300):
    try:
        os.kill(int(sys.argv[1]), 0)
        with urllib.request.urlopen('http://127.0.0.1:8080/api/quran/1/1', timeout=1) as r:
            if r.status == 200: break
    except (OSError, ValueError): time.sleep(.1)
else: raise SystemExit('server did not become ready')
PY
then cat "$test_dir/server.log"; exit 1; fi
if ! REMINDER_TEST_URL=http://127.0.0.1:8080 go test -count=1 ./test; then
 cat "$test_dir/server.log"
 exit 1
fi
