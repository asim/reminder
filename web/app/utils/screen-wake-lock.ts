// Called only while audio is playing. Browsers may decline a wake lock (for
// example in battery saver); that must never prevent audio playback.
export function keepScreenAwake() {
  if (!('wakeLock' in navigator)) return () => {};
  let stopped = false;
  let requesting = false;
  let lock: WakeLockSentinel | undefined;
  const release = () => {
    const current = lock;
    lock = undefined;
    if (current) void current.release().catch(() => {});
  };
  const acquire = async () => {
    if (stopped || requesting || (lock && !lock.released) || document.visibilityState !== 'visible') return;
    requesting = true;
    try {
      const granted = await navigator.wakeLock.request('screen');
      if (stopped || document.visibilityState !== 'visible') {
        await granted.release();
      } else {
        lock = granted;
      }
    } catch { /* Optional capability; keep playback working if denied. */ }
    finally { requesting = false; }
  };
  const visibilityChanged = () => {
    if (document.visibilityState === 'visible') void acquire();
    else release();
  };
  document.addEventListener('visibilitychange', visibilityChanged);
  void acquire();
  return () => {
    stopped = true;
    document.removeEventListener('visibilitychange', visibilityChanged);
    release();
  };
}
