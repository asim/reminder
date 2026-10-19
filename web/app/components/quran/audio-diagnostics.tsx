import { useEffect, useState, type RefObject } from 'react';

const key = 'reminder_audio_diagnostics';
const version = import.meta.env.VITE_APP_VERSION || 'development';

// Opt-in only. Nothing is sent to a server; the bounded log stays in this tab.
export function AudioDiagnostics({ audioRef }: { audioRef: RefObject<HTMLAudioElement | null> }) {
  const [enabled, setEnabled] = useState(false);
  const [report, setReport] = useState('');
  const [worker, setWorker] = useState('checking');
  const [update, setUpdate] = useState('');
  useEffect(() => {
    if (new URLSearchParams(location.search).get('audioDebug') !== '1') return;
    setEnabled(true);
    const audio = audioRef.current;
    if (!audio) return;
    let entries: object[] = [];
    try { const saved = JSON.parse(sessionStorage.getItem(key) || '[]'); if (Array.isArray(saved)) entries = saved.slice(-80); } catch { /* storage unavailable */ }
    const record = (event: string) => {
      const buffered = Array.from({ length: audio.buffered.length }, (_, i) => [audio.buffered.start(i), audio.buffered.end(i)]);
      entries.push({ at: new Date().toISOString(), version, event, visibility: document.visibilityState,
        online: navigator.onLine, time: audio.currentTime, duration: Number.isFinite(audio.duration) ? audio.duration : null,
        paused: audio.paused, ended: audio.ended, ready: audio.readyState, network: audio.networkState,
        buffered, source: audio.currentSrc, error: audio.error && { code: audio.error.code, message: audio.error.message } });
      entries = entries.slice(-80);
      const text = JSON.stringify(entries, null, 2);
      try { sessionStorage.setItem(key, text); } catch { /* storage unavailable */ }
      setReport(JSON.stringify({ version, userAgent: navigator.userAgent }, null, 2) + '\n' + text);
    };
    const media = (event: Event) => record(event.type);
    let lastSample = 0;
    const sample = () => { if (Date.now() - lastSample >= 5000) { lastSample = Date.now(); record('progress'); } };
    const click = (event: Event) => {
      const button = (event.target as Element).closest('button');
      if (button) record('control: ' + (button.getAttribute('aria-label') || button.textContent));
    };
    const events = ['loadstart', 'loadedmetadata', 'playing', 'pause', 'waiting', 'stalled', 'suspend', 'ended', 'error', 'emptied', 'seeking', 'seeked'];
    events.forEach(event => audio.addEventListener(event, media));
    audio.addEventListener('timeupdate', sample);
    const container = audio.parentElement;
    container?.addEventListener('click', click, true);
    document.addEventListener('visibilitychange', media);
    window.addEventListener('pagehide', media);
    window.addEventListener('online', media);
    window.addEventListener('offline', media);
    record('diagnostics mounted');
    let channel: MessageChannel | undefined;
    let timeout: ReturnType<typeof setTimeout> | undefined;
    const checkWorker = () => {
      channel?.port1.close();
      if (timeout) clearTimeout(timeout);
      const controller = navigator.serviceWorker?.controller;
      if (!controller) { setWorker('no controlling service worker'); return; }
      channel = new MessageChannel();
      channel.port1.onmessage = event => { clearTimeout(timeout); setWorker(String(event.data.appVersion)); channel?.port1.close(); };
      controller.postMessage({ type: 'GET_VERSION' }, [channel.port2]);
      timeout = setTimeout(() => setWorker('older worker: no version response'), 3000);
    };
    checkWorker();
    navigator.serviceWorker?.addEventListener('controllerchange', checkWorker);
    return () => {
      record('diagnostics unmounted');
      events.forEach(event => audio.removeEventListener(event, media));
      audio.removeEventListener('timeupdate', sample);
      container?.removeEventListener('click', click, true);
      document.removeEventListener('visibilitychange', media);
      for (const event of ['pagehide', 'online', 'offline']) window.removeEventListener(event, media);
      navigator.serviceWorker?.removeEventListener('controllerchange', checkWorker);
      channel?.port1.close();
      clearTimeout(timeout);
    };
  }, [audioRef]);

  if (!enabled) return null;
  return <details className="mt-3 text-xs text-gray-600">
    <summary>Audio diagnostics</summary>
    <p>Player: {version} · Worker: {worker}</p>
    <p>{worker === version ? 'Player and worker versions match.' : 'Player and worker match has not been confirmed.'}</p>
    <button type="button" className="underline mr-3" onClick={async () => {
      setUpdate('Checking for an update…');
      try {
        const registration = await navigator.serviceWorker?.getRegistration();
        if (!registration) { setUpdate('No service worker registered. Reload while online.'); return; }
        await registration.update();
        setUpdate('Update check complete. If a worker is installing, wait for its version above to change before reloading.');
      } catch { setUpdate('Update check failed. Check your connection.'); }
    }}>Check for update</button>
    <button type="button" className="underline mr-3" onClick={() => location.reload()}>Reload page (stops audio)</button>
    <button type="button" className="underline" onClick={async () => {
      try { await navigator.clipboard.writeText(`Player: ${version}\nWorker: ${worker}\n${report}`); setUpdate('Report copied.'); }
      catch { setUpdate('Select and copy the report below.'); }
    }}>Copy report</button>
    <p role="status">{update}</p>
    <textarea aria-label="Audio diagnostic report" readOnly value={report} className="w-full h-40 font-mono" />
  </details>;
}
